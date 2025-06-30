package shell

import (
	"encoding/json"
	"os"
	"os/exec"
	"path"

	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
)

const (
	resultFilename  = `result.json`
	contextFilename = `context.json`
)

type shellCommandContext struct {
	Command []string `yaml:"command,omitempty" json:"command,omitempty"`
}

type shellJobContext struct {
	Arguments []string `yaml:"arguments,omitempty" json:"arguments,omitempty"`
}

type runtimeContext struct {
	Job     *context.Context     `yaml:"job,omitempty" json:"job,omitempty"`
	Command *shellCommandContext `yaml:"command,omitempty" json:"command,omitempty"`
	Cluster *context.Context     `yaml:"cluster,omitempty" json:"cluster,omitempty"`
}

func New(commandContext *context.Context) (plugin.Handler, error) {

	s := &shellCommandContext{}

	if commandContext != nil {
		if err := commandContext.Unmarshal(s); err != nil {
			return nil, err
		}
	}

	return s.handler, nil

}

func (s *shellCommandContext) handler(r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {

	// let's unmarshal job context
	jc := &shellJobContext{}
	if j.Context != nil {
		if err := j.Context.Unmarshal(jc); err != nil {
			return err
		}
	}

	// let's write job context as context.json to the job's working directory
	// this will enable us to pass any context into shell command and...
	// parse those keys in the shell script using `jq` utility
	rc := &runtimeContext{
		Job:     j.Context,
		Command: s,
		Cluster: c.Context,
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
	cmd := exec.Command(commandWithArguments[0], commandWithArguments[1:]...)

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
