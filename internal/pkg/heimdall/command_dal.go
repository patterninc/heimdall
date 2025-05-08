package heimdall

import (
	_ "embed"
	"encoding/json"

	_ "github.com/lib/pq"

	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/internal/pkg/object/command"
)

//go:embed queries/command/insert.sql
var queryCommandInsert string

//go:embed queries/command/tags_delete.sql
var queryCommandTagsDelete string

//go:embed queries/command/tags_insert.sql
var queryCommandTagsInsert string

//go:embed queries/command/cluster_tags_delete.sql
var queryCommandClusterTagsDelete string

//go:embed queries/command/cluster_tags_insert.sql
var queryCommandClusterTagsInsert string

//go:embed queries/command/select_statuses.sql
var queryCommandStatusesSelect string

//go:embed queries/command/select_commands.sql
var queryCommandsSelect string

//go:embed queries/command/tags_select.sql
var queryCommandTagsSelect string

//go:embed queries/command/cluster_tags_select.sql
var queryCommandClusterTagsSelect string

var (
	commandsFilterConfig = &database.FilterConfig{
		Join: " and\n    ",
		Parameters: map[string]*database.FilterParameter{
			`username`: {
				Value: `c.username like $%d`,
			},
			`id`: {
				Value: `c.command_id like $%d`,
			},
			`name`: {
				Value: `c.command_name like $%d`,
			},
			`version`: {
				Value: `c.command_version like $%d`,
			},
			`plugin`: {
				Value: `c.command_plugin like $%d`,
			},
			`status`: {
				IsSlice: true,
				Item:    `$%d`,
				Join:    `, `,
				Value:   `cs.command_status_name in ({{ .Slice }})`,
			},
		},
	}
)

func (h *Heimdall) commandInsert(c *command.Command) error {

	// open connection
	sess, err := h.Database.NewSession(true)
	if err != nil {
		return err
	}
	defer sess.Close()

	// upsert command row
	commandID, err := sess.InsertRow(queryCommandInsert, c.Status, c.ID, c.Name, c.Version, c.Plugin, c.Description, c.Context.String(), c.User, c.IsSync)
	if err != nil {
		return err
	}

	// delete all tags for the upserted command
	if err := sess.Exec(queryCommandTagsDelete, commandID); err != nil {
		return err
	}

	// insert command tags
	insertTagsQuery, tagItems, err := database.PrepareSliceQuery(queryCommandTagsInsert, `($1, $%d)`, c.Tags.SliceAny(), commandID)
	if err != nil {
		return err
	}

	if len(tagItems) > 0 {
		if err := sess.Exec(insertTagsQuery, tagItems...); err != nil {
			return err
		}
	}

	// delete all cluster tags for the upserted command
	if err := sess.Exec(queryCommandClusterTagsDelete, commandID); err != nil {
		return err
	}

	// insert command cluster tags
	insertClusterTagsQuery, clusterTagItems, err := database.PrepareSliceQuery(queryCommandClusterTagsInsert, `($1, $%d)`, c.ClusterTags.SliceAny(), commandID)
	if err != nil {
		return err
	}

	if len(clusterTagItems) > 0 {
		if err := sess.Exec(insertClusterTagsQuery, clusterTagItems...); err != nil {
			return err
		}
	}

	return sess.Commit()

}

func (h *Heimdall) getCommands(f *database.Filter) (any, error) {

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	query, args, err := f.Render(queryCommandsSelect, commandsFilterConfig)
	if err != nil {
		return nil, err
	}

	rows, err := sess.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*command.Command, 0, 100)

	for rows.Next() {

		commandContext := ``
		r := &command.Command{}

		if err := rows.Scan(&r.SystemID, &r.Status, &r.ID, &r.Name, &r.Version, &r.Plugin, &r.Description,
			&commandContext, &r.User, &r.IsSync, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}

		if err := commandParseContextAndTags(r, commandContext, sess); err != nil {
			return nil, err
		}

		result = append(result, r)

	}

	return &resultset{
		Data: result,
	}, nil

}

func (h *Heimdall) getCommandStatuses(_ *database.Filter) (any, error) {

	return database.GetSlice(h.Database, queryCommandStatusesSelect)

}

func commandParseContextAndTags(c *command.Command, commandContext string, sess *database.Session) (err error) {

	// ...and add command context
	if commandContext != `` {
		if err = json.Unmarshal([]byte(commandContext), &c.Context); err != nil {
			return err
		}
	}

	// let's add tags
	if c.Tags, err = sess.SelectSet(queryCommandTagsSelect, c.SystemID); err != nil {
		return err
	}

	// let's add cluster criteria
	if c.ClusterTags, err = sess.SelectSet(queryCommandClusterTagsSelect, c.SystemID); err != nil {
		return err
	}

	return nil

}
