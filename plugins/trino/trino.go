package main

import (
	"github.com/patterninc/heimdall/internal/pkg/object/command/trino"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/plugin"
)

// New creates a new trino plugin handler
func New(commandContext *context.Context) (plugin.Handler, error) {
	return trino.New(commandContext)
}
