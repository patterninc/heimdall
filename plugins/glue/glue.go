package main

import (
	"github.com/patterninc/heimdall/internal/pkg/context"
	"github.com/patterninc/heimdall/internal/pkg/object/command/glue"
	"github.com/patterninc/heimdall/pkg/plugin"
)

func New(commandContext *context.Context) (plugin.Handler, error) {
	return glue.New(commandContext)
}
