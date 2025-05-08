package main

import (
	"github.com/patterninc/heimdall/internal/pkg/context"
	"github.com/patterninc/heimdall/internal/pkg/object/command/dynamo"
	"github.com/patterninc/heimdall/pkg/plugin"
)

// New creates a new dynamo plugin handler.
func New(_ *context.Context) (plugin.Handler, error) {
	return dynamo.New(nil)
}
