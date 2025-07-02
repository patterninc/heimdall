package janitor

import (
	_ "embed"

	"github.com/patterninc/heimdall/internal/pkg/database"
)

//go:embed queries/stale_jobs_select.sql
var queryStaleJobsSelect string

//go:embed queries/jobs_set_failed.sql
var queryFailStaleJobs string

//go:embed queries/stale_jobs_delete.sql
var queryStaleJobsDelete string

func (j *Janitor) cleanupStaleJobs() error {

	// let's find the jobs we'll be cleaning up...
	sess, err := j.db.NewSession(false)
	if err != nil {
		return err
	}
	defer sess.Close()

	rows, err := sess.Query(queryStaleJobsSelect, j.StaleJob)
	if err != nil {
		return err
	}
	defer rows.Close()

	staleJobIDs := make([]any, 0, 100)

	for rows.Next() {

		var jobID int

		if err := rows.Scan(&jobID); err != nil {
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
		return err
	}

	if _, err := sess.Exec(updateStaleJobs, jobSystemIDs...); err != nil {
		return err
	}

	// delete stale jobs from active jobs
	deleteStaleJobs, jobSystemIDs, err := database.PrepareSliceQuery(queryStaleJobsDelete, `$%d`, staleJobIDs)
	if err != nil {
		return err
	}

	if _, err := sess.Exec(deleteStaleJobs, jobSystemIDs...); err != nil {
		return err
	}

	return nil

}
