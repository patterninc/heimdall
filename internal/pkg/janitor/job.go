package janitor

import (
	_ "embed"
	"fmt"

	"github.com/patterninc/heimdall/internal/pkg/database"
)

//go:embed queries/stale_jobs_select.sql
var queryStaleJobsSelect string

//go:embed queries/jobs_set_failed.sql
var queryFailStaleJobs string

//go:embed queries/stale_jobs_delete.sql
var queryStaleJobsDelete string

//go:embed queries/old_jobs_cluster_tags_delete.sql
var queryOldJobsClusterTagsDelete string

//go:embed queries/old_jobs_command_tags_delete.sql
var queryOldJobsCommandTagsDelete string

//go:embed queries/old_jobs_tags_delete.sql
var queryOldJobsTagsDelete string

//go:embed queries/old_jobs_delete.sql
var queryOldJobsDelete string

var (
	queriesForOldJobsCleanup = []string{
		queryOldJobsClusterTagsDelete,
		queryOldJobsCommandTagsDelete,
		queryOldJobsTagsDelete,
		queryOldJobsDelete,
	}
)

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

func (j *Janitor) cleanupFinishedJobs() error {
	if j.FinishedJobRetentionDays == 0 {
		return nil
	}
	// Start transactional session
	sess, err := j.db.NewSession(true)
	if err != nil {
		return err
	}
	defer sess.Close()

	defer func() {
		_ = sess.Rollback()
	}()

	exec := func(query string, args ...any) error {
		if _, err := sess.Exec(query, args...); err != nil {
			return fmt.Errorf("failed to exec query %q: %w", query, err)
		}
		return nil
	}

	for _, q := range queriesForOldJobsCleanup {
		if err := exec(q, j.FinishedJobRetentionDays); err != nil {
			return err
		}
	}

	if err := sess.Commit(); err != nil {
		return fmt.Errorf("failed to commit cleanup transaction: %w", err)
	}
	
	return nil
}
