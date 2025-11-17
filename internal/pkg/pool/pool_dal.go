package pool

import (
	_ "embed"
)

//go:embed queries/cancelling_jobs_select.sql
var queryCancellingJobsSelect string

//go:embed queries/job_status_cancelled_update.sql
var queryJobStatusCancelledUpdate string

// getCancellingJobs retrieves jobs in CANCELLING state from database
func (p *Pool[T]) getCancellingJobs() []string {

	sess, err := p.db.NewSession(false)
	if err != nil {
		return nil
	}
	defer sess.Close()

	rows, err := sess.Query(queryCancellingJobsSelect)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var jobIDs []string
	for rows.Next() {
		var jobID string
		if err := rows.Scan(&jobID); err != nil {
			continue
		}
		jobIDs = append(jobIDs, jobID)
	}

	return jobIDs
}

// updateJobStatusToCancelled updates job status to CANCELLED in database
func (p *Pool[T]) updateJobStatusToCancelled(jobID string) error {
	if p.db == nil {
		return nil
	}

	sess, err := p.db.NewSession(true)
	if err != nil {
		return err
	}
	defer sess.Close()

	_, err = sess.Exec(queryJobStatusCancelledUpdate, jobID)

	if err == nil {
		return sess.Commit()
	}

	return err
}
