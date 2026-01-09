package main

import (
	"github.com/patterninc/heimdall/internal/pkg/object/command/dynamo"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/plugin"
)

// New creates a new dynamo plugin handler.
func New(_ *context.Context) (*plugin.Handlers, error) {
	handler, err := dynamo.New(nil)
	if err != nil {
		return nil, err
	}
	return &plugin.Handlers{
		Handler:        handler,
		CleanupHandler: nil,
	}, nil
}
