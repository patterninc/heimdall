package janitor

import (
	"fmt"
	"os"
	"testing"

	"github.com/patterninc/heimdall/internal/pkg/database"
	jobStatus "github.com/patterninc/heimdall/pkg/object/job/status"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const testDatabaseEnvVar = `HEIMDALL_TEST_DATABASE`

// lease states of a job's active_jobs row
const (
	noRow     = `no row`
	claimed   = `claimed`
	released  = `released`  // no agent, heartbeat kept from the agent that lost it
	unclaimed = `unclaimed` // no agent and no heartbeat: never claimed
)

// TestStaleJobs runs the janitor's query, split, release and update steps until nothing is left to act on,
// inside a transaction that is rolled back. It needs a Postgres with the heimdall schema in HEIMDALL_TEST_DATABASE.
func TestStaleJobs(t *testing.T) {

	connectionString := os.Getenv(testDatabaseEnvVar)
	if connectionString == `` {
		t.Skipf("%s is not set", testDatabaseEnvVar)
	}

	sess, err := (&database.Database{ConnectionString: connectionString}).NewSession(true)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	const (
		resumable = `janitor-test-resumable`
		plain     = `janitor-test-plain`
		cluster   = `janitor-test-cluster`
	)
	for _, command := range []string{resumable, plain} {
		mustExec(t, sess, `insert into commands (command_status_id, command_id, command_name, command_version, command_plugin, username, is_sync) values (1, $1, $1, '1', 'test', 'test', false)`, command)
	}
	mustExec(t, sess, `insert into clusters (cluster_status_id, cluster_id, cluster_name, cluster_version, username) values (1, $1, $1, '1', 'test')`, cluster)

	tests := []struct {
		name         string
		command      string
		status       jobStatus.Status
		lease        string
		heartbeatAge int // seconds since the last heartbeat of a claimed or released job
		wantStatus   jobStatus.Status
		wantLease    string
	}{
		{`running resumable job stale past stale_job is released`, resumable, jobStatus.Running, claimed, 60, jobStatus.Running, released},
		{`running resumable job with a fresh heartbeat is left alone`, resumable, jobStatus.Running, claimed, 30, jobStatus.Running, claimed},
		{`released job waiting for an agent is left alone`, resumable, jobStatus.Running, released, 60, jobStatus.Running, released},
		{`released job stale past awol_timeout is failed`, resumable, jobStatus.Running, released, 4000, jobStatus.Failed, noRow},
		{`running resumable job stale past awol_timeout is failed`, resumable, jobStatus.Running, claimed, 4000, jobStatus.Failed, noRow},
		{`accepted resumable job stale past stale_job is released`, resumable, jobStatus.Accepted, claimed, 60, jobStatus.Accepted, released},
		{`accepted resumable job stale past awol_timeout is failed`, resumable, jobStatus.Accepted, claimed, 4000, jobStatus.Failed, noRow},
		{`running non-resumable job with a fresh heartbeat is left alone`, plain, jobStatus.Running, claimed, 30, jobStatus.Running, claimed},
		{`running non-resumable job stale past stale_job is failed`, plain, jobStatus.Running, claimed, 60, jobStatus.Failed, noRow},
		{`accepted non-resumable job stale past stale_job is failed`, plain, jobStatus.Accepted, claimed, 60, jobStatus.Failed, noRow},
		{`canceling job without an active_jobs row is canceled`, resumable, jobStatus.Canceling, noRow, 0, jobStatus.Canceled, noRow},
		{`canceling released job is canceled`, resumable, jobStatus.Canceling, released, 60, jobStatus.Canceled, noRow},
		{`canceling job with a fresh heartbeat is left to its agent`, resumable, jobStatus.Canceling, claimed, 30, jobStatus.Canceling, claimed},
		{`canceling resumable job stale past stale_job is canceled`, resumable, jobStatus.Canceling, claimed, 60, jobStatus.Canceled, noRow},
		{`canceling non-resumable job stale past stale_job is canceled`, plain, jobStatus.Canceling, claimed, 60, jobStatus.Canceled, noRow},
		{`new job never claimed is left alone`, resumable, jobStatus.New, unclaimed, 0, jobStatus.New, unclaimed},
	}

	jobID := func(i int) string { return fmt.Sprintf("janitor-test-%d", i) }
	for i, tt := range tests {
		systemID, err := sess.InsertRow(`insert into jobs (job_command_id, job_cluster_id, job_status_id, job_id, job_name, job_version, username, is_sync)
			select cm.system_command_id, cl.system_cluster_id, $3, $4, 'test', '1', 'test', false
			from commands cm, clusters cl where cm.command_id = $1 and cl.cluster_id = $2
			returning system_job_id`, tt.command, cluster, tt.status, jobID(i))
		if err != nil {
			t.Fatal(err)
		}
		switch tt.lease {
		case claimed:
			mustExec(t, sess, `insert into active_jobs (system_job_id, agent_name, last_heartbeat) values ($1, 'agent', extract(epoch from now())::int - $2::int)`, systemID, tt.heartbeatAge)
		case released:
			mustExec(t, sess, `insert into active_jobs (system_job_id, last_heartbeat) values ($1, extract(epoch from now())::int - $2::int)`, systemID, tt.heartbeatAge)
		case unclaimed:
			mustExec(t, sess, `insert into active_jobs (system_job_id) values ($1)`, systemID)
		}
	}

	handlers := map[string]plugin.Handler{resumable: resumer{}, plain: handler{}}
	j := &Janitor{StaleJob: 45, AWOLTimeout: 3600, commandHandlers: handlers}
	// the worker loops without sleeping while the query returns jobs, so every job it returns must be acted on
	for pass := 0; ; pass++ {
		if pass == 20 {
			t.Fatal("janitor query keeps returning jobs it already acted on")
		}
		jobs, awol, err := j.queryJobs(sess)
		if err != nil {
			t.Fatal(err)
		}
		if len(jobs) == 0 {
			break
		}
		jobs, released := j.splitResumable(jobs, awol)
		if err := j.releaseJobs(sess, released); err != nil {
			t.Fatal(err)
		}
		if err := j.updateJobs(sess, jobs); err != nil {
			t.Fatal(err)
		}
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row, err := sess.QueryRow(`select j.job_status_id, aj.system_job_id is not null, aj.agent_name is not null, aj.last_heartbeat is not null
				from jobs j left join active_jobs aj on aj.system_job_id = j.system_job_id where j.job_id = $1`, jobID(i))
			if err != nil {
				t.Fatal(err)
			}
			var status jobStatus.Status
			var hasRow, hasAgent, hasHeartbeat bool
			if err := row.Scan(&status, &hasRow, &hasAgent, &hasHeartbeat); err != nil {
				t.Fatal(err)
			}
			lease := noRow
			switch {
			case hasRow && hasAgent:
				lease = claimed
			case hasRow && hasHeartbeat:
				lease = released
			case hasRow:
				lease = unclaimed
			}
			if status != tt.wantStatus || lease != tt.wantLease {
				t.Errorf("got status %v, lease %s; want status %v, lease %s", status, lease, tt.wantStatus, tt.wantLease)
			}
		})
	}

}

func mustExec(t *testing.T, sess *database.Session, query string, args ...any) {
	t.Helper()
	if _, err := sess.Exec(query, args...); err != nil {
		t.Fatal(err)
	}
}
