package ping

import (
	ct "context"
	"fmt"

	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
)

const (
	messageFormat = `Hello, %s!`
)

type pingCommandContext struct{}

func New(_ *context.Context) (plugin.Handler, error) {

	p := &pingCommandContext{}
	return p.handler, nil

}

func (p *pingCommandContext) handler(ct ct.Context, _ *plugin.Runtime, j *job.Job, _ *cluster.Cluster) (err error) {

	j.Result, err = result.FromMessage(fmt.Sprintf(messageFormat, j.User))
	return

}
