package main

import (
	"flag"
	"fmt"

	"github.com/babourine/x/pkg/config"
	"github.com/babourine/x/pkg/process"

	"github.com/patterninc/heimdall/internal/pkg/database"
)

var (
	configFile string
)

func init() {

	flag.StringVar(&configFile, `conf`, `/etc/pattern.d/heimdall.yaml`, `config file`)
	flag.Parse()

	if len(flag.Args()) != 1 {
		process.Bail(`usage`, fmt.Errorf(`deploydb [-conf <config file>] <list file>`))
	}

}

type deploy struct {
	Database *database.Database `yaml:"database,omitempty" json:"database,omitempty"`
}

func main() {

	d := deploy{}

	if err := config.LoadYAML(configFile, &d); err != nil {
		process.Bail(`config`, err)
	}

	if err := d.Database.Deploy(flag.Arg(0)); err != nil {
		process.Bail(`deploy`, err)
	}

}
