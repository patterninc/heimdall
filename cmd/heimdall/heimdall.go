package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strconv"
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
	defaultReadTimeout       = 15   // seconds
	defaulWriteTimeout       = 15   // seconds
	defaultIdleTimeout       = 30   // seconds
	defaultReadHeaderTimeout = 2    // seconds
	defaultJanitorKeepalive  = 5    // seconds
	defaultStaleJob          = 45   // seconds
	defaultAWOLTimeout       = 3600 // seconds
	defaultCleanInterval     = 60   // seconds
	defaultShutdownTimeout   = 30   // seconds
	shutdownTimeoutEnvVar    = `HEIMDALL_SHUTDOWN_TIMEOUT`
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
			AWOLTimeout:   defaultAWOLTimeout,
			CleanInterval: defaultCleanInterval,
		},
		ShutdownTimeout: defaultShutdownTimeout,
	}

	// size the drain to the container's stop timeout when the deployment passes it in
	if v, err := strconv.Atoi(os.Getenv(shutdownTimeoutEnvVar)); err == nil && v > 0 {
		h.ShutdownTimeout = v
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

	// the first SIGTERM or interrupt starts a graceful drain; a second one exits immediately
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	go func() {
		<-ctx.Done()
		stop()
	}()

	// start proxy
	if err := h.Start(ctx); err != nil {
		process.Bail(`server`, err)
	}

}
