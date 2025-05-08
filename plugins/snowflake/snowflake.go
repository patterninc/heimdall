package main

import (
	"github.com/patterninc/heimdall/internal/pkg/context"
	"github.com/patterninc/heimdall/internal/pkg/object/command/snowflake"
	"github.com/patterninc/heimdall/pkg/plugin"
)

func New(_ *context.Context) (plugin.Handler, error) {
	return snowflake.New(nil)
}
