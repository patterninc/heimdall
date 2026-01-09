package janitor

import (
	"context"
	_ "embed"
	"time"

	"github.com/go-faster/errors"
	"github.com/hladush/go-telemetry/pkg/telemetry"
	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/pkg/object/job"
	jobStatus "github.com/patterninc/heimdall/pkg/object/job/status"
)

var (
	cleanupJobsMethod       = telemetry.NewMethod("db_connection", "cleanup_jobs")
	queryAndSendJobsMethod  = telemetry.NewMethod("db_connection", "query_and_send_jobs")
	cleanupWorkerJobsMethod = telemetry.NewMethod("janitor", "cleanup_worker_jobs")
	ctx                     = context.Background()
)

//go:embed queries/stale_and_canceling_jobs_select.sql
var queryStaleAndCancelingJobsSelect string

//go:embed queries/active_jobs_delete.sql
var queryActiveJobsDelete string

//go:embed queries/jobs_set_canceled.sql
var queryJobsSetCanceled string

//go:embed queries/jobs_set_failed.sql
var queryJobsSetFailed string

func (j *Janitor) cleanup(jb *job.Job) error {

	cleanupJobsMethod.CountRequest()

	// Call cleanup handler
	handler := j.commandHandlers[jb.CommandID]
	if handler.CleanupHandler != nil {
		cluster := j.clusters[jb.ClusterID]
		if err := handler.CleanupHandler(ctx, jb, cluster); err != nil {
			cleanupJobsMethod.CountError("cleanup_handler")
			return errors.Wrap(err, "cleanup_handler")
		}
	} else {
		// count requests for jobs that don't have a cleanup handler
		cleanupJobsMethod.CountRequest("no_cleanup_handler")
	}

	return nil

}

func (j *Janitor) queryAndSendJobs(jobChan chan<- *job.Job) error {

	// track session creation
	defer queryAndSendJobsMethod.RecordLatency(time.Now())
	queryAndSendJobsMethod.CountRequest()

	sess, err := j.db.NewSession(true)
	if err != nil {
		queryAndSendJobsMethod.CountError("new_session")
		return errors.Wrap(err, "new_session")
	}
	defer sess.Close()

	// query all stale and canceling jobs and return job ids and statuses
	rows, err := sess.Query(queryStaleAndCancelingJobsSelect, j.StaleJob, defaultJobLimit)
	if err != nil {
		queryAndSendJobsMethod.CountError("query_stale_and_canceling_jobs_select")
		return errors.Wrap(err, "query_stale_and_canceling_jobs_select")
	}
	defer rows.Close()

	// collect all jobs
	jobs := make([]*job.Job, 0)
	allSystemIDs := make([]any, 0)
	cancelingSystemIDs := make([]any, 0)
	otherSystemIDs := make([]any, 0)

	// organize jobs by status
	for rows.Next() {
		jb := &job.Job{}
		if err := rows.Scan(&jb.SystemID, &jb.ID, &jb.Status, &jb.CommandID, &jb.ClusterID); err != nil {
			queryAndSendJobsMethod.CountError("scan")
			continue
		}

		jobs = append(jobs, jb)
		allSystemIDs = append(allSystemIDs, jb.SystemID)

		// separate by status for bulk updates
		if jb.Status == jobStatus.Canceling {
			cancelingSystemIDs = append(cancelingSystemIDs, jb.SystemID)
		} else {
			// accepted or running (stale jobs)
			otherSystemIDs = append(otherSystemIDs, jb.SystemID)
		}
	}

	// if no jobs found, return early
	if len(jobs) == 0 {
		return nil
	}

	// bulk delete all active jobs from this set of jobs
	if len(allSystemIDs) > 0 {
		query, args, err := database.PrepareSliceQuery(queryActiveJobsDelete, "$%d", allSystemIDs)
		if err != nil {
			queryAndSendJobsMethod.CountError("prepare_active_jobs_delete")
			return errors.Wrap(err, "prepare_active_jobs_delete")
		}
		if _, err := sess.Exec(query, args...); err != nil {
			queryAndSendJobsMethod.CountError("delete_active_jobs")
			return errors.Wrap(err, "delete_active_jobs")
		}
	}

	// bulk update all canceling jobs to canceled
	if len(cancelingSystemIDs) > 0 {
		query, args, err := database.PrepareSliceQuery(queryJobsSetCanceled, "$%d", cancelingSystemIDs)
		if err != nil {
			queryAndSendJobsMethod.CountError("prepare_jobs_set_canceled")
			return errors.Wrap(err, "prepare_jobs_set_canceled")
		}
		if _, err := sess.Exec(query, args...); err != nil {
			queryAndSendJobsMethod.CountError("update_jobs_set_canceled")
			return errors.Wrap(err, "update_jobs_set_canceled")
		}
	}

	// bulk update all other jobs to failed
	if len(otherSystemIDs) > 0 {
		query, args, err := database.PrepareSliceQuery(queryJobsSetFailed, "$%d", otherSystemIDs)
		if err != nil {
			queryAndSendJobsMethod.CountError("prepare_jobs_set_failed")
			return errors.Wrap(err, "prepare_jobs_set_failed")
		}
		if _, err := sess.Exec(query, args...); err != nil {
			queryAndSendJobsMethod.CountError("update_jobs_set_failed")
			return errors.Wrap(err, "update_jobs_set_failed")
		}
	}

	// commit transaction to release locks and persist changes
	if err := sess.Commit(); err != nil {
		queryAndSendJobsMethod.CountError("commit")
		return errors.Wrap(err, "commit")
	}

	// send jobs to cleanup job channel
	for _, jb := range jobs {
		jobChan <- jb
	}

	return nil

}

func (j *Janitor) cleanupWorker(jobChan <-chan *job.Job) {

	// process jobs from channel
	for jb := range jobChan {
		if err := j.cleanup(jb); err != nil {
			cleanupWorkerJobsMethod.LogAndCountError(err, "cleanup")
			continue // don't fail if cleanup fails
		}
	}

}
