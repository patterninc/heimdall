package ping

import (
	"fmt"

	"github.com/patterninc/heimdall/internal/pkg/context"
	"github.com/patterninc/heimdall/internal/pkg/object/cluster"
	"github.com/patterninc/heimdall/internal/pkg/object/job"
	"github.com/patterninc/heimdall/internal/pkg/result"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	messageFormat = `Hello, %s!`
)

type pingCommandContext struct{}

func New(_ *context.Context) (plugin.Handler, error) {

	p := &pingCommandContext{}
	return p.handler, nil

}

func (p *pingCommandContext) handler(_ *plugin.Runtime, j *job.Job, _ *cluster.Cluster) (err error) {

	j.Result, err = result.FromMessage(fmt.Sprintf(messageFormat, j.User))
	return

}
