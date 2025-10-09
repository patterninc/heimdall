package clickhouse

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/hladush/go-telemetry/pkg/telemetry"
	hdctx "github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/object/job/status"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
	"github.com/patterninc/heimdall/pkg/result/column"
)

type commandContext struct {
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
}

type clusterContext struct {
	Endpoints []string `yaml:"endpoints" json:"endpoints"`
	Database  string   `yaml:"database,omitempty" json:"database,omitempty"`
}

type jobContext struct {
	Query        string            `yaml:"query" json:"query"`
	Params       map[string]string `yaml:"params,omitempty" json:"params,omitempty"`
	ReturnResult bool              `yaml:"return_result,omitempty" json:"return_result,omitempty"`
}

type executionContext struct {
	query        string
	params       map[string]string
	returnResult bool
	conn         driver.Conn
}

const (
	serviceName = "clickhouse"
)

var (
	dummyRowsInstance    = dummyRows()
	handleMethod         = telemetry.NewMethod("handle", serviceName)
	createExcMethod      = telemetry.NewMethod("createExc", serviceName)
	collectResultsMethod = telemetry.NewMethod("collectResults", serviceName)
)

// New creates a new clickhouse plugin handler
func New(ctx *hdctx.Context) (plugin.Handler, error) {
	t := &commandContext{}

	if ctx != nil {
		if err := ctx.Unmarshal(t); err != nil {
			return nil, err
		}
	}

	return t.handler, nil
}

func (cmd *commandContext) handler(r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {
	ctx := context.Background()

	exc, err := cmd.createExecutionContext(j, c)
	if err != nil {
		handleMethod.LogAndCountError(err, "create_exc")
		return err
	}

	rows, err := exc.exec(ctx)
	if err != nil {
		handleMethod.LogAndCountError(err, "exec")
		return err
	}
	res, err := CollectResults(rows)
	if err != nil {
		handleMethod.LogAndCountError(err, "collect_results")
		return err
	}
	j.Result = res
	j.Status = status.Succeeded

	return nil
}

func (cmd *commandContext) createExecutionContext(j *job.Job, c *cluster.Cluster) (*executionContext, error) {
	// get cluster context
	clusterCtx := &clusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterCtx); err != nil {
			createExcMethod.CountError("unmarshal_cluster_context")
			return nil, fmt.Errorf("failed to unmarshal cluster context: %v", err)
		}
	}

	// get job context
	jobCtx := &jobContext{}
	if j.Context != nil {
		if err := j.Context.Unmarshal(jobCtx); err != nil {
			createExcMethod.CountError("unmarshal_job_context")
			return nil, fmt.Errorf("failed to unmarshal job context: %v", err)
		}
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: clusterCtx.Endpoints,
		Auth: clickhouse.Auth{
			Database: clusterCtx.Database,
			Username: cmd.Username,
			Password: cmd.Password,
		},
	})
	if err != nil {
		createExcMethod.CountError("open_connection")
		return nil, fmt.Errorf("failed to open ClickHouse connection: %v", err)
	}

	return &executionContext{
		query:        jobCtx.Query,
		params:       jobCtx.Params,
		returnResult: jobCtx.ReturnResult,
		conn:         conn,
	}, nil
}

func (exc *executionContext) exec(ctx context.Context) (driver.Rows, error) {
	var args []any
	for k, v := range exc.params {
		args = append(args, clickhouse.Named(k, v))
	}
	if exc.returnResult {
		return exc.conn.Query(ctx, exc.query, args...)
	}
	return dummyRowsInstance, exc.conn.Exec(ctx, exc.query, args...)

}

func CollectResults(rows driver.Rows) (*result.Result, error) {
	defer rows.Close()

	cols := rows.Columns()
	colTypes := rows.ColumnTypes()

	out := &result.Result{
		Columns: make([]*column.Column, len(cols)),
		Data:    make([][]any, 0, 128),
	}
	for i, c := range cols {
		base, _ := unwrapCHType(colTypes[i].DatabaseTypeName())
		columnTypeName := colTypes[i].DatabaseTypeName()
		if val, ok := chTypeToResultTypeName[base]; ok {
			columnTypeName = val
		}
		out.Columns[i] = &column.Column{
			Name: c,
			Type: column.Type(columnTypeName),
		}
	}

	// For each column we keep: scan target and a reader that returns a normalized interface{}

	for rows.Next() {
		scanTargets := make([]any, len(cols))
		readers := make([]func() any, len(cols))

		for i, ct := range colTypes {
			base, nullable := unwrapCHType(ct.DatabaseTypeName())

			if handler, ok := chTypeHandlers[base]; ok {
				scanTargets[i], readers[i] = handler(nullable)
			} else {
				// Fallback (covers unknown + legacy decimal detection)
				scanTargets[i], readers[i] = handleDefault(nullable)
			}
		}

		if err := rows.Scan(scanTargets...); err != nil {
			collectResultsMethod.CountError("row_scan")
			return nil, fmt.Errorf("row scan error: %w", err)
		}

		row := make([]any, len(cols))
		for i := range readers {
			row[i] = readers[i]()
		}
		out.Data = append(out.Data, row)
	}
	return out, nil
}
