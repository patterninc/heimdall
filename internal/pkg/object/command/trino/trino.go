package trino

import (
	"fmt"
	"time"

	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	finishedState       = `FINISHED`
	defaultPollInterval = 150 // ms
)

type commandContext struct {
	PollInterval int `yaml:"poll_interval,omitempty" json:"poll_interval,omitempty"`
}

type clusterContext struct {
	Endpoint string `yaml:"endpoint" json:"endpoint"`
	Catalog  string `yaml:"catalog,omitempty" json:"catalog,omitempty"`
}

type jobContext struct {
	Query string `yaml:"query" json:"query"`
}

// New creates a new trino plugin handler
func New(ctx *context.Context) (plugin.Handler, error) {

	t := &commandContext{
		PollInterval: defaultPollInterval,
	}

	if ctx != nil {
		if err := ctx.Unmarshal(t); err != nil {
			return nil, err
		}
	}

	return t.handler, nil

}

func (t *commandContext) handler(r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {

	// let's submit our query to trino
	req, err := newRequest(r, j, c)
	if err != nil {
		return err
	}

	// now let's keep pooling until we get the full result...
	for req.nextUri != `` {
		time.Sleep(time.Duration(t.PollInterval) * time.Millisecond)
		if err := req.poll(); err != nil {
			return err
		}
	}

	// return query result
	j.Result = req.result

	// did the query succeed?
	if s := req.state; s != finishedState {
		return fmt.Errorf("query finished with unexpected status: %s", s)
	}

	return nil

}
