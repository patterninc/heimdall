package heimdall

import (
	_ "embed"
	"encoding/json"

	_ "github.com/lib/pq"

	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/pkg/object/cluster"
)

//go:embed queries/cluster/insert.sql
var queryClusterInsert string

//go:embed queries/cluster/tags_delete.sql
var queryClusterTagsDelete string

//go:embed queries/cluster/tags_insert.sql
var queryClusterTagsInsert string

//go:embed queries/cluster/select_statuses.sql
var queryClusterStatusesSelect string

//go:embed queries/cluster/select_clusters.sql
var queryClustersSelect string

//go:embed queries/cluster/tags_select.sql
var queryClusterTagsSelect string

var (
	clustersFilterConfig = &database.FilterConfig{
		Join: " and\n    ",
		Parameters: map[string]*database.FilterParameter{
			`username`: {
				Value: `c.username like $%d`,
			},
			`id`: {
				Value: `c.cluster_id like $%d`,
			},
			`name`: {
				Value: `c.cluster_name like $%d`,
			},
			`version`: {
				Value: `c.cluster_version like $%d`,
			},
			`status`: {
				IsSlice: true,
				Item:    `$%d`,
				Join:    `, `,
				Value:   `cs.cluster_status_name in ({{ .Slice }})`,
			},
		},
	}
)

func (h *Heimdall) clusterInsert(c *cluster.Cluster) error {

	// open connection
	sess, err := h.Database.NewSession(true)
	if err != nil {
		return err
	}
	defer sess.Close()

	// upsert cluster row
	clusterID, err := sess.InsertRow(queryClusterInsert, c.Status, c.ID, c.Name, c.Version, c.Description, c.Context.String(), c.User)
	if err != nil {
		return err
	}

	// delete all tags for the upserted cluster
	if err := sess.Exec(queryClusterTagsDelete, clusterID); err != nil {
		return err
	}

	// insert cluster tags
	insertTagsQuery, tagItems, err := database.PrepareSliceQuery(queryClusterTagsInsert, `($1, $%d)`, c.Tags.SliceAny(), clusterID)
	if err != nil {
		return err
	}

	if len(tagItems) > 0 {
		if err := sess.Exec(insertTagsQuery, tagItems...); err != nil {
			return err
		}
	}

	return sess.Commit()

}

func (h *Heimdall) getClusters(f *database.Filter) (any, error) {

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	query, args, err := f.Render(queryClustersSelect, clustersFilterConfig)
	if err != nil {
		return nil, err
	}

	rows, err := sess.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*cluster.Cluster, 0, 100)

	for rows.Next() {

		clusterContext := ``
		r := &cluster.Cluster{}

		if err := rows.Scan(&r.SystemID, &r.Status, &r.ID, &r.Name, &r.Version, &r.Description,
			&clusterContext, &r.User, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}

		if err := clusterParseContextAndTags(r, clusterContext, sess); err != nil {
			return nil, err
		}

		result = append(result, r)

	}

	return &resultset{
		Data: result,
	}, nil

}

func (h *Heimdall) getClusterStatuses(_ *database.Filter) (any, error) {

	return database.GetSlice(h.Database, queryClusterStatusesSelect)

}

func clusterParseContextAndTags(c *cluster.Cluster, clusterContext string, sess *database.Session) (err error) {

	// ...and add cluster context
	if clusterContext != `` {
		if err = json.Unmarshal([]byte(clusterContext), &c.Context); err != nil {
			return err
		}
	}

	// let's add tags
	if c.Tags, err = sess.SelectSet(queryClusterTagsSelect, c.SystemID); err != nil {
		return err
	}

	return nil

}
