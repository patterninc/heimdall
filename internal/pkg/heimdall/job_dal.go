package heimdall

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hladush/go-telemetry/pkg/telemetry"
	_ "github.com/lib/pq"

	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/pkg/object"
	"github.com/patterninc/heimdall/pkg/object/job"
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
	ErrUnknownJobID    = fmt.Errorf(`unknown job_id`)
	ErrInvalidCursor   = fmt.Errorf(`invalid pagination cursor`)
	insertJobMethod    = telemetry.NewMethod("db_connection", "insert_job")
	getJobMethod       = telemetry.NewMethod("db_connection", "get_job")
	getJobsMethod      = telemetry.NewMethod("db_connection", "get_jobs")
	getJobStatusMethod = telemetry.NewMethod("db_connection", "get_job_status")
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

	// Whitelisted sortable columns -> SQL expr + value type; unknown order_by falls back to system_job_id (no injection).
	jobsSortColumns = map[string]sortColumn{
		`id`:         {expr: `j.job_id`, isInt: false},
		`created_at`: {expr: `j.created_at`, isInt: true},
		`updated_at`: {expr: `j.updated_at`, isInt: true},
	}

	// tag filters use correlated EXISTS — one clause per value, all must match (AND semantics)
	jobsTagsFilterConfig = map[string]string{
		`tags`: `exists (select 1 from job_tags jt where jt.system_job_id = j.system_job_id and jt.job_tag = $%d)`,
	}
)

type jobRequest struct {
	ID   string `yaml:"id,omitempty" json:"id,omitempty"`
	File string `yaml:"file,omitempty" json:"file,omitempty"`
	User string `yaml:"user,omitempty" json:"user,omitempty"`
}

func (h *Heimdall) insertJob(j *job.Job, clusterID, commandID string) (int64, error) {

	// Track DB connection for job insert operation
	defer insertJobMethod.RecordLatency(time.Now())
	insertJobMethod.CountRequest()

	// open connection
	sess, err := h.Database.NewSession(true)
	if err != nil {
		insertJobMethod.LogAndCountError(err, "new_session")
		return 0, err
	}
	defer sess.Close()

	// insert job row
	jobID, err := sess.InsertRow(queryJobInsert, clusterID, commandID, j.Status, j.ID, j.Name, j.Version, j.Description, j.Context.String(), j.Error, j.User, j.IsSync, j.StoreResultSync, j.CanceledBy)
	if err != nil {
		return 0, err
	}

	// is this an async job?
	if !j.IsSync {
		if _, err := sess.Exec(queryActiveJobInsert, jobID); err != nil {
			return 0, err
		}
	}

	// insert job tags
	insertTagsQuery, tagItems, err := database.PrepareSliceQuery(queryJobTagsInsert, `($1, $%d)`, j.Tags.SliceAny(), jobID)
	if err != nil {
		return 0, err
	}

	if len(tagItems) > 0 {
		if _, err := sess.Exec(insertTagsQuery, tagItems...); err != nil {
			return 0, err
		}
	}

	// insert job cluster tags based on cluster criteria
	insertClusterTagsQuery, tagClusterItems, err := database.PrepareSliceQuery(queryJobClusterTagsInsert, `($1, $%d)`, j.ClusterCriteria.SliceAny(), jobID)
	if err != nil {
		return 0, err
	}

	if len(tagClusterItems) > 0 {
		if _, err := sess.Exec(insertClusterTagsQuery, tagClusterItems...); err != nil {
			return 0, err
		}
	}

	// insert job cluster tags based on cluster criteria
	insertCommandTagsQuery, tagCommandItems, err := database.PrepareSliceQuery(queryJobCommandTagsInsert, `($1, $%d)`, j.CommandCriteria.SliceAny(), jobID)
	if err != nil {
		return 0, err
	}

	if len(tagCommandItems) > 0 {
		if _, err := sess.Exec(insertCommandTagsQuery, tagCommandItems...); err != nil {
			return 0, err
		}
	}

	if err := sess.Commit(); err != nil {
		insertJobMethod.LogAndCountError(err, "commit")
		return 0, err
	}

	insertJobMethod.CountSuccess()
	return jobID, nil

}

func (h *Heimdall) getJob(ctx context.Context, j *jobRequest) (any, error) {

	// Track DB connection for job get operation
	defer getJobMethod.RecordLatency(time.Now())
	getJobMethod.CountRequest()

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		getJobMethod.LogAndCountError(err, "new_session")
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
		&r.CreatedAt, &r.UpdatedAt, &r.CommandID, &r.CommandName, &r.ClusterID, &r.ClusterName, &r.StoreResultSync, &r.CanceledBy); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUnknownJobID
		} else {
			return nil, err
		}
	}

	if err := jobParseContextAndTags(r, jobContext, sess); err != nil {
		getJobMethod.LogAndCountError(err, "job_parse_context_and_tags")
		return nil, err
	}

	getJobMethod.CountSuccess()
	return r, nil

}

const defaultPageSize = 101

// sortColumn is a whitelisted sortable column: SQL expr and whether its value is an integer (for cursor decoding).
type sortColumn struct {
	expr  string
	isInt bool
}

type resolvedSort struct {
	orderKey  string     
	column    sortColumn 
	sorted    bool       
	expr      string     
	direction string     
	cmp       string     
}

func resolveSortColumns(f *database.Filter) resolvedSort {
	orderKey := (*f)[`order_by`]
	sortCol, sorted := jobsSortColumns[orderKey]
	expr := `j.system_job_id`
	if sorted {
		expr = sortCol.expr
	}
	delete(*f, `order_by`)

	direction, cmp := `desc`, `<`
	if strings.EqualFold((*f)[`direction`], `asc`) {
		direction, cmp = `asc`, `>`
	}
	delete(*f, `direction`)

	return resolvedSort{
		orderKey:  orderKey,
		column:    sortCol,
		sorted:    sorted,
		expr:      expr,
		direction: direction,
		cmp:       cmp,
	}
}

// jobsCursor is the opaque keyset position: last row's sort value + system_job_id tiebreaker, base64-JSON encoded.
type jobsCursor struct {
	Value any   `json:"v"`
	ID    int64 `json:"i"`
}

func encodeJobsCursor(value any, id int64) string {
	// json.Marshal of these primitive types never fails
	raw, _ := json.Marshal(jobsCursor{Value: value, ID: id})
	return base64.StdEncoding.EncodeToString(raw)
}

func decodeJobsCursor(s string) (*jobsCursor, error) {
	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, ErrInvalidCursor
	}
	var c jobsCursor
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, ErrInvalidCursor
	}
	return &c, nil
}

// appendWhereClause inserts a condition into the WHERE, before the trailing ORDER BY (adds WHERE if absent).
func appendWhereClause(query, clause string) string {
	idx := strings.Index(query, "\norder by")
	if idx < 0 {
		idx = len(query)
	}
	before, after := query[:idx], query[idx:]
	if strings.Contains(before, "where") {
		return before + " and\n    " + clause + after
	}
	return before + "\nwhere\n    " + clause + after
}

func injectTagsFilter(f *database.Filter, key, existsTemplate string, query string, args []any) (string, []any) {
	v, ok := (*f)[key]
	if !ok {
		return query, args
	}
	delete(*f, key)

	var clauses []string
	for _, t := range strings.Split(v, `,`) {
		if t = strings.TrimSpace(t); t != `` {
			clauses = append(clauses, fmt.Sprintf(existsTemplate, len(args)+1))
			args = append(args, t)
		}
	}
	if len(clauses) == 0 {
		return query, args
	}

	return appendWhereClause(query, strings.Join(clauses, " and\n    ")), args
}

func applyKeysetSeek(query string, args []any, sort resolvedSort, cursor *jobsCursor) (string, []any, error) {
	if cursor == nil {
		return query, args, nil
	}

	if !sort.sorted {
		query = appendWhereClause(query, fmt.Sprintf(`j.system_job_id %s $%d`, sort.cmp, len(args)+1))
		return query, append(args, cursor.ID), nil
	}

	value := cursor.Value
	if sort.column.isInt {
		fv, ok := value.(float64)
		if !ok {
			return query, args, ErrInvalidCursor
		}
		value = int64(fv)
	}

	valIdx := len(args) + 1
	args = append(args, value, cursor.ID)
	query = appendWhereClause(query, fmt.Sprintf(
		`(%s %s $%d or (%s = $%d and j.system_job_id %s $%d))`,
		sort.expr, sort.cmp, valIdx, sort.expr, valIdx, sort.cmp, valIdx+1))
	return query, args, nil
}

func (h *Heimdall) getJobs(ctx context.Context, f *database.Filter) (any, error) {

	// Track DB connection for jobs list operation
	defer getJobsMethod.RecordLatency(time.Now())
	getJobsMethod.CountRequest()

	// pull pagination/sort params out before rendering the filter WHERE
	pageSize := defaultPageSize
	if v, ok := (*f)[`limit`]; ok {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			if n > defaultPageSize {
				n = defaultPageSize
			}
			pageSize = n
		}
		delete(*f, `limit`)
	}

	sort := resolveSortColumns(f)

	// decode opaque keyset cursor (absent = first page)
	var cursor *jobsCursor
	if cursorStr := (*f)[`cursor`]; cursorStr != `` {
		c, err := decodeJobsCursor(cursorStr)
		if err != nil {
			getJobsMethod.LogAndCountError(err, "cursor_decode")
			return nil, err
		}
		cursor = c
	}
	delete(*f, `cursor`)

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		getJobsMethod.LogAndCountError(err, "new_session")
		return nil, err
	}
	defer sess.Close()

	query, args, err := f.Render(queryJobsSelect, jobsFilterConfig)
	if err != nil {
		getJobsMethod.LogAndCountError(err, "query")
		return nil, err
	}

	for key, tmpl := range jobsTagsFilterConfig {
		query, args = injectTagsFilter(f, key, tmpl, query, args)
	}

	// keyset seek: rows after the cursor, sort column + system_job_id tiebreaker (no dropped/repeated rows)
	query, args, err = applyKeysetSeek(query, args, sort, cursor)
	if err != nil {
		getJobsMethod.LogAndCountError(err, "keyset_seek")
		return nil, err
	}

	// swap in chosen ORDER BY (+ system_job_id tiebreaker); fetch one extra row to detect a next page
	orderIdx := strings.Index(query, "\norder by")
	if orderIdx < 0 {
		orderIdx = len(query)
	}
	base := strings.TrimRight(query[:orderIdx], "\n; \t")
	query = fmt.Sprintf("%s\norder by\n    %s %s, j.system_job_id %s\nlimit $%d",
		base, sort.expr, sort.direction, sort.direction, len(args)+1)
	args = append(args, pageSize+1)

	rows, err := sess.Query(query, args...)
	if err != nil {
		getJobsMethod.LogAndCountError(err, "query")
		return nil, err
	}
	defer rows.Close()

	result := make([]*job.Job, 0, pageSize+1)

	for rows.Next() {

		jobContext := ``
		r := &job.Job{}

		if err := rows.Scan(&r.SystemID, &r.ID, &r.Status, &r.Name, &r.Version, &r.Description, &jobContext, &r.Error, &r.User, &r.IsSync,
			&r.CreatedAt, &r.UpdatedAt, &r.CommandID, &r.CommandName, &r.ClusterID, &r.ClusterName, &r.StoreResultSync, &r.CanceledBy); err != nil {
			getJobsMethod.LogAndCountError(err, "scan")
			return nil, err
		}

		if err := jobParseContextAndTags(r, jobContext, sess); err != nil {
			getJobsMethod.LogAndCountError(err, "job_parse_context_and_tags")
			return nil, err
		}

		result = append(result, r)

	}

	// extra row => next page exists; trim it and emit a cursor from the last returned row
	rs := &resultset{Data: result}
	if len(result) > pageSize {
		last := result[pageSize-1]
		rs.HasMore = true
		rs.NextCursor = encodeJobsCursor(jobSortValue(sort.orderKey, last), last.SystemID)
		rs.Data = result[:pageSize]
	}

	getJobsMethod.CountSuccess()
	return rs, nil

}

// jobSortValue returns a job's value for the sort column, for the next-page cursor.
func jobSortValue(orderKey string, j *job.Job) any {
	switch orderKey {
	case `id`:
		return j.ID
	case `created_at`:
		return j.CreatedAt
	case `updated_at`:
		return j.UpdatedAt
	default:
		return j.SystemID
	}
}

func (h *Heimdall) getJobStatus(ctx context.Context, j *jobRequest) (any, error) {

	// Track DB connection for job status operation
	defer getJobStatusMethod.RecordLatency(time.Now())
	getJobStatusMethod.CountRequest()

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		getJobStatusMethod.LogAndCountError(err, "new_session")
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
			getJobStatusMethod.LogAndCountError(ErrUnknownJobID, "query")
			return nil, ErrUnknownJobID
		} else {
			getJobStatusMethod.LogAndCountError(err, "query")
			return nil, err
		}
	}

	getJobStatusMethod.CountSuccess()
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

func (h *Heimdall) getJobStatuses(ctx context.Context, _ *database.Filter) (any, error) {

	return database.GetSlice(h.Database, queryJobStatusesSelect)

}
