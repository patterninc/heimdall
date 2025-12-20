package plugin

import (
	"context"

	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
)

type Handlers struct {
	Handler        Handler
	CleanupHandler CleanupHandler
}

type Handler func(context.Context, *Runtime, *job.Job, *cluster.Cluster) error

type CleanupHandler func(ctx context.Context, j *job.Job, c *cluster.Cluster) error
