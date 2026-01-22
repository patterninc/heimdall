package shell

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path"

	heimdallContext "github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
)

const (
	resultFilename  = `result.json`
	contextFilename = `context.json`
)

type commandContext struct {
	Command []string `yaml:"command,omitempty" json:"command,omitempty"`
}

type jobContext struct {
	Arguments []string `yaml:"arguments,omitempty" json:"arguments,omitempty"`
}

type runtimeContext struct {
	Job     *job.Job         `yaml:"job,omitempty" json:"job,omitempty"`
	Command *commandContext  `yaml:"command,omitempty" json:"command,omitempty"`
	Cluster *cluster.Cluster `yaml:"cluster,omitempty" json:"cluster,omitempty"`
	Runtime *plugin.Runtime  `yaml:"runtime,omitempty" json:"runtime,omitempty"`
}

func New(commandCtx *heimdallContext.Context) (plugin.Handler, error) {

	s := &commandContext{}

	if commandCtx != nil {
		if err := commandCtx.Unmarshal(s); err != nil {
			return nil, err
		}
	}

	return s, nil

}

// Execute implements the plugin.Handler interface
func (s *commandContext) Execute(ctx context.Context, r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {

	// let's unmarshal job context
	jc := &jobContext{}
	if j.Context != nil {
		if err := j.Context.Unmarshal(jc); err != nil {
			return err
		}
	}

	// let's write job context as context.json to the job's working directory
	// this will enable us to pass any context into shell command and...
	// parse those keys in the shell script using `jq` utility
	rc := &runtimeContext{
		Job:     j,
		Command: s,
		Cluster: c,
		Runtime: r,
	}

	contextContent, err := json.MarshalIndent(rc, ``, `  `)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path.Join(r.WorkingDirectory, contextFilename), contextContent, 0600); err != nil {
		return err
	}

	// construct command slice
	commandWithArguments := make([]string, 0, len(s.Command)+len(jc.Arguments))
	commandWithArguments = append(commandWithArguments, s.Command...)
	commandWithArguments = append(commandWithArguments, jc.Arguments...)

	// configure command
	cmd := exec.CommandContext(ctx, commandWithArguments[0], commandWithArguments[1:]...)

	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr
	cmd.Dir = r.WorkingDirectory

	// run command
	if err := cmd.Run(); err != nil {
		return err
	}

	// if we have result file, return job's result
	resultFile := path.Join(r.WorkingDirectory, resultFilename)
	if fi, err := os.Stat(resultFile); !os.IsNotExist(err) && !fi.IsDir() {
		fileContent, err := os.ReadFile(resultFile)
		if err != nil {
			return err
		}
		j.Result = &result.Result{}
		return json.Unmarshal(fileContent, j.Result)
	}

	return nil

}

func (s *commandContext) Cleanup(ctx context.Context, jobID string, c *cluster.Cluster) error {
	// Implement cleanup if needed
	return nil
}
