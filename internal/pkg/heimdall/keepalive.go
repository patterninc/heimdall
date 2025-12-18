package heimdall

import (
	_ "embed"
	"time"

	"github.com/hladush/go-telemetry/pkg/telemetry"
)

const (
	defaultJanitorKeepalive = 5 // seconds
)

var (
	jobKeepaliveMethod = telemetry.NewMethod("db_connection", "job_keepalive")
)

//go:embed queries/job/active_keepalive.sql
var queryActiveJobKeepalive string

func (h *Heimdall) jobKeepalive(done <-chan struct{}, jobID int64, agentName string) {

	keepaliveSeconds := defaultJanitorKeepalive
	if h.Janitor != nil && h.Janitor.Keepalive > 0 {
		keepaliveSeconds = h.Janitor.Keepalive
	}

	ticker := time.NewTicker(time.Duration(keepaliveSeconds) * time.Second)
	defer ticker.Stop()

	// Track DB connection for job keepalive
	defer jobKeepaliveMethod.RecordLatency(time.Now())
	jobKeepaliveMethod.CountRequest()

	// set the db session
	sess, err := h.Database.NewSession(false)
	if err != nil {
		jobKeepaliveMethod.LogAndCountError(err, "new_session")
		sess = nil
	}
	defer sess.Close()

	for {
		select {
		case <-ticker.C:
			// let's update job's keepalive timestamp
			if sess != nil {
				// make the best effort to keep the job "alive" (eat the error!)
				sess.Exec(queryActiveJobKeepalive, jobID, agentName)
			}
		case _, stillOpen := <-done:
			if !stillOpen {
				jobKeepaliveMethod.CountSuccess()
				return
			}
		}
	}

}
