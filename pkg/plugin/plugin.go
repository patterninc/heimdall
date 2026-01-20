package plugin

import (
	"context"

	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
)

type Handler interface {
	Execute(context.Context, *Runtime, *job.Job, *cluster.Cluster) error
	Cleanup(ctx context.Context, jobID string, c *cluster.Cluster) error
}
