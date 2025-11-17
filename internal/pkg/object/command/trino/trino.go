package trino

import (
	ct "context"
	"fmt"
	"log"
	"time"

	"github.com/hladush/go-telemetry/pkg/telemetry"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	trinoPlugin         = `trino`
	finishedState       = `FINISHED`
	defaultPollInterval = 150 // ms
)

var (
	canBeExecutedMethod = telemetry.NewMethod("canBeExecuted", trinoPlugin)
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

func (t *commandContext) handler(ct ct.Context, r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {

	// get job context
	jobCtx := &jobContext{}
	if j.Context != nil {
		if err := j.Context.Unmarshal(jobCtx); err != nil {
			return err
		}
	}
	jobCtx.Query = normalizeTrinoQuery(jobCtx.Query)

	if !canQueryBeExecuted(jobCtx.Query, j.User, j.ID, c) {
		log.Printf("user %s is not allowed to run the query", j.User)
		// this code will be enabled in prod after some testing
	}
	// let's submit our query to trino
	req, err := newRequest(r, j, c, jobCtx)
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

func canQueryBeExecuted(query, user, id string, c *cluster.Cluster) bool {
	defer canBeExecutedMethod.RecordLatency(time.Now())
	if query == `` {
		return false
	}

	for _, rbac := range c.RBACs {
		allowed, err := rbac.HasAccess(user, query)
		if err != nil {
			err = fmt.Errorf("error checking an access for user %s and jobID %s cause: %w", user, id, err)
			canBeExecutedMethod.LogAndCountError(err)
			return false
		}
		if !allowed {
			err = fmt.Errorf("access denied for user %s to run the query in jobID %s", user, id)
			canBeExecutedMethod.LogAndCountError(err)
			return false
		}
	}

	canBeExecutedMethod.CountSuccess()
	return true
}
