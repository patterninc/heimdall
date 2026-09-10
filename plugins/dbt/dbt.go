package main

import (
	"github.com/patterninc/heimdall/internal/pkg/object/command/dbt"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/plugin"
)

func New(ctx *context.Context) (plugin.Handler, error) {
	return dbt.New(ctx)
}
