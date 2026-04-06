package heimdall

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/command"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/object/status"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	healthCheckTimeout = 30 * time.Second
	healthCheckUser    = `heimdall-health`
	healthStatusOK     = `ok`
	healthStatusError  = `error`
)

type healthCheckResult struct {
	CommandID string `json:"command_id"`
	ClusterID string `json:"cluster_id"`
	Status    string `json:"status"`
	LatencyMs int64  `json:"latency_ms"`
	Error     string `json:"error,omitempty"`
}

type healthChecksResponse struct {
	Healthy bool                `json:"healthy"`
	Checks  []healthCheckResult `json:"checks"`
}

type healthPair struct {
	cmd     *command.Command
	cluster *cluster.Cluster
	handler plugin.Handler
}

func (h *Heimdall) healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), healthCheckTimeout)
	defer cancel()

	results := h.runHealthChecks(ctx, h.resolveHealthPairs(ctx))

	healthy := true
	for _, res := range results {
		if res.Status == healthStatusError {
			healthy = false
			break
		}
	}

	resp := healthChecksResponse{Healthy: healthy, Checks: results}
	data, _ := json.Marshal(resp)

	w.Header().Set(contentTypeKey, contentTypeJSON)
	if healthy {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	w.Write(data)
}

func (h *Heimdall) resolveHealthPairs(ctx context.Context) []*healthPair {
	var pairs []*healthPair
	for _, cmd := range h.Commands {
		if !cmd.HealthCheck {
			continue
		}
		// check DB for command status to avoid unnecessary health checks for inactive commands
		dbCmd, err := h.getCommandStatus(ctx, cmd)
		if err != nil {
			continue
		}
		if dbCmd.(*command.Command).Status != status.Active {
			continue
		}
		// find active clusters matching command's cluster tags
		for _, cluster := range h.Clusters {
			if cluster.Status != status.Active {
				continue
			}
			if cluster.Tags.Contains(cmd.ClusterTags) {
				pairs = append(pairs, &healthPair{cmd, cluster, h.commandHandlers[cmd.ID]})
			}
		}
	}
	return pairs
}

func (h *Heimdall) runHealthChecks(ctx context.Context, pairs []*healthPair) []healthCheckResult {
	results := make([]healthCheckResult, len(pairs))
	var wg sync.WaitGroup
	for i, pair := range pairs {
		wg.Add(1)
		go func(i int, pair *healthPair) {
			defer wg.Done()
			results[i] = h.checkPair(ctx, pair)
		}(i, pair)
	}
	wg.Wait()
	return results
}

func (h *Heimdall) checkPair(ctx context.Context, pair *healthPair) healthCheckResult {
	start := time.Now()
	res := healthCheckResult{CommandID: pair.cmd.ID, ClusterID: pair.cluster.ID}

	var err error
	if hc, ok := pair.handler.(plugin.HealthChecker); ok {
		err = hc.HealthCheck(ctx, pair.cluster)
	} else {
		err = h.pluginProbe(ctx, pair.cluster, pair.handler)
	}

	res.LatencyMs = time.Since(start).Milliseconds()
	if err != nil {
		res.Status = healthStatusError
		res.Error = err.Error()
	} else {
		res.Status = healthStatusOK
	}
	return res
}

func (h *Heimdall) pluginProbe(ctx context.Context, cl *cluster.Cluster, handler plugin.Handler) error {
	j := &job.Job{}
	j.ID = uuid.NewString()
	j.User = healthCheckUser

	tmpDir, err := os.MkdirTemp("", "heimdall-health-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	runtime := &plugin.Runtime{
		WorkingDirectory: tmpDir,
		ResultDirectory:  tmpDir + separator + "result",
		Version:          h.Version,
		UserAgent:        fmt.Sprintf(formatUserAgent, h.Version),
	}

	if err := runtime.Set(); err != nil {
		return err
	}
	defer runtime.Stdout.Close()
	defer runtime.Stderr.Close()

	return handler.Execute(ctx, runtime, j, cl)
}
