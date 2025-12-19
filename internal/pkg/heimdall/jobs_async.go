package heimdall

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/hladush/go-telemetry/pkg/telemetry"
	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/object/job/status"
)

const (
	formatErrUnknownCommand = "unknown command: %s"
	formatErrUnknownCluster = "unknown cluster: %s"
)

var (
	getAsyncJobsMethod         = telemetry.NewMethod("db_connection", "get_async_jobs")
	runAsyncJobMethod          = telemetry.NewMethod("db_connection", "run_async_job")
	updateAsyncJobStatusMethod = telemetry.NewMethod("db_connection", "update_async_job_status")
)

//go:embed queries/job/active_select.sql
var queryActiveJobSelect string

//go:embed queries/job/status_slice_update.sql
var queryJobSliceStatusUpdate string

//go:embed queries/job/status_update.sql
var queryJobStatusUpdate string

//go:embed queries/job/active_delete.sql
var queryActiveJobDelete string

func (h *Heimdall) getAsyncJobs(limit int) ([]*job.Job, error) {

	// Track DB connection for async jobs retrieval
	defer getAsyncJobsMethod.RecordLatency(time.Now())
	getAsyncJobsMethod.CountRequest()

	// open connection
	sess, err := h.Database.NewSession(true)
	if err != nil {
		getAsyncJobsMethod.LogAndCountError(err, "new_session")
		return nil, err
	}
	defer sess.Close()

	// let's get a batch of async jobs
	rows, err := sess.Query(queryActiveJobSelect, limit, h.agentName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// process rows
	result := make([]*job.Job, 0, limit)
	jobIDs := make([]any, 0, limit)

	for rows.Next() {

		jobContext, j := ``, &job.Job{}

		if err := rows.Scan(&j.SystemID, &j.CommandID, &j.ClusterID, &j.Status, &j.ID, &j.Name,
			&j.Version, &j.Description, &jobContext, &j.User, &j.IsSync, &j.CreatedAt, &j.UpdatedAt, &j.StoreResultSync); err != nil {
			return nil, err
		}

		// ...and add job context
		if err := jobParseContextAndTags(j, jobContext, sess); err != nil {
			return nil, err
		}

		result = append(result, j)
		jobIDs = append(jobIDs, j.SystemID)

	}

	if len(result) > 0 {
		// let's update status of the jobs
		updateJobSliceStatusQuery, jobSystemIDs, err := database.PrepareSliceQuery(queryJobSliceStatusUpdate, `$%d`, jobIDs, status.Accepted)
		if err != nil {
			return nil, err
		}

		if len(jobSystemIDs) > 0 {
			if _, err := sess.Exec(updateJobSliceStatusQuery, jobSystemIDs...); err != nil {
				return nil, err
			}
		}
	}

	// commit transaction
	if err := sess.Commit(); err != nil {
		getAsyncJobsMethod.LogAndCountError(err, "commit")
		return nil, err
	}

	getAsyncJobsMethod.CountSuccess()
	return result, nil

}

func (h *Heimdall) runAsyncJob(ctx context.Context, j *job.Job) error {

	// Track DB connection for async job execution
	defer runAsyncJobMethod.RecordLatency(time.Now())
	runAsyncJobMethod.CountRequest()

	// let's updte job status that we're running it...
	sess, err := h.Database.NewSession(false)
	if err != nil {
		runAsyncJobMethod.LogAndCountError(err, "new_session")
		return h.updateAsyncJobStatus(j, err)
	}
	defer sess.Close()

	if _, err := sess.Exec(queryJobStatusUpdate, status.Running, ``, j.SystemID); err != nil {
		return h.updateAsyncJobStatus(j, err)
	}

	// do we have the command?
	command, found := h.Commands[j.CommandID]
	if !found {
		return h.updateAsyncJobStatus(j, fmt.Errorf(formatErrUnknownCommand, j.CommandID))
	}

	// do we have hte cluster?
	cluster, found := h.Clusters[j.ClusterID]
	if !found {
		return h.updateAsyncJobStatus(j, fmt.Errorf(formatErrUnknownCluster, j.ClusterID))
	}

	runAsyncJobMethod.CountSuccess()
	return h.updateAsyncJobStatus(j, h.runJob(ctx, j, command, cluster))

}

func (h *Heimdall) updateAsyncJobStatus(j *job.Job, jobError error) error {

	// Track DB connection for async job status update
	defer updateAsyncJobStatusMethod.RecordLatency(time.Now())
	updateAsyncJobStatusMethod.CountRequest()

	// now we update that status in the database
	sess, err := h.Database.NewSession(true)
	if err != nil {
		updateAsyncJobStatusMethod.LogAndCountError(err, "new_session")
		fmt.Println(`session error:`, err)
		return jobError // Return early if session creation fails
	}
	defer sess.Close()

	if _, err := sess.Exec(queryJobStatusUpdate, j.Status, j.Error, j.SystemID); err != nil {
		// TODO: implement proper logging
		fmt.Println(`job status update error:`, err)
	}

	if _, err := sess.Exec(queryActiveJobDelete, j.SystemID); err != nil {
		// TODO: implement proper logging
		fmt.Println(`active job delete error:`, err)
	}

	if err := sess.Commit(); err != nil {
		// TODO: implement proper logging
		fmt.Println(`session commit error:`, err)
	}

	updateAsyncJobStatusMethod.CountSuccess()
	return jobError

}
