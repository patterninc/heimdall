package main

import (
	sparkeks "github.com/patterninc/heimdall/internal/pkg/object/command/sparkeks"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/plugin"
)

// New creates a new Spark EKS plugin handler.
func New(commandContext *context.Context) (plugin.Handler, error) {
	return sparkeks.New(commandContext)
}
