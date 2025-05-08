package plugin

import (
	"github.com/patterninc/heimdall/internal/pkg/object/cluster"
	"github.com/patterninc/heimdall/internal/pkg/object/job"
)

type Handler func(*Runtime, *job.Job, *cluster.Cluster) error
