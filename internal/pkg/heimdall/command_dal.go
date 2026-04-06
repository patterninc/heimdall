package heimdall

import (
	"context"
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
	ErrUnknownCommandID       = fmt.Errorf(`unknown command_id`)
	upsertCommandMethod       = telemetry.NewMethod("db_connection", "upsert_command")
	getCommandMethod          = telemetry.NewMethod("db_connection", "get_command")
	getCommandStatusMethod    = telemetry.NewMethod("db_connection", "get_command_status")
	updateCommandStatusMethod = telemetry.NewMethod("db_connection", "update_command_status")
	getCommandsMethod         = telemetry.NewMethod("db_connection", "get_commands")
)

func (h *Heimdall) submitCommand(ctx context.Context, c *command.Command) (any, error) {

	if err := h.commandUpsert(c); err != nil {
		return nil, err
	}

	return h.getCommand(ctx, &command.Command{Object: object.Object{ID: c.ID}})

}

func (h *Heimdall) commandUpsert(c *command.Command) error {

	// Track DB connection for command upsert operation
	defer upsertCommandMethod.RecordLatency(time.Now())
	upsertCommandMethod.CountRequest()

	// open connection
	sess, err := h.Database.NewSession(true)
	if err != nil {
		upsertCommandMethod.LogAndCountError(err, "new_session")
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
		upsertCommandMethod.LogAndCountError(err, "commit")
		return err
	}

	upsertCommandMethod.CountSuccess()
	return nil

}

func (h *Heimdall) getCommand(ctx context.Context, c *command.Command) (any, error) {

	// Track DB connection for get command operation
	defer getCommandMethod.RecordLatency(time.Now())
	getCommandMethod.CountRequest()

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		getCommandMethod.LogAndCountError(err, "new_session")
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
			return nil, ErrUnknownCommandID
		} else {
			return nil, err
		}
	}

	if err := commandParseContextAndTags(r, commandContext, sess); err != nil {
		getCommandMethod.LogAndCountError(err, "command_parse_context_and_tags")
		return nil, err
	}

	getCommandMethod.CountSuccess()
	return r, nil

}

func (h *Heimdall) getCommandStatus(ctx context.Context, c *command.Command) (any, error) {

	// Track DB connection for command status operation
	defer getCommandStatusMethod.RecordLatency(time.Now())
	getCommandStatusMethod.CountRequest()

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		getCommandStatusMethod.LogAndCountError(err, "new_session")
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
			return nil, ErrUnknownCommandID
		} else {
			return nil, err
		}
	}

	getCommandStatusMethod.CountSuccess()
	return r, nil

}

func (h *Heimdall) updateCommandStatus(ctx context.Context, c *command.Command) (any, error) {

	// Track DB connection for command status update operation
	defer updateCommandStatusMethod.RecordLatency(time.Now())
	updateCommandStatusMethod.CountRequest()

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		updateCommandStatusMethod.LogAndCountError(err, "new_session")
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

	updateCommandStatusMethod.CountSuccess()
	return h.getCommandStatus(ctx, c)

}

func (h *Heimdall) getCommands(ctx context.Context, f *database.Filter) (any, error) {

	// Track DB connection for commands list operation
	defer getCommandsMethod.RecordLatency(time.Now())
	getCommandsMethod.CountRequest()

	// open connection
	sess, err := h.Database.NewSession(false)
	if err != nil {
		getCommandsMethod.LogAndCountError(err, "new_session")
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

	getCommandsMethod.CountSuccess()
	return &resultset{
		Data: result,
	}, nil

}

func (h *Heimdall) getCommandStatuses(ctx context.Context, _ *database.Filter) (any, error) {

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
