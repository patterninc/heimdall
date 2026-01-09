package main

import (
	"github.com/patterninc/heimdall/internal/pkg/object/command/shell"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/plugin"
)

func New(commandContext *context.Context) (*plugin.Handlers, error) {
	handler, err := shell.New(commandContext)
	if err != nil {
		return nil, err
	}
	return &plugin.Handlers{
		Handler:        handler,
		CleanupHandler: nil, // exec.CommandContext handles cancellation automatically
	}, nil
}
