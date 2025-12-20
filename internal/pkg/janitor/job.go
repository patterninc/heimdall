package janitor

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hladush/go-telemetry/pkg/telemetry"
	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/pkg/object/job"
)

var (
	cleanUpStaleJobsMethod      = telemetry.NewMethod("db_connection", "cleanup_stale_jobs")
	cleanUpCancellingJobsMethod = telemetry.NewMethod("db_connection", "cleanup_cancelling_jobs")
	cleanupMethod               = telemetry.NewMethod("db_connection", "cleanup_jobs")
	ctx                         = context.Background()
)

//go:embed queries/stale_jobs_select.sql
var queryStaleJobsSelect string

//go:embed queries/stale_jobs_delete.sql
var queryStaleJobsDelete string

//go:embed queries/cancelling_jobs_select.sql
var queryCancellingJobsSelect string

//go:embed queries/job_set_cancelled.sql
var queryJobSetCancelled string

// clean up jobs in parallel and return their system IDs for updating
func (j *Janitor) cleanupJobs(jobs []*job.Job) []any {

	var wg sync.WaitGroup
	jobIDsChan := make(chan int64, len(jobs))

	for _, jb := range jobs {
		wg.Add(1)
		go func(job *job.Job) {
			defer wg.Done()

			// look up handlers directly by command ID
			handlers := j.commandHandlers[job.CommandID]
			if handlers != nil && handlers.CleanupHandler != nil {
				// look up cluster from map
				cl, found := j.clusters[job.ClusterID]
				if found {
					// attempt to clean up job resources (log errors but continue)
					if err := handlers.CleanupHandler(ctx, job, cl); err != nil {
						cleanupMethod.LogAndCountError(fmt.Errorf("cleanup failed for job %s: %w", job.ID, err), "cleanup")
					}
				} else {
					cleanupMethod.LogAndCountError(fmt.Errorf("unknown cluster_id: %s for job %s", job.ClusterID, job.ID), "cluster_lookup")
				}
			}

			// collect job ID (regardless of cleanup success)
			jobIDsChan <- job.SystemID
		}(jb)
	}

	// wait for all cleanup operations to complete
	wg.Wait()
	close(jobIDsChan)

	// collect all job IDs from channel
	jobIDs := make([]any, 0, len(jobs))
	for id := range jobIDsChan {
		jobIDs = append(jobIDs, id)
	}

	return jobIDs
}

func (j *Janitor) cleanupStaleJobs() error {

	// use transaction for FOR UPDATE SKIP LOCKED row locking
	sess, err := j.db.NewSession(true)
	if err != nil {
		cleanUpStaleJobsMethod.LogAndCountError(err, "new_session")
		return err
	}
	defer sess.Close()

	// query stale jobs
	rows, err := sess.Query(queryStaleJobsSelect, j.StaleJob, defaultJobLimit)
	if err != nil {
		cleanUpStaleJobsMethod.LogAndCountError(err, "query")
		return err
	}
	defer rows.Close()

	// collect all jobs
	jobs := make([]*job.Job, 0, defaultJobLimit)
	for rows.Next() {
		var cancellationCtxJSON []byte
		jb := &job.Job{}

		if err := rows.Scan(&jb.SystemID, &jb.ID, &cancellationCtxJSON, &jb.CommandID, &jb.ClusterID); err != nil {
			cleanUpStaleJobsMethod.LogAndCountError(err, "scan")
			continue
		}

		// parse cancellation context (JSONB from PostgreSQL)
		if len(cancellationCtxJSON) > 0 {
			if err := json.Unmarshal(cancellationCtxJSON, &jb.CancellationCtx); err != nil {
				cleanUpStaleJobsMethod.LogAndCountError(err, "parse_cancellation_ctx")
				// continue anyway
			}
		}

		jobs = append(jobs, jb)
	}

	// no jobs found, commit and return early
	if len(jobs) == 0 {
		if err := sess.Commit(); err != nil {
			cleanUpStaleJobsMethod.LogAndCountError(err, "commit")
			return err
		}
		cleanUpStaleJobsMethod.CountSuccess()
		return nil
	}

	// clean up jobs and get their IDs
	jobIDs := j.cleanupJobs(jobs)

	// delete stale jobs from active_jobs
	deleteStaleJobs, jobSystemIDs, err := database.PrepareSliceQuery(queryStaleJobsDelete, `$%d`, jobIDs)
	if err != nil {
		cleanUpStaleJobsMethod.LogAndCountError(err, "prepare_slice_query")
		return err
	}

	if _, err := sess.Exec(deleteStaleJobs, jobSystemIDs...); err != nil {
		cleanUpStaleJobsMethod.LogAndCountError(err, "exec")
		return err
	}

	// commit transaction to release locks
	if err := sess.Commit(); err != nil {
		cleanUpStaleJobsMethod.LogAndCountError(err, "commit")
		return err
	}

	cleanUpStaleJobsMethod.CountSuccess()
	return nil
}

func (j *Janitor) cleanupCancellingJobs() error {

	// Start a new session
	sess, err := j.db.NewSession(true)
	if err != nil {
		cleanUpCancellingJobsMethod.LogAndCountError(err, "new_session")
		return err
	}
	defer sess.Close()

	// query cancelling jobs
	rows, err := sess.Query(queryCancellingJobsSelect, defaultJobLimit)
	if err != nil {
		cleanUpCancellingJobsMethod.LogAndCountError(err, "query")
		return err
	}
	defer rows.Close()

	// collect all jobs
	jobs := make([]*job.Job, 0, defaultJobLimit)
	for rows.Next() {
		var cancellationCtxJSON []byte
		jb := &job.Job{}

		if err := rows.Scan(&jb.SystemID, &jb.ID, &cancellationCtxJSON, &jb.CommandID, &jb.ClusterID); err != nil {
			cleanUpCancellingJobsMethod.LogAndCountError(err, "scan")
			continue
		}

		// parse cancellation context
		if len(cancellationCtxJSON) > 0 {
			if err := json.Unmarshal(cancellationCtxJSON, &jb.CancellationCtx); err != nil {
				cleanUpCancellingJobsMethod.LogAndCountError(err, "parse_cancellation_ctx")
				// continue anyway
			}
		}

		jobs = append(jobs, jb)
	}

	// no jobs found, commit and return early
	if len(jobs) == 0 {
		if err := sess.Commit(); err != nil {
			cleanUpCancellingJobsMethod.LogAndCountError(err, "commit")
			return err
		}
		cleanUpCancellingJobsMethod.CountSuccess()
		return nil
	}

	// clean up jobs and get their IDs
	jobIDs := j.cleanupJobs(jobs)

	// update status to cancelled
	updateQuery, jobSystemIDs, err := database.PrepareSliceQuery(queryJobSetCancelled, `$%d`, jobIDs)
	if err != nil {
		cleanUpCancellingJobsMethod.LogAndCountError(err, "prepare_slice_query")
		return err
	}

	if _, err := sess.Exec(updateQuery, jobSystemIDs...); err != nil {
		cleanUpCancellingJobsMethod.LogAndCountError(err, "exec")
		return err
	}

	// commit transaction to release locks
	if err := sess.Commit(); err != nil {
		cleanUpCancellingJobsMethod.LogAndCountError(err, "commit")
		return err
	}

	cleanUpCancellingJobsMethod.CountSuccess()
	return nil
}
