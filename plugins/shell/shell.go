package main

import (
	"github.com/patterninc/heimdall/internal/pkg/object/command/shell"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/plugin"
)

func New(commandContext *context.Context) (plugin.Handler, error) {
	return shell.New(commandContext)
}
