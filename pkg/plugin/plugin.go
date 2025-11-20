package plugin

import (
	"context"

	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
)

type Handler func(context.Context, *Runtime, *job.Job, *cluster.Cluster) error
