package ping

import (
	"context"
	"fmt"

	heimdallContext "github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
)

const (
	messageFormat = `Hello, %s!`
)

type commandContext struct{}

func New(_ *heimdallContext.Context) (plugin.Handler, error) {

	p := &commandContext{}
	return p, nil

}

// Execute implements the plugin.Handler interface
func (p *commandContext) Execute(ctx context.Context, _ *plugin.Runtime, j *job.Job, _ *cluster.Cluster) (err error) {

	j.Result, err = result.FromMessage(fmt.Sprintf(messageFormat, j.User))
	return

}

// Cleanup implements the plugin.Handler interface
func (p *commandContext) Cleanup(ctx context.Context, jobID string, c *cluster.Cluster) error {
	// TODO: Implement cleanup if needed
	return nil
}
