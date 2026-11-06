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
		newJob(2, `resumable`, jobStatus.Accepted),  // never started, failed as today
		newJob(3, `resumable`, jobStatus.Canceling), // canceled as today
		newJob(4, `plain`, jobStatus.Running),       // no Resumer, failed as today
		newJob(5, `unknown`, jobStatus.Running),     // command no longer loaded
	})

	if len(resumable) != 1 || resumable[0].SystemID != 1 {
		t.Fatalf("resumable: got %v, want job 1 only", ids(resumable))
	}
	if got := ids(rest); len(got) != 4 || got[0] != 2 || got[3] != 5 {
		t.Fatalf("rest: got %v, want jobs 2 to 5", got)
	}

}

func ids(jobs []*job.Job) []int64 {
	out := make([]int64, 0, len(jobs))
	for _, jb := range jobs {
		out = append(out, jb.SystemID)
	}
	return out
}
