package janitor

import (
	"time"

	"github.com/hladush/go-telemetry/pkg/telemetry"
	"github.com/patterninc/heimdall/internal/pkg/database"
)

var (
	startMethod = telemetry.NewMethod("Start", "Janitor")
)

type Janitor struct {
	Keepalive                int `yaml:"keepalive,omitempty" json:"keepalive,omitempty"`
	StaleJob                 int `yaml:"stale_job,omitempty" json:"stale_job,omitempty"`
	FinishedJobRetentionDays int `yaml:"finished_job_retention_days,omitempty" json:"finished_job_retention_days,omitempty"`
	db                       *database.Database
}

func (j *Janitor) Start(d *database.Database) error {

	// record database context
	j.db = d

	// let's run jobs cleanup once before we start it as a go routine
	if err := j.cleanupStaleJobs(); err != nil {
		return err
	}

	// start cleanup loop
	go func() {

		for {

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
