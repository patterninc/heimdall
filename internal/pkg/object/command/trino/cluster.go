package trino

import "github.com/patterninc/heimdall/pkg/object/cluster"

type clusterContext struct {
	Endpoint string  `yaml:"endpoint"`
	Token    *string `yaml:"token,omitempty"`
}

func newClusterContext(c *cluster.Cluster) (*clusterContext, error) {
	if c.Context == nil {
		return nil, ErrEmptyClusterContext
	}

	ctx := &clusterContext{}
	if err := c.Context.Unmarshal(ctx); err != nil {
		return nil, err
	}
	return ctx, nil
}
