package heimdall

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/internal/pkg/object"
	"github.com/patterninc/heimdall/internal/pkg/object/job"
)

//go:embed queries/job/insert.sql
var queryJobInsert string

//go:embed queries/job/active_insert.sql
var queryActiveJobInsert string

//go:embed queries/job/tags_insert.sql
var queryJobTagsInsert string

//go:embed queries/job/tags_select.sql
var queryJobTagsSelect string

//go:embed queries/job/cluster_tags_insert.sql
var queryJobClusterTagsInsert string

//go:embed queries/job/cluster_tags_select.sql
var queryJobClusterTagsSelect string

//go:embed queries/job/command_tags_insert.sql
var queryJobCommandTagsInsert string

//go:embed queries/job/command_tags_select.sql
var queryJobCommandTagsSelect string

//go:embed queries/job/select.sql
var queryJobSelect string

//go:embed queries/job/select_jobs.sql
var queryJobsSelect string

//go:embed queries/job/status_select.sql
var queryJobStatusSelect string

//go:embed queries/job/select_statuses.sql
var queryJobStatusesSelect string

var (
	ErrUnknownJobID = fmt.Errorf(`unknown job_id`)
)

var (
	jobsFilterConfig = &database.FilterConfig{
		Join: " and\n    ",
		Parameters: map[string]*database.FilterParameter{
			`username`: {
				Value: `j.username like $%d`,
			},
			`id`: {
				Value: `j.job_id like $%d`,
			},
			`name`: {
				Value: `j.job_name like $%d`,
			},
			`version`: {
				Value: `j.job_version like $%d`,
			},
			`command`: {
				Value: `cm.command_name like $%d`,
			},
			`cluster`: {
				Value: `cl.cluster_name like $%d`,
			},
			`status`: {
				IsSlice: true,
				Item:    `$%d`,
				Join:    `, `,
				Value:   `js.job_status_name in ({{ .Slice }})`,
			},
		},
	}
)

type jobRequest struct {
	ID   string `yaml:"id,omitempty" json:"id,omitempty"`
	File string `yaml:"file,omitempty" json:"file,omitempty"`
}

func (h *Heimdall) insertJob(j *job.Job, clusterID, commandID string) (int64, error) {

	// open connection
	sess, err := h.Database.NewSession(true)
	if err != nil {
		return 0, err
	}
	defer sess.Close()

	// insert job row
	jobID, err := sess.InsertRow(queryJobInsert, clusterID, commandID, j.Status, j.ID, j.Name, j.Version, j.Description, j.Context.String(), j.Error, j.User, j.IsSync)
	if err != nil {
		return 0, err
	}

	// is this an async job?
	if !j.IsSync {
		if err := sess.Exec(queryActiveJobInsert, jobID); err != nil {
			return 0, err
		}
	}

	// insert job tags
	insertTagsQuery, tagItems, err := database.PrepareSliceQuery(queryJobTagsInsert, `($1, $%d)`, j.Tags.SliceAny(), jobID)
	if err != nil {
		return 0, err
	}

	if len(tagItems) > 0 {
		if err := sess.Exec(insertTagsQuery, tagItems...); err != nil {
			return 0, err
		}
	}

	// insert job cluster tags based on cluster criteria
	insertClusterTagsQuery, tagClusterItems, err := database.PrepareSliceQuery(queryJobClusterTagsInsert, `($1, $%d)`, j.ClusterCriteria.SliceAny(), jobID)
	if err != nil {
		return 0, err
	}

	if len(tagClusterItems) > 0 {
		if err := sess.Exec(insertClusterTagsQuery, tagClusterItems...); err != nil {
			return 0, err
		}
	}

	// insert job cluster tags based on cluster criteria
	insertCommandTagsQuery, tagCommandItems, err := database.PrepareSliceQuery(queryJobCommandTagsInsert, `($1, $%d)`, j.CommandCriteria.SliceAny(), jobID)
	if err != nil {
		return 0, err
	}

	if len(tagCommandItems) > 0 {
		if err := sess.Exec(insertCommandTagsQuery, tagCommandItems...); err != nil {
			return 0, err
		}
	}

	if err := sess.Commit(); err != nil {
		return 0, err
	}

	return jobID, nil

}

func (h *Heimdall) getJob(j *jobRequest) (any, error) {

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	row, err := sess.QueryRow(queryJobSelect, j.ID)

	if err != nil {
		return nil, err
	}

	r := &job.Job{
		Object: object.Object{
			ID: j.ID,
		},
	}

	var jobContext string

	if err := row.Scan(&r.SystemID, &r.Status, &r.Name, &r.Version, &r.Description, &jobContext, &r.Error, &r.User, &r.IsSync,
		&r.CreatedAt, &r.UpdatedAt, &r.CommandID, &r.CommandName, &r.CluserID, &r.ClusterName); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUnknownJobID
		} else {
			return nil, err
		}
	}

	if err := jobParseContextAndTags(r, jobContext, sess); err != nil {
		return nil, err
	}

	return r, nil

}

func (h *Heimdall) getJobs(f *database.Filter) (any, error) {

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	query, args, err := f.Render(queryJobsSelect, jobsFilterConfig)
	if err != nil {
		return nil, err
	}

	rows, err := sess.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*job.Job, 0, 100)

	for rows.Next() {

		jobContext := ``
		r := &job.Job{}

		if err := rows.Scan(&r.SystemID, &r.ID, &r.Status, &r.Name, &r.Version, &r.Description, &jobContext, &r.Error, &r.User, &r.IsSync,
			&r.CreatedAt, &r.UpdatedAt, &r.CommandID, &r.CommandName, &r.CluserID, &r.ClusterName); err != nil {
			return nil, err
		}

		if err := jobParseContextAndTags(r, jobContext, sess); err != nil {
			return nil, err
		}

		result = append(result, r)

	}

	return &resultset{
		Data: result,
	}, nil

}

func (h *Heimdall) getJobStatus(j *jobRequest) (any, error) {

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	row, err := sess.QueryRow(queryJobStatusSelect, j.ID)

	if err != nil {
		return nil, err
	}

	r := &job.Job{}

	if err := row.Scan(&r.Status, &r.Error, &r.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUnknownJobID
		} else {
			return nil, err
		}
	}

	return r, nil

}

func jobParseContextAndTags(j *job.Job, jobContext string, sess *database.Session) (err error) {

	// ...and add job context
	if jobContext != `` {
		if err = json.Unmarshal([]byte(jobContext), &j.Context); err != nil {
			return err
		}
	}

	// let's add tags
	if j.Tags, err = sess.SelectSet(queryJobTagsSelect, j.SystemID); err != nil {
		return err
	}

	// let's add cluster criteria
	if j.ClusterCriteria, err = sess.SelectSet(queryJobClusterTagsSelect, j.SystemID); err != nil {
		return err
	}

	// let's add command criteria
	if j.CommandCriteria, err = sess.SelectSet(queryJobCommandTagsSelect, j.SystemID); err != nil {
		return err
	}

	return nil

}

func (h *Heimdall) getJobStatuses(_ *database.Filter) (any, error) {

	return database.GetSlice(h.Database, queryJobStatusesSelect)

}
