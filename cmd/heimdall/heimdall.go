package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/babourine/x/pkg/config"
	"github.com/babourine/x/pkg/process"

	"github.com/patterninc/heimdall/internal/pkg/heimdall"
	"github.com/patterninc/heimdall/internal/pkg/janitor"
	"github.com/patterninc/heimdall/internal/pkg/server"
)

const (
	defaultAddress           = `0.0.0.0:9090`
	defaultBuild             = `0.0.0`
	defaultReadTimeout       = 15 // seconds
	defaulWriteTimeout       = 15 // seconds
	defaultIdleTimeout       = 30 // seconds
	defaultReadHeaderTimeout = 2  // seconds
	defaultJanitorKeepalive  = 5  // seconds
	defaultStaleJob          = 45 // seconds
	defaultCleanInterval     = 60 // seconds
)

var (
	configFile string
	Build      string
)

func init() {

	flag.StringVar(&configFile, `conf`, `/etc/heimdall/heimdall.yaml`, `config file`)
	flag.Parse()

}

func main() {

	// setup defaults before we unmarshal config
	h := heimdall.Heimdall{
		Server: &server.Server{
			Address:           defaultAddress,
			ReadTimeout:       defaultReadTimeout,
			WriteTimeout:      defaulWriteTimeout,
			IdleTimeout:       defaultIdleTimeout,
			ReadHeaderTimeout: defaultReadHeaderTimeout,
		},
		Janitor: &janitor.Janitor{
			Keepalive:     defaultJanitorKeepalive,
			StaleJob:      defaultStaleJob,
			CleanInterval: defaultCleanInterval,
		},
	}

	// load config file
	if err := config.LoadYAML(configFile, &h); err != nil {
		process.Bail(`config`, err)
	}

	// setup version
	if Build != `` {
		h.Version = Build
	} else {
		h.Version = defaultBuild
	}

	// create context that cancels on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// initialize internal components (janitor, pool, etc.)
	if err := h.Init(ctx); err != nil {
		process.Bail(`init`, err)
	}

	// start proxy and return when server exits or signal received
	if err := h.Start(ctx); err != nil {
		process.Bail(`server`, err)
	}

}
