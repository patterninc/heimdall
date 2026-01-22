package janitor

import (
	"time"

	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	defaultJobLimit = 3
)

type Janitor struct {
	Keepalive       int `yaml:"keepalive,omitempty" json:"keepalive,omitempty"`
	StaleJob        int `yaml:"stale_job,omitempty" json:"stale_job,omitempty"`
	CleanInterval   int `yaml:"clean_interval,omitempty" json:"clean_interval,omitempty"`
	db              *database.Database
	commandHandlers map[string]plugin.Handler
	clusters        cluster.Clusters
}

func (j *Janitor) Start(d *database.Database, commandHandlers map[string]plugin.Handler, clusters cluster.Clusters) error {

	// record database context
	j.db = d
	j.commandHandlers = commandHandlers
	j.clusters = clusters

	// kick off janitor worker in the background.
	go func() {
		for {
			jobsFound := j.worker()

			// if no jobs are found, sleep before checking again
			if !jobsFound {
				time.Sleep(time.Duration(j.CleanInterval) * time.Second)
			}
		}
	}()

	return nil

}
