package main

import (
	"github.com/patterninc/heimdall/internal/pkg/object/command/spark"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/plugin"
)

// New creates a new Spark plugin handler.
func New(commandContext *context.Context) (*plugin.Handlers, error) {
	handler, err := spark.New(commandContext)
	if err != nil {
		return nil, err
	}
	return &plugin.Handlers{
		Handler:        handler,
		CleanupHandler: nil,
	}, nil
}
