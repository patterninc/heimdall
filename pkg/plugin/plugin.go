package plugin

import (
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
)

type Handler func(*Runtime, *job.Job, *cluster.Cluster) error
