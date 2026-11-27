package janitor

import (
	"context"
	"testing"

	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	jobStatus "github.com/patterninc/heimdall/pkg/object/job/status"
	"github.com/patterninc/heimdall/pkg/plugin"
)

type handler struct{}

func (handler) Execute(context.Context, *plugin.Runtime, *job.Job, *cluster.Cluster) error {
	return nil
}
func (handler) Cleanup(context.Context, string, *cluster.Cluster) error { return nil }

type resumer struct{ handler }

func (resumer) Resume(context.Context, *plugin.Runtime, *job.Job, *cluster.Cluster) error { return nil }

func TestSplitResumable(t *testing.T) {

	j := &Janitor{commandHandlers: map[string]plugin.Handler{
		`plain`:     handler{},
		`resumable`: resumer{},
	}}

	newJob := func(id int64, command string, status jobStatus.Status) *job.Job {
		jb := &job.Job{CommandID: command, Status: status}
		jb.SystemID = id
		return jb
	}

	rest, resumable := j.splitResumable([]*job.Job{
		newJob(1, `resumable`, jobStatus.Running),   // released
		newJob(2, `resumable`, jobStatus.Accepted),  // released, run again from the start
		newJob(3, `resumable`, jobStatus.Canceling), // canceled as today
		newJob(4, `plain`, jobStatus.Running),       // no Resumer, failed as today
		newJob(5, `unknown`, jobStatus.Running),     // command no longer loaded
		newJob(6, `resumable`, jobStatus.Running),   // past awol_timeout, failed
	}, map[int64]bool{6: true})

	if got := ids(resumable); len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("resumable: got %v, want jobs 1 and 2", got)
	}
	if got := ids(rest); len(got) != 4 || got[0] != 3 || got[3] != 6 {
		t.Fatalf("rest: got %v, want jobs 3 to 6", got)
	}

}

func ids(jobs []*job.Job) []int64 {
	out := make([]int64, 0, len(jobs))
	for _, jb := range jobs {
		out = append(out, jb.SystemID)
	}
	return out
}
