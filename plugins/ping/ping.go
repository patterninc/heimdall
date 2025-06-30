package main

import (
	"github.com/patterninc/heimdall/internal/pkg/object/command/ping"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/plugin"
)

func New(_ *context.Context) (plugin.Handler, error) {
	return ping.New(nil)
}
