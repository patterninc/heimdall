package main

import (
	"flag"

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
)

var (
	configFile string
	Build      string
)

func init() {

	flag.StringVar(&configFile, `conf`, `/Users/ivanhladush/git/heimdall/configs/local.yaml`, `config file`)
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
			Keepalive: defaultJanitorKeepalive,
			StaleJob:  defaultStaleJob,
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

	// start proxy
	if err := h.Start(); err != nil {
		process.Bail(`server`, err)
	}

}
