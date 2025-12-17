package janitor

import (
	_ "embed"
	"time"

	"github.com/hladush/go-telemetry/pkg/telemetry"
	"github.com/patterninc/heimdall/internal/pkg/database"
)

var (
	janitorMetrics = telemetry.NewMethod("db_connection", "janitor")
)

//go:embed queries/stale_jobs_select.sql
var queryStaleJobsSelect string

//go:embed queries/jobs_set_failed.sql
var queryFailStaleJobs string

//go:embed queries/stale_jobs_delete.sql
var queryStaleJobsDelete string

func (j *Janitor) cleanupStaleJobs() error {

	// Track DB connection for janitor cleanup operations
	defer janitorMetrics.RecordLatency(time.Now(), "operation", "cleanup_stale_jobs")
	janitorMetrics.CountRequest("operation", "cleanup_stale_jobs")

	// let's find the jobs we'll be cleaning up...
	sess, err := j.db.NewSession(false)
	if err != nil {
		janitorMetrics.LogAndCountError(err, "operation", "cleanup_session_create")
		return err
	}
	defer sess.Close()

	rows, err := sess.Query(queryStaleJobsSelect, j.StaleJob)
	if err != nil {
		janitorMetrics.LogAndCountError(err, "operation", "cleanup_select_stale_jobs")
		return err
	}
	defer rows.Close()

	staleJobIDs := make([]any, 0, 100)

	for rows.Next() {

		var jobID int

		if err := rows.Scan(&jobID); err != nil {
			janitorMetrics.LogAndCountError(err, "operation", "cleanup_scan_stale_jobs")
			return err
		}

		if jobID != 0 {
			staleJobIDs = append(staleJobIDs, jobID)
		}

	}

	// do we have any stale jobs?
	if len(staleJobIDs) == 0 {
		return nil
	}

	// prepare query to update job statuses
	updateStaleJobs, jobSystemIDs, err := database.PrepareSliceQuery(queryFailStaleJobs, `$%d`, staleJobIDs)
	if err != nil {
		janitorMetrics.LogAndCountError(err, "operation", "cleanup_prepare_update_query")
		return err
	}

	if _, err := sess.Exec(updateStaleJobs, jobSystemIDs...); err != nil {
		janitorMetrics.LogAndCountError(err, "operation", "cleanup_update_stale_jobs")
		return err
	}

	// delete stale jobs from active jobs
	deleteStaleJobs, jobSystemIDs, err := database.PrepareSliceQuery(queryStaleJobsDelete, `$%d`, staleJobIDs)
	if err != nil {
		janitorMetrics.LogAndCountError(err, "operation", "cleanup_prepare_delete_query")
		return err
	}

	if _, err := sess.Exec(deleteStaleJobs, jobSystemIDs...); err != nil {
		janitorMetrics.LogAndCountError(err, "operation", "cleanup_delete_stale_jobs")
		return err
	}

	janitorMetrics.CountSuccess("operation", "cleanup_stale_jobs")
	return nil

}
