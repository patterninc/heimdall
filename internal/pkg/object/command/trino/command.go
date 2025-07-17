package trino

import (
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
	"github.com/patterninc/heimdall/pkg/result/column"
)

type commandContext struct{}

func newCommandContext(_ *context.Context) (*commandContext, error) {
	return &commandContext{}, nil
}

func (ctx *commandContext) handler(r *plugin.Runtime, j *job.Job, c *cluster.Cluster) (err error) {
	r.Stdout.WriteString("Loading cluster context\n")
	clusterCtx, err := newClusterContext(c)
	if err != nil {
		return err
	}

	r.Stdout.WriteString("Loading job context\n")
	jobCtx, err := newJobContext(j)
	if err != nil {
		r.Stderr.WriteString(err.Error())
		return err
	}

	r.Stdout.WriteString("Making trino request\n")
	trinoClient := newTrinoClient(clusterCtx.Endpoint, j.User, clusterCtx.Token, r)
	output, err := trinoClient.Query(jobCtx.Query)
	if err != nil {
		r.Stderr.WriteString(err.Error())
		return err
	}

	result := result.Result{
		Columns: make([]*column.Column, len(output.Columns)),
		Data:    output.Data,
	}
	for i, c := range output.Columns {
		result.Columns[i] = &column.Column{Name: c.Name, Type: column.Type(c.Type)}
	}
	j.Result = &result

	return nil
}
