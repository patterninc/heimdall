package heimdall

import (
	_ "embed"
	"fmt"

	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/object/job/status"
)

const (
	formatErrUnknownCommand = "unknown command: %s"
	formatErrUnknownCluster = "unknown cluster: %s"
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

	// open connection
	sess, err := h.Database.NewSession(true)
	if err != nil {
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

		if err := rows.Scan(&j.SystemID, &j.CommandID, &j.CluserID, &j.Status, &j.ID, &j.Name,
			&j.Version, &j.Description, &jobContext, &j.User, &j.IsSync, &j.CreatedAt, &j.UpdatedAt); err != nil {
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
			if err := sess.Exec(updateJobSliceStatusQuery, jobSystemIDs...); err != nil {
				return nil, err
			}
		}
	}

	// commit transaction
	if err := sess.Commit(); err != nil {
		return nil, err
	}

	return result, nil

}

func (h *Heimdall) runAsyncJob(j *job.Job) error {

	// let's updte job status that we're running it...
	sess, err := h.Database.NewSession(false)
	if err != nil {
		return h.updateAsyncJobStatus(j, err)
	}
	defer sess.Close()

	if err := sess.Exec(queryJobStatusUpdate, status.Running, ``, j.SystemID); err != nil {
		return h.updateAsyncJobStatus(j, err)
	}

	// do we have the command?
	command, found := h.Commands[j.CommandID]
	if !found {
		return h.updateAsyncJobStatus(j, fmt.Errorf(formatErrUnknownCommand, j.CommandID))
	}

	// do we have hte cluster?
	cluster, found := h.Clusters[j.CluserID]
	if !found {
		return h.updateAsyncJobStatus(j, fmt.Errorf(formatErrUnknownCluster, j.CluserID))
	}

	return h.updateAsyncJobStatus(j, h.runJob(j, command, cluster))

}

func (h *Heimdall) updateAsyncJobStatus(j *job.Job, jobError error) error {

	// we updte the final job status based on presence of the error
	if jobError == nil {
		j.Status = status.Succeeded
	} else {
		j.Status = status.Failed
		j.Error = jobError.Error()
	}

	// now we update that status in the database
	sess, err := h.Database.NewSession(true)
	if err != nil {
		// TODO: implement proper logging
		fmt.Println(`session error:`, err)
	}
	defer sess.Close()

	if err := sess.Exec(queryJobStatusUpdate, j.Status, j.Error, j.SystemID); err != nil {
		// TODO: implement proper logging
		fmt.Println(`job status update error:`, err)
	}

	if err := sess.Exec(queryActiveJobDelete, j.SystemID); err != nil {
		// TODO: implement proper logging
		fmt.Println(`active job delete error:`, err)
	}

	if err := sess.Commit(); err != nil {
		// TODO: implement proper logging
		fmt.Println(`session commit error:`, err)
	}

	return jobError

}
