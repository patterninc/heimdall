package janitor

import (
	"time"

	"github.com/hladush/go-telemetry/pkg/telemetry"
	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	defaultJobLimit   = 25
	defaultNumWorkers = 3
)

var (
	janitorStartMethod = telemetry.NewMethod("janitor", "start")
)

type Janitor struct {
	Keepalive       int `yaml:"keepalive,omitempty" json:"keepalive,omitempty"`
	StaleJob        int `yaml:"stale_job,omitempty" json:"stale_job,omitempty"`
	CleanInterval   int `yaml:"clean_interval,omitempty" json:"clean_interval,omitempty"`
	db              *database.Database
	commandHandlers map[string]*plugin.Handlers
	clusters        cluster.Clusters
}

func (j *Janitor) Start(d *database.Database, commandHandlers map[string]*plugin.Handlers, clusters cluster.Clusters) error {

	// record database context
	j.db = d
	j.commandHandlers = commandHandlers
	j.clusters = clusters

	// create channel for cleanup jobs
	jobChan := make(chan *job.Job, defaultJobLimit*2)

	// start cleanup workers
	for i := 0; i < defaultNumWorkers; i++ {
		go j.cleanupWorker(jobChan)
	}

	// send jobs to channel
	go func() {
		for {
			if err := j.queryAndSendJobs(jobChan); err != nil {
				janitorStartMethod.CountError("query_and_send_jobs")
			}
			time.Sleep(time.Duration(j.CleanInterval) * time.Second)
		}
	}()

	return nil

}
