package trino

import (
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/plugin"
)

// New creates a new trino plugin handler
func New(ctx *context.Context) (plugin.Handler, error) {
	commandContext, err := newCommandContext(ctx)
	if err != nil {
		return nil, err
	}
	return commandContext.handler, nil
}
