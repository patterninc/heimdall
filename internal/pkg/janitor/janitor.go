package janitor

import (
	"fmt"
	"time"

	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	defaultJobLimit = 25
)

type Janitor struct {
	Keepalive       int `yaml:"keepalive,omitempty" json:"keepalive,omitempty"`
	StaleJob        int `yaml:"stale_job,omitempty" json:"stale_job,omitempty"`
	db              *database.Database
	commandHandlers map[string]*plugin.Handlers
	clusters        cluster.Clusters
}

func (j *Janitor) Start(d *database.Database, commandHandlers map[string]*plugin.Handlers, clusters cluster.Clusters) error {

	// record database context
	j.db = d
	j.commandHandlers = commandHandlers
	j.clusters = clusters

	// let's run jobs cleanup once before we start it as a go routine
	if err := j.cleanupStaleJobs(); err != nil {
		return err
	}

	// start cleanup loops
	runCleanupLoop := func(cleanupFn func() error) {
		for {
			if err := cleanupFn(); err != nil {
				fmt.Printf("Janitor error: %v\n", err)
			}
			time.Sleep(60 * time.Second)
		}
	}

	go runCleanupLoop(j.cleanupStaleJobs)
	go runCleanupLoop(j.cleanupCancellingJobs)

	return nil

}
