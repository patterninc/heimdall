package janitor

import (
	"time"

	"github.com/hladush/go-telemetry/pkg/telemetry"
	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	defaultJobLimit = 3
)

var (
	startMethod = telemetry.NewMethod("Start", "janitor")
)

type Janitor struct {
	Keepalive                int `yaml:"keepalive,omitempty" json:"keepalive,omitempty"`
	StaleJob                 int `yaml:"stale_job,omitempty" json:"stale_job,omitempty"`
	FinishedJobRetentionDays int `yaml:"finished_job_retention_days,omitempty" json:"finished_job_retention_days,omitempty"`
	db                       *database.Database
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

			if err := j.cleanupStaleJobs(); err != nil {
				startMethod.LogAndCountError(err, "cleanup_stale_jobs")
			}

			if err := j.cleanupFinishedJobs(); err != nil {
				startMethod.LogAndCountError(err, "cleanup_finished_jobs")
			}
			time.Sleep(60 * time.Second)

		}
	}()

	return nil

}
