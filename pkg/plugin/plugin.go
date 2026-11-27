package plugin

import (
	"context"
	"io"

	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
)

type Handler interface {
	Execute(context.Context, *Runtime, *job.Job, *cluster.Cluster) error
	Cleanup(ctx context.Context, jobID string, c *cluster.Cluster) error
}

type HealthChecker interface {
	HealthCheck(ctx context.Context, c *cluster.Cluster) error
}

type ResultStreamer interface {
	StreamResult(ctx context.Context, w io.Writer, r *Runtime, j *job.Job, c *cluster.Cluster) error
}

// Resumer is implemented by plugins whose remote work outlives the agent that started it.
// Resume reattaches to the work Execute started for the job and waits for it as Execute would.
type Resumer interface {
	Resume(ctx context.Context, r *Runtime, j *job.Job, c *cluster.Cluster) error
}

// StreamResumer is the Resumer of a ResultStreamer: ResumeStream reattaches to the work StreamResult started
// for the job and streams its result to w as StreamResult would.
type StreamResumer interface {
	ResumeStream(ctx context.Context, w io.Writer, r *Runtime, j *job.Job, c *cluster.Cluster) error
}

// IsResumable reports whether the janitor can release a stale job of the handler for another agent to resume.
func IsResumable(h Handler) bool {
	switch h.(type) {
	case Resumer, StreamResumer:
		return true
	}
	return false
}
