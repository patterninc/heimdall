package heimdall

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hladush/go-telemetry/pkg/telemetry"
	_ "github.com/lib/pq"

	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/pkg/object"
	"github.com/patterninc/heimdall/pkg/object/command"
)

//go:embed queries/command/upsert.sql
var queryCommandUpsert string

//go:embed queries/command/select.sql
var queryCommandSelect string

//go:embed queries/command/status_select.sql
var queryCommandStatusSelect string

//go:embed queries/command/status_update.sql
var queryCommandStatusUpdate string

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

var (
	ErrUnknownCommandID = fmt.Errorf(`unknown command_id`)
	commandDALMetrics   = telemetry.NewMethod("db_connection", "command_dal")
)

func (h *Heimdall) submitCommand(c *command.Command) (any, error) {

	if err := h.commandUpsert(c); err != nil {
		return nil, err
	}

	return h.getCommand(&command.Command{Object: object.Object{ID: c.ID}})

}

func (h *Heimdall) commandUpsert(c *command.Command) error {

	// Track DB connection for command upsert operation
	defer commandDALMetrics.RecordLatency(time.Now(), "operation", "upsert_command")
	commandDALMetrics.CountRequest("operation", "upsert_command")

	// open connection
	sess, err := h.Database.NewSession(true)
	if err != nil {
		commandDALMetrics.LogAndCountError(err, "operation", "upsert_command")
		return err
	}
	defer sess.Close()

	// upsert command row
	commandID, err := sess.InsertRow(queryCommandUpsert, c.Status, c.ID, c.Name, c.Version, c.Plugin, c.Description, c.Context.String(), c.User, c.IsSync)
	if err != nil {
		return err
	}

	// delete all tags for the upserted command
	if _, err := sess.Exec(queryCommandTagsDelete, commandID); err != nil {
		return err
	}

	// insert command tags
	insertTagsQuery, tagItems, err := database.PrepareSliceQuery(queryCommandTagsInsert, `($1, $%d)`, c.Tags.SliceAny(), commandID)
	if err != nil {
		return err
	}

	if len(tagItems) > 0 {
		if _, err := sess.Exec(insertTagsQuery, tagItems...); err != nil {
			return err
		}
	}

	// delete all cluster tags for the upserted command
	if _, err := sess.Exec(queryCommandClusterTagsDelete, commandID); err != nil {
		return err
	}

	// insert command cluster tags
	insertClusterTagsQuery, clusterTagItems, err := database.PrepareSliceQuery(queryCommandClusterTagsInsert, `($1, $%d)`, c.ClusterTags.SliceAny(), commandID)
	if err != nil {
		return err
	}

	if len(clusterTagItems) > 0 {
		if _, err := sess.Exec(insertClusterTagsQuery, clusterTagItems...); err != nil {
			return err
		}
	}

	if err := sess.Commit(); err != nil {
		commandDALMetrics.LogAndCountError(err, "operation", "upsert_command")
		return err
	}

	commandDALMetrics.CountSuccess("operation", "upsert_command")
	return nil

}

func (h *Heimdall) getCommand(c *command.Command) (any, error) {

	// Track DB connection for get command operation
	defer commandDALMetrics.RecordLatency(time.Now(), "operation", "get_command")
	commandDALMetrics.CountRequest("operation", "get_command")

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		commandDALMetrics.LogAndCountError(err, "operation", "get_command")
		return nil, err
	}
	defer sess.Close()

	row, err := sess.QueryRow(queryCommandSelect, c.ID)

	if err != nil {
		return nil, err
	}

	r := &command.Command{
		Object: object.Object{
			ID: c.ID,
		},
	}

	var commandContext string

	if err := row.Scan(&r.SystemID, &r.Status, &r.Name, &r.Version, &r.Plugin, &r.Description, &commandContext,
		&r.User, &r.IsSync, &r.CreatedAt, &r.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			commandDALMetrics.LogAndCountError(ErrUnknownCommandID, "operation", "get_command")
			return nil, ErrUnknownCommandID
		} else {
			commandDALMetrics.LogAndCountError(err, "operation", "get_command")
			return nil, err
		}
	}

	if err := commandParseContextAndTags(r, commandContext, sess); err != nil {
		commandDALMetrics.LogAndCountError(err, "operation", "get_command")
		return nil, err
	}

	commandDALMetrics.CountSuccess("operation", "get_command")
	return r, nil

}

func (h *Heimdall) getCommandStatus(c *command.Command) (any, error) {

	// Track DB connection for command status operation
	defer commandDALMetrics.RecordLatency(time.Now(), "operation", "get_command_status")
	commandDALMetrics.CountRequest("operation", "get_command_status")

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		commandDALMetrics.LogAndCountError(err, "operation", "get_command_status")
		return nil, err
	}
	defer sess.Close()

	row, err := sess.QueryRow(queryCommandStatusSelect, c.ID)

	if err != nil {
		return nil, err
	}

	r := &command.Command{}

	if err := row.Scan(&r.Status, &r.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			commandDALMetrics.LogAndCountError(ErrUnknownCommandID, "operation", "get_command_status")
			return nil, ErrUnknownCommandID
		} else {
			commandDALMetrics.LogAndCountError(err, "operation", "get_command_status")
			return nil, err
		}
	}

	commandDALMetrics.CountSuccess("operation", "get_command_status")
	return r, nil

}

func (h *Heimdall) updateCommandStatus(c *command.Command) (any, error) {

	// Track DB connection for command status update operation
	commandDALMetrics.CountRequest("operation", "update_command_status")

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		commandDALMetrics.LogAndCountError(err, "operation", "update_command_status")
		return nil, err
	}
	defer sess.Close()

	rowsAffected, err := sess.Exec(queryCommandStatusUpdate, c.ID, c.Status, c.User)
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUnknownCommandID
	}

	return h.getCommandStatus(c)

}

func (h *Heimdall) getCommands(f *database.Filter) (any, error) {

	// Track DB connection for commands list operation
	defer commandDALMetrics.RecordLatency(time.Now(), "operation", "get_commands")
	commandDALMetrics.CountRequest("operation", "get_commands")

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		commandDALMetrics.LogAndCountError(err, "operation", "get_commands")
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

	commandDALMetrics.CountSuccess("operation", "get_commands")
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
