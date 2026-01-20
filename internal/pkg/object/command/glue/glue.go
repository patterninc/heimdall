package glue

import (
	"context"

	"github.com/patterninc/heimdall/internal/pkg/aws"
	heimdallContext "github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
)

type commandContext struct {
	CatalogID string `yaml:"catalog_id,omitempty" json:"catalog_id,omitempty"`
}

type jobContext struct {
	TableName string `yaml:"table_name,omitempty" json:"table_name,omitempty"`
}

func New(commandCtx *heimdallContext.Context) (plugin.Handler, error) {

	g := &commandContext{}

	if commandCtx != nil {
		if err := commandCtx.Unmarshal(g); err != nil {
			return nil, err
		}
	}

	return g, nil

}

// Execute implements the plugin.Handler interface
func (g *commandContext) Execute(ctx context.Context, _ *plugin.Runtime, j *job.Job, _ *cluster.Cluster) (err error) {

	// let's unmarshal job context
	jc := &jobContext{}
	if j.Context != nil {
		if err = j.Context.Unmarshal(jc); err != nil {
			return
		}
	}

	// let's get our metadata
	metadata, err := aws.GetTableMetadata(ctx, g.CatalogID, jc.TableName)
	if err != nil {
		return
	}

	// return it...
	j.Result, err = result.FromMessage(string(metadata))
	return

}

// Cleanup implements the plugin.Handler interface
func (g *commandContext) Cleanup(ctx context.Context, jobID string, c *cluster.Cluster) error {
	// TODO: Implement cleanup if needed
	return nil
}
