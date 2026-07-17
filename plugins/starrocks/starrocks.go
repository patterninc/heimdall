package main

import (
	"github.com/patterninc/heimdall/internal/pkg/object/command/starrocks"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/plugin"
)

// New creates a new starrocks plugin handler
func New(commandContext *context.Context) (plugin.Handler, error) {
	return starrocks.New(commandContext)
}
