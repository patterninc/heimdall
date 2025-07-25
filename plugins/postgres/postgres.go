package main

import (
	"github.com/patterninc/heimdall/internal/pkg/object/command/postgres"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/plugin"
)

// New creates a new instance of the postgres plugin.
func New(c *context.Context) (plugin.Handler, error) {
	return postgres.New(c)
}
