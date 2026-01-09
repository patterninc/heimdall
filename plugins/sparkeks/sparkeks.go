package main

import (
	sparkeks "github.com/patterninc/heimdall/internal/pkg/object/command/sparkeks"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/plugin"
)

// New creates a new Spark EKS plugin handler.
func New(commandContext *context.Context) (*plugin.Handlers, error) {
	handler, err := sparkeks.New(commandContext)
	if err != nil {
		return nil, err
	}
	return &plugin.Handlers{
		Handler:        handler,
		CleanupHandler: nil,
	}, nil
}
