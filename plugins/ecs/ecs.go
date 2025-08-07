package main

import (
	"github.com/patterninc/heimdall/internal/pkg/object/command/ecs"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/plugin"
)

// New creates a new ECS plugin handler.
func New(commandContext *context.Context) (plugin.Handler, error) {
	return ecs.New(commandContext)
}
