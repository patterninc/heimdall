package snowflake

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	sf "github.com/snowflakedb/gosnowflake"

	"github.com/patterninc/heimdall/internal/pkg/aws"
	heimdallContext "github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
)

const (
	privateKeyType      = `PRIVATE KEY`
	snowflakeDriverName = `snowflake`
	formatQueryTag      = `heimdall-job-%s`
	s3Prefix            = `s3://`
	resultFilename      = `result.json` // the file heimdall serves for a job, under its result directory

	// result_limit picks the user's most recent queries before the where clause runs, so it is set to the maximum
	queryByTagSelect = `select query_id from table(information_schema.query_history_by_user(result_limit => 10000)) where query_tag = ? order by start_time desc limit 1`
	queryCancel      = `select system$cancel_query(?)`
)

var (
	ErrFailedToDecodePEMBlock = fmt.Errorf(`failed to decode PEM block`)
	ErrInvalidKeyType         = fmt.Errorf(`invalid key type`)
)

// commandContext represents command-level configuration from the YAML config file.
// Role is defined here for security - users cannot override
type commandContext struct {
	Role string `yaml:"role,omitempty" json:"role,omitempty"`
}

type jobContext struct {
	Query string `yaml:"query,omitempty" json:"query,omitempty"`
}

type clusterContext struct {
	Account    string `yaml:"account,omitempty" json:"account,omitempty"`
	User       string `yaml:"user,omitempty" json:"user,omitempty"`
	Database   string `yaml:"database,omitempty" json:"database,omitempty"`
	Warehouse  string `yaml:"warehouse,omitempty" json:"warehouse,omitempty"`
	PrivateKey string `yaml:"private_key,omitempty" json:"private_key,omitempty"`
}

func parsePrivateKey(privateKeyBytes []byte) (*rsa.PrivateKey, error) {

	// Decode the key into a "block"
	privateBlock, _ := pem.Decode(privateKeyBytes)
	if privateBlock == nil || privateBlock.Type != privateKeyType {
		return nil, ErrFailedToDecodePEMBlock
	}

	// Parse the private key from the block
	privateKey, err := x509.ParsePKCS8PrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// Check the type of the key
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidKeyType
	}

	return rsaPrivateKey, nil

}

func New(cmdCtx *heimdallContext.Context) (plugin.Handler, error) {
	s := &commandContext{}

	// Parse command context from YAML config (contains role configuration)
	if cmdCtx != nil {
		if err := cmdCtx.Unmarshal(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Execute implements the plugin.Handler interface
func (s *commandContext) Execute(ctx context.Context, r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {
	return fmt.Errorf("snowflake: Execute not implemented for job %q -- this plugin only supports StreamResult", j.ID)
}

// StreamResult implements plugin.ResultStreamer: rows stream straight to w via result.StreamRows, so j.Result's full [][]any never exists.
func (s *commandContext) StreamResult(ctx context.Context, w io.Writer, r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {

	// let's unmarshal job context
	jobContext := &jobContext{}
	if err := j.Context.Unmarshal(jobContext); err != nil {
		return err
	}

	db, err := s.open(c, r.UserAgent)
	if err != nil {
		return err
	}
	defer db.Close()

	return streamQuery(ctx, w, db, j.ID, jobContext.Query)

}

// Resume implements plugin.Resumer. Snowflake keeps running a query after its client disconnects,
// so the job's query is found by its query tag and its result fetched by query ID; a query that was never
// submitted runs now.
func (s *commandContext) Resume(ctx context.Context, r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {

	return writeResult(ctx, r, func(w io.Writer) error {
		return s.resume(ctx, w, r, j, c)
	})

}

func (s *commandContext) resume(ctx context.Context, w io.Writer, r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {

	jobContext := &jobContext{}
	if err := j.Context.Unmarshal(jobContext); err != nil {
		return err
	}

	db, err := s.open(c, r.UserAgent)
	if err != nil {
		return err
	}
	defer db.Close()

	queryID, err := queryByTag(ctx, db, j.ID)
	if err != nil {
		return err
	}
	if queryID == `` {
		return streamQuery(ctx, w, db, j.ID, jobContext.Query)
	}

	// the driver polls a running query until it finishes; a failed or canceled query returns its error
	rows, err := db.QueryContext(sf.WithFetchResultByID(ctx, queryID), ``)
	if err != nil {
		return err
	}

	return result.StreamRows(w, rows)

}

func streamQuery(ctx context.Context, w io.Writer, db *sql.DB, jobID, query string) error {

	rows, err := db.QueryContext(sf.WithQueryTag(ctx, queryTag(jobID)), query)
	if err != nil {
		return err
	}

	return result.StreamRows(w, rows)

}

// writeResult streams the result to the file heimdall serves for the job.
func writeResult(ctx context.Context, r *plugin.Runtime, write func(io.Writer) error) error {

	name := r.ResultDirectory + "/" + resultFilename
	if strings.HasPrefix(r.ResultDirectory, s3Prefix) {
		return aws.StreamToS3(ctx, name, write)
	}

	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	writeErr := write(f)
	closeErr := f.Close()
	if writeErr != nil {
		return writeErr
	}

	return closeErr

}

func queryTag(jobID string) string {
	return fmt.Sprintf(formatQueryTag, jobID)
}

// queryByTag returns the ID of the job's query, or an empty string if it was never submitted.
func queryByTag(ctx context.Context, db *sql.DB, jobID string) (string, error) {

	var queryID string
	if err := db.QueryRowContext(ctx, queryByTagSelect, queryTag(jobID)).Scan(&queryID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return ``, err
	}

	return queryID, nil

}

// open connects to the cluster's warehouse as its user, with the command's role.
func (s *commandContext) open(c *cluster.Cluster, application string) (*sql.DB, error) {

	clusterContext := &clusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterContext); err != nil {
			return nil, err
		}
	}

	// prepare snowflake config
	privateKeyBytes, err := os.ReadFile(clusterContext.PrivateKey)
	if err != nil {
		return nil, err
	}

	// Parse the private key
	privateKey, err := parsePrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	// s.Role from command context; empty string = Snowflake uses user's default role
	dsn, err := sf.DSN(&sf.Config{
		Account:       clusterContext.Account,
		User:          clusterContext.User,
		Database:      clusterContext.Database,
		Warehouse:     clusterContext.Warehouse,
		Role:          s.Role,
		Authenticator: sf.AuthTypeJwt,
		PrivateKey:    privateKey,
		Application:   application,
	})
	if err != nil {
		return nil, err
	}

	return sql.Open(snowflakeDriverName, dsn)

}

// HealthCheck implements the plugin.HealthChecker interface
func (s *commandContext) HealthCheck(ctx context.Context, c *cluster.Cluster) error {
	clusterCtx := &clusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterCtx); err != nil {
			return err
		}
	}

	privateKeyBytes, err := os.ReadFile(clusterCtx.PrivateKey)
	if err != nil {
		return err
	}

	privateKey, err := parsePrivateKey(privateKeyBytes)
	if err != nil {
		return err
	}

	dsn, err := sf.DSN(&sf.Config{
		Account:       clusterCtx.Account,
		User:          clusterCtx.User,
		Database:      clusterCtx.Database,
		Warehouse:     clusterCtx.Warehouse,
		Role:          s.Role,
		Authenticator: sf.AuthTypeJwt,
		PrivateKey:    privateKey,
	})
	if err != nil {
		return err
	}

	db, err := sql.Open(snowflakeDriverName, dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, fmt.Sprintf("SHOW WAREHOUSES LIKE '%s'", clusterCtx.Warehouse))
	if err != nil {
		return fmt.Errorf("snowflake SHOW WAREHOUSES failed: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	stateIdx, autoResumeIdx := -1, -1
	for i, col := range cols {
		switch col {
		case "state":
			stateIdx = i
		case "auto_resume":
			autoResumeIdx = i
		}
	}
	if stateIdx < 0 {
		return fmt.Errorf("snowflake SHOW WAREHOUSES: state column not found")
	}

	if !rows.Next() {
		return fmt.Errorf("snowflake warehouse %q not found", clusterCtx.Warehouse)
	}

	vals := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	if err := rows.Scan(ptrs...); err != nil {
		return err
	}

	state := scanToString(vals[stateIdx])
	switch state {
	case "STARTED", "RESUMING":
		return nil
	case "SUSPENDED":
		if autoResumeIdx >= 0 && scanToString(vals[autoResumeIdx]) == "true" {
			return nil
		}
		return fmt.Errorf("snowflake warehouse %q is suspended and auto_resume is not enabled", clusterCtx.Warehouse)
	default:
		return fmt.Errorf("snowflake warehouse %q is not ready: state=%s", clusterCtx.Warehouse, state)
	}
}

func scanToString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Cleanup cancels the job's query if it is still running, so a query left by an agent that died stops using the warehouse.
func (s *commandContext) Cleanup(ctx context.Context, jobID string, c *cluster.Cluster) error {

	db, err := s.open(c, ``)
	if err != nil {
		return err
	}
	defer db.Close()

	queryID, err := queryByTag(ctx, db, jobID)
	if err != nil || queryID == `` {
		return err
	}

	_, err = db.ExecContext(ctx, queryCancel, queryID)
	return err

}
