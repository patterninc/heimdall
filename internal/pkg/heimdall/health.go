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
	"github.com/gorilla/mux"

	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/command"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/object/status"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	healthCheckTimeout     = 30 * time.Second
	healthCheckUser        = `heimdall-health`
	healthStatusOK         = `ok`
	healthStatusError      = `error`
	healthCheckConcurrency = 10
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
	h.writeHealthResponse(w, ctx, nil)
}

func (h *Heimdall) commandHealthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), healthCheckTimeout)
	defer cancel()

	commandID := mux.Vars(r)[idKey]
	cmd, found := h.Commands[commandID]
	if !found {
		writeAPIError(w, ErrUnknownCommandID, nil)
		return
	}
	if !cmd.HealthCheck {
		writeAPIError(w, fmt.Errorf("command %s has not opted into health checks", commandID), nil)
		return
	}

	h.writeHealthResponse(w, ctx, &commandID)
}

func (h *Heimdall) writeHealthResponse(w http.ResponseWriter, ctx context.Context, commandID *string) {
	pairs, err := h.resolveHealthPairs(ctx, commandID)
	if err != nil {
		writeAPIError(w, fmt.Errorf("error resolving health check pairs: %w", err), nil)
		return
	}
	results := h.runHealthChecks(ctx, pairs)

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

func (h *Heimdall) resolveHealthPairs(ctx context.Context, commandID *string) ([]*healthPair, error) {
	var pairs []*healthPair
	if commandID != nil {
		cmd, found := h.Commands[*commandID]
		if !found {
			return pairs, nil
		}
		return h.resolveHealthPairsForCommand(ctx, cmd)
	}
	for _, cmd := range h.Commands {
		cmdPairs, err := h.resolveHealthPairsForCommand(ctx, cmd)
		if err != nil {
			return pairs, err
		}
		pairs = append(pairs, cmdPairs...)
	}
	return pairs, nil
}

func (h *Heimdall) resolveHealthPairsForCommand(ctx context.Context, cmd *command.Command) ([]*healthPair, error) {
	var pairs []*healthPair
	if cmd == nil {
		return pairs, nil
	}
	if !cmd.HealthCheck {
		return pairs, nil
	}
	// check DB for command status to avoid unnecessary health checks for inactive commands
	dbCmd, err := h.getCommandStatus(ctx, cmd)
	if err != nil {
		return pairs, err
	}
	if dbCmd.(*command.Command).Status != status.Active {
		return pairs, nil
	}
	for _, cl := range h.Clusters {
		if cl.Status != status.Active {
			continue
		}
		if cl.Tags.Contains(cmd.ClusterTags) {
			pairs = append(pairs, &healthPair{cmd, cl, h.commandHandlers[cmd.ID]})
		}
	}
	return pairs, nil
}

func (h *Heimdall) runHealthChecks(ctx context.Context, pairs []*healthPair) []healthCheckResult {
	results := make([]healthCheckResult, len(pairs))
	sem := make(chan struct{}, healthCheckConcurrency)
	var wg sync.WaitGroup
	for i, pair := range pairs {
		wg.Add(1)
		go func(i int, pair *healthPair) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
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
