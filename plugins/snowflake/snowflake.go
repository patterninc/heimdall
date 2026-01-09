package main

import (
	"github.com/patterninc/heimdall/internal/pkg/object/command/snowflake"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/plugin"
)

func New(_ *context.Context) (*plugin.Handlers, error) {
	handler, err := snowflake.New(nil)
	if err != nil {
		return nil, err
	}
	return &plugin.Handlers{
		Handler:        handler,
		CleanupHandler: nil,
	}, nil
}
