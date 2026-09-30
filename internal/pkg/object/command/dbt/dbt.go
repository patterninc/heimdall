package dbt

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	heimdallContext "github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
)

// commandContext represents command-level configuration from the YAML config file
type commandContext struct {
	DbtProjectPath string `yaml:"dbt_project_path,omitempty" json:"dbt_project_path,omitempty"`
	DbtProfilesDir string `yaml:"dbt_profiles_dir,omitempty" json:"dbt_profiles_dir,omitempty"`
}

// jobContext represents job-level configuration from Airflow
type jobContext struct {
	Command     string            `json:"command"`      // build, run, test, snapshot
	Selector    string            `json:"selector,omitempty"`
	Models      []string          `json:"models,omitempty"`
	Exclude     []string          `json:"exclude,omitempty"`
	Threads     int               `json:"threads,omitempty"`
	FullRefresh bool              `json:"full_refresh,omitempty"`
	Vars        map[string]string `json:"vars,omitempty"`
	Target      string            `json:"target,omitempty"` // prod, dev, etc.
}

// clusterContext represents Snowflake cluster configuration
type clusterContext struct {
	Account    string `yaml:"account,omitempty" json:"account,omitempty"`
	User       string `yaml:"user,omitempty" json:"user,omitempty"`
	Database   string `yaml:"database,omitempty" json:"database,omitempty"`
	Warehouse  string `yaml:"warehouse,omitempty" json:"warehouse,omitempty"`
	PrivateKey string `yaml:"private_key,omitempty" json:"private_key,omitempty"`
	Role       string `yaml:"role,omitempty" json:"role,omitempty"`
	Schema     string `yaml:"schema,omitempty" json:"schema,omitempty"`
}

func New(cmdCtx *heimdallContext.Context) (plugin.Handler, error) {
	cmd := &commandContext{}

	// Parse command context from YAML config
	if cmdCtx != nil {
		if err := cmdCtx.Unmarshal(cmd); err != nil {
			return nil, err
		}
	}

	return cmd.handler, nil
}

func (c *commandContext) handler(ctx context.Context, r *plugin.Runtime, j *job.Job, cl *cluster.Cluster) error {
	// Parse cluster context (Snowflake credentials)
	clusterContext := &clusterContext{}
	if cl.Context != nil {
		if err := cl.Context.Unmarshal(clusterContext); err != nil {
			return fmt.Errorf("failed to unmarshal cluster context: %w", err)
		}
	}

	// Parse job context (dbt command parameters)
	jobContext := &jobContext{}
	if err := j.Context.Unmarshal(jobContext); err != nil {
		return fmt.Errorf("failed to unmarshal job context: %w", err)
	}

	// Set defaults
	if jobContext.Threads == 0 {
		jobContext.Threads = 4
	}
	if jobContext.Target == "" {
		jobContext.Target = "prod"
	}
	if jobContext.Command == "" {
		jobContext.Command = "build"
	}

	// Generate profiles.yml
	profilesContent, err := generateProfilesYML(clusterContext, jobContext.Target)
	if err != nil {
		return fmt.Errorf("failed to generate profiles.yml: %w", err)
	}

	// Ensure profiles directory exists
	if err := os.MkdirAll(c.DbtProfilesDir, 0755); err != nil {
		return fmt.Errorf("failed to create profiles directory: %w", err)
	}

	// Write profiles.yml
	profilesPath := filepath.Join(c.DbtProfilesDir, "profiles.yml")
	if err := os.WriteFile(profilesPath, []byte(profilesContent), 0600); err != nil {
		return fmt.Errorf("failed to write profiles.yml: %w", err)
	}

	// Build dbt command
	dbtCmd := buildDbtCommand(jobContext, c.DbtProjectPath, c.DbtProfilesDir)

	// Execute dbt command
	cmd := exec.CommandContext(ctx, dbtCmd[0], dbtCmd[1:]...)
	cmd.Dir = c.DbtProjectPath
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("DBT_PROFILES_DIR=%s", c.DbtProfilesDir),
	)

	// Redirect output to runtime's stdout/stderr
	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr

	// Execute and wait
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dbt command failed: %w", err)
	}

	return nil
}

func generateProfilesYML(clusterCtx *clusterContext, target string) (string, error) {
	if clusterCtx.Account == "" || clusterCtx.User == "" || clusterCtx.Database == "" || clusterCtx.Warehouse == "" {
		return "", fmt.Errorf("missing required Snowflake connection parameters")
	}

	if clusterCtx.PrivateKey == "" {
		return "", fmt.Errorf("missing private_key path")
	}

	role := clusterCtx.Role
	if role == "" {
		role = "PUBLIC" // Default role if not specified
	}

	schema := clusterCtx.Schema
	if schema == "" {
		schema = "PUBLIC" // Default schema if not specified
	}

	// Generate profiles.yml content
	profiles := fmt.Sprintf(`analytics_poc_profile:
  target: %s
  outputs:
    prod:
      type: snowflake
      account: %s
      user: %s
      role: %s
      database: %s
      warehouse: %s
      schema: %s
      authenticator: snowflake_jwt
      private_key_path: %s
      threads: 4
      client_session_keep_alive: false
    dev:
      type: snowflake
      account: %s
      user: %s
      role: %s
      database: analytics_dev_db
      warehouse: %s
      schema: %s
      authenticator: snowflake_jwt
      private_key_path: %s
      threads: 4
      client_session_keep_alive: false
`,
		target,
		// prod output
		clusterCtx.Account,
		clusterCtx.User,
		role,
		clusterCtx.Database,
		clusterCtx.Warehouse,
		schema,
		clusterCtx.PrivateKey,
		// dev output
		clusterCtx.Account,
		clusterCtx.User,
		role,
		clusterCtx.Warehouse,
		schema,
		clusterCtx.PrivateKey,
	)

	return profiles, nil
}

func buildDbtCommand(jobCtx *jobContext, projectPath, profilesDir string) []string {
	cmd := []string{"dbt", jobCtx.Command}

	// Add selector
	if jobCtx.Selector != "" {
		cmd = append(cmd, "--selector", jobCtx.Selector)
	}

	// Add models
	if len(jobCtx.Models) > 0 {
		cmd = append(cmd, "--models")
		cmd = append(cmd, jobCtx.Models...)
	}

	// Add exclude
	if len(jobCtx.Exclude) > 0 {
		cmd = append(cmd, "--exclude")
		cmd = append(cmd, jobCtx.Exclude...)
	}

	// Add threads
	if jobCtx.Threads > 0 {
		cmd = append(cmd, "--threads", fmt.Sprintf("%d", jobCtx.Threads))
	}

	// Add full refresh
	if jobCtx.FullRefresh {
		cmd = append(cmd, "--full-refresh")
	}

	// Add vars
	if len(jobCtx.Vars) > 0 {
		varsStr := ""
		for k, v := range jobCtx.Vars {
			if varsStr != "" {
				varsStr += ","
			}
			varsStr += fmt.Sprintf("%s:%s", k, v)
		}
		cmd = append(cmd, "--vars", fmt.Sprintf("{%s}", varsStr))
	}

	// Add target
	if jobCtx.Target != "" {
		cmd = append(cmd, "--target", jobCtx.Target)
	}

	// Add project dir
	cmd = append(cmd, "--project-dir", projectPath)

	// Add profiles dir
	cmd = append(cmd, "--profiles-dir", profilesDir)

	return cmd
}
