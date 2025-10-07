package main

import (
	"github.com/patterninc/heimdall/internal/pkg/object/command/clickhouse"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/plugin"
)

// New creates a new clickhouse plugin handler
func New(commandContext *context.Context) (plugin.Handler, error) {
	return clickhouse.New(commandContext)
}
