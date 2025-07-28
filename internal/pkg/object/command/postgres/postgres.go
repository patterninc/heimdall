package postgres

import (
	"fmt"
	"strings"
	"sync"

	"github.com/patterninc/heimdall/internal/pkg/database"
	pkgcontext "github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
	"github.com/patterninc/heimdall/pkg/result/column"
)

// postgresJobContext represents the context for a PostgreSQL job
type postgresJobContext struct {
	Query        string `yaml:"query,omitempty" json:"query,omitempty"`
	ReturnResult bool   `yaml:"return_result,omitempty" json:"return_result,omitempty"`
}

type postgresClusterContext struct {
	ConnectionString string `yaml:"connection_string,omitempty" json:"connection_string,omitempty"`
}

type postgresCommandContext struct {
	mu sync.Mutex
}

// New creates a new PostgreSQL plugin handler.
func New(_ *pkgcontext.Context) (plugin.Handler, error) {
	p := &postgresCommandContext{}
	return p.handler, nil
}

// Handler for the PostgreSQL query execution.
func (p *postgresCommandContext) handler(r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {
	jobContext, err := validateJobContext(j)
	if err != nil {
		return err
	}

	clusterContext, err := validateClusterContext(c)
	if err != nil {
		return err
	}

	db := &database.Database{ConnectionString: clusterContext.ConnectionString}

	if jobContext.ReturnResult {
		return executeSyncQuery(db, jobContext.Query, j)
	}
	return p.executeAsyncQueries(db, jobContext.Query, j)
}

func validateJobContext(j *job.Job) (*postgresJobContext, error) {
	jobContext := &postgresJobContext{}
	if j.Context != nil {
		if err := j.Context.Unmarshal(jobContext); err != nil {
			return nil, fmt.Errorf("failed to unmarshal job context: %w", err)
		}
	}
	if jobContext.Query == "" {
		return nil, fmt.Errorf("query is required in job context")
	}
	return jobContext, nil
}

func validateClusterContext(c *cluster.Cluster) (*postgresClusterContext, error) {
	clusterContext := &postgresClusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterContext); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cluster context: %w", err)
		}
	}
	if clusterContext.ConnectionString == "" {
		return nil, fmt.Errorf("connection_string is required in cluster context")
	}
	return clusterContext, nil
}

func executeSyncQuery(db *database.Database, query string, j *job.Job) error {
	// Allow a single query, even if it ends with a semicolon
	queries := splitAndTrimQueries(query)
	if len(queries) != 1 {
		return fmt.Errorf("multiple queries are not allowed when return_result is true")
	}

	sess, err := db.NewSession(false)
	if err != nil {
		return fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}
	defer sess.Close()

	rows, err := sess.Query(queries[0])
	if err != nil {
		return fmt.Errorf("PostgreSQL query execution failed: %w", err)
	}
	defer rows.Close()

	rowsResult, err := result.FromRows(rows)
	if err != nil {
		return fmt.Errorf("failed to process PostgreSQL query results: %w", err)
	}

	j.Result = rowsResult
	return nil
}

func (p *postgresCommandContext) executeAsyncQueries(db *database.Database, query string, j *job.Job) error {
	sess, err := db.NewSession(false)
	if err != nil {
		j.Result = &result.Result{
			Columns: []*column.Column{{
				Name: "error",
				Type: column.Type("string"),
			}},
			Data: [][]any{{fmt.Sprintf("Async PostgreSQL connection error: %v", err)}},
		}
		return fmt.Errorf("Async PostgreSQL connection error: %v", err)
	}
	defer sess.Close()

	_, err = sess.Exec(query)
	if err != nil {
		j.Result = &result.Result{
			Columns: []*column.Column{{
				Name: "error",
				Type: column.Type("string"),
			}},
			Data: [][]any{{fmt.Sprintf("Async PostgreSQL query execution error: %v", err)}},
		}
		return fmt.Errorf("Async PostgreSQL query execution error: %v", err)
	}

	j.Result = &result.Result{
		Columns: []*column.Column{{
			Name: "message",
			Type: column.Type("string"),
		}},
		Data: [][]any{{"All queries executed successfully"}},
	}
	return nil
}

func splitAndTrimQueries(query string) []string {
	queries := []string{}
	for _, q := range strings.Split(query, ";") {
		q = strings.TrimSpace(q)
		if q != "" {
			queries = append(queries, q)
		}
	}
	return queries
}
