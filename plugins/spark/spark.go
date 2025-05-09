package main

import (
	"github.com/patterninc/heimdall/internal/pkg/context"
	"github.com/patterninc/heimdall/internal/pkg/object/command/spark"
	"github.com/patterninc/heimdall/pkg/plugin"
)

// New creates a new Spark plugin handler.
func New(commandContext *context.Context) (plugin.Handler, error) {
	return spark.New(commandContext)
}
