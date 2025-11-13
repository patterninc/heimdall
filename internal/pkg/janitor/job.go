package janitor

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"sync"
	"time"

	"github.com/go-faster/errors"
	"github.com/hladush/go-telemetry/pkg/telemetry"
	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/pkg/object/job"
	jobStatus "github.com/patterninc/heimdall/pkg/object/job/status"
)

var (
	cleanupJobsMethod = telemetry.NewMethod("db_connection", "cleanup_jobs")
	workerMethod      = telemetry.NewMethod("janitor", "worker")
	ctx               = context.Background()
)

//go:embed queries/stale_and_canceling_jobs_select.sql
var queryStaleAndCancelingJobsSelect string

//go:embed queries/active_jobs_delete.sql
var queryActiveJobsDelete string

//go:embed queries/jobs_set_canceled.sql
var queryJobsSetCanceled string

//go:embed queries/jobs_set_failed.sql
var queryJobsSetFailed string

//go:embed queries/old_jobs_cluster_tags_delete.sql
var queryOldJobsClusterTagsDelete string

//go:embed queries/old_jobs_command_tags_delete.sql
var queryOldJobsCommandTagsDelete string

//go:embed queries/old_jobs_tags_delete.sql
var queryOldJobsTagsDelete string

//go:embed queries/old_jobs_delete.sql
var queryOldJobsDelete string

//go:embed queries/old_job_biggest_id.sql
var queryOldJobsBiggestID string

var (
	queriesForOldJobsCleanup = []string{
		queryOldJobsClusterTagsDelete,
		queryOldJobsCommandTagsDelete,
		queryOldJobsTagsDelete,
		queryOldJobsDelete,
	}
)

func (j *Janitor) worker() bool {
	// track worker cycle
	workerMethod.CountRequest()
	defer workerMethod.RecordLatency(time.Now())

	// create database session with transaction
	sess, err := j.db.NewSession(true)
	if err != nil {
		workerMethod.LogAndCountError(err, "new_session")
		return false
	}
	defer sess.Close()

	// query for stale and canceling jobs
	jobs, err := j.queryJobs(sess)
	if err != nil {
		workerMethod.LogAndCountError(err, "query_jobs")
		return false
	}

	// if no jobs found, return false
	if len(jobs) == 0 {
		return false
	}

	// call Cleanup for each job in parallel
	var wg sync.WaitGroup
	for i, jb := range jobs {
		wg.Add(1)
		go func(idx int, job *job.Job) {
			defer wg.Done()
			if err := j.cleanup(job); err != nil {
				cleanupJobsMethod.LogAndCountError(err, "cleanup")
			}
		}(i, jb)
	}

	// wait for all cleanups to complete
	wg.Wait()

	// update database: remove from active_jobs and update status
	if err := j.updateJobs(sess, jobs); err != nil {
		workerMethod.LogAndCountError(err, "update_jobs")
		return false
	}

	// commit transaction to persist changes
	if err := sess.Commit(); err != nil {
		workerMethod.LogAndCountError(err, "commit")
		return false
	}

	workerMethod.CountSuccess()

	return true

}

func (j *Janitor) queryJobs(sess *database.Session) ([]*job.Job, error) {

	// query all stale and canceling jobs
	rows, err := sess.Query(queryStaleAndCancelingJobsSelect, j.StaleJob, defaultJobLimit)
	if err != nil {
		return nil, errors.Wrap(err, "query_stale_and_canceling_jobs_select")
	}
	defer rows.Close()

	// collect jobs
	jobs := make([]*job.Job, 0)
	for rows.Next() {
		jb := &job.Job{}
		if err := rows.Scan(&jb.SystemID, &jb.ID, &jb.Status, &jb.CommandID, &jb.ClusterID); err != nil {
			workerMethod.CountError("scan")
			continue
		}
		jobs = append(jobs, jb)
	}

	return jobs, nil

}

// cleanup calls the cleanup handler for a job
func (j *Janitor) cleanup(jb *job.Job) error {

	cleanupJobsMethod.CountRequest()

	// Call cleanup handler
	handler := j.commandHandlers[jb.CommandID]
	if handler != nil {
		cluster := j.clusters[jb.ClusterID]
		if err := handler.Cleanup(ctx, jb.ID, cluster); err != nil {
			cleanupJobsMethod.CountError("cleanup_handler")
			return errors.Wrap(err, "cleanup_handler")
		}
	} else {
		// count requests for jobs that don't have a cleanup handler
		cleanupJobsMethod.CountRequest("no_cleanup_handler")
	}

	return nil

}

// updateJobs performs bulk updates: removes from active_jobs and updates job status
func (j *Janitor) updateJobs(sess *database.Session, jobs []*job.Job) error {

	// organize jobs by status for bulk updates
	allSystemIDs := make([]any, 0, len(jobs))
	cancelingSystemIDs := make([]any, 0)
	otherSystemIDs := make([]any, 0)

	for _, jb := range jobs {
		allSystemIDs = append(allSystemIDs, jb.SystemID)
		if jb.Status == jobStatus.Canceling {
			cancelingSystemIDs = append(cancelingSystemIDs, jb.SystemID)
		} else {
			// accepted or running (stale jobs)
			otherSystemIDs = append(otherSystemIDs, jb.SystemID)
		}
	}

	// bulk delete all active jobs from this set
	if len(allSystemIDs) > 0 {
		query, args, err := database.PrepareSliceQuery(queryActiveJobsDelete, "$%d", allSystemIDs)
		if err != nil {
			return errors.Wrap(err, "prepare_active_jobs_delete")
		}
		if _, err := sess.Exec(query, args...); err != nil {
			return errors.Wrap(err, "delete_active_jobs")
		}
	}

	// bulk update all canceling jobs to canceled
	if len(cancelingSystemIDs) > 0 {
		query, args, err := database.PrepareSliceQuery(queryJobsSetCanceled, "$%d", cancelingSystemIDs)
		if err != nil {
			return errors.Wrap(err, "prepare_jobs_set_canceled")
		}
		if _, err := sess.Exec(query, args...); err != nil {
			return errors.Wrap(err, "update_jobs_set_canceled")
		}
	}

	// bulk update all other jobs to failed
	if len(otherSystemIDs) > 0 {
		query, args, err := database.PrepareSliceQuery(queryJobsSetFailed, "$%d", otherSystemIDs)
		if err != nil {
			return errors.Wrap(err, "prepare_jobs_set_failed")
		}
		if _, err := sess.Exec(query, args...); err != nil {
			return errors.Wrap(err, "update_jobs_set_failed")
		}
	}

	return nil

}

func (j *Janitor) cleanupFinishedJobs() error {
	if j.FinishedJobRetentionDays <= 0 {
		return nil
	}
	// open session
	sess, err := j.db.NewSession(false)
	if err != nil {
		return err
	}
	defer sess.Close()

	retentionTimestamp := time.Now().AddDate(0, 0, -j.FinishedJobRetentionDays).Unix()

	// get biggest ID of old jobs
	row, err := sess.QueryRow(queryOldJobsBiggestID, retentionTimestamp)
	if err != nil {
		return fmt.Errorf("failed to get biggest ID of old jobs: %w", err)
	}

	var biggestID sql.NullInt64
	if err := row.Scan(&biggestID); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("failed to get biggest ID of old jobs: %w", err)
	}

	if !biggestID.Valid || biggestID.Int64 == 0 {
		return nil
	}

	// remove old jobs data
	for _, q := range queriesForOldJobsCleanup {
		for {
			affectedRows, err := sess.Exec(q, biggestID.Int64)
			if err != nil {
				return err
			}
			if affectedRows == 0 {
				break
			}
		}
	}

	return nil
}
