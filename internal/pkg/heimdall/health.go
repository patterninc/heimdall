package heimdall

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/status"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	healthCheckTimeout     = 30 * time.Second
	healthCheckConcurrency = 10
	healthStatusOK         = `ok`
	healthStatusError      = `error`
	healthStatusUnchecked  = `unchecked`
)

type clusterProbe struct {
	cluster    *cluster.Cluster
	handler    plugin.Handler
	pluginName string
}

type healthCheckResult struct {
	ClusterID   string `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
	Plugin      string `json:"plugin"`
	Status      string `json:"status"`
	LatencyMs   int64  `json:"latency_ms"`
	Error       string `json:"error,omitempty"`
}

type healthChecksResponse struct {
	Healthy bool                `json:"healthy"`
	Checks  []healthCheckResult `json:"checks"`
}

func (h *Heimdall) healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), healthCheckTimeout)
	defer cancel()

	probes := h.resolveClusterProbes()
	results := h.runHealthChecks(ctx, probes)

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

func (h *Heimdall) resolveClusterProbes() []*clusterProbe {
	var probes []*clusterProbe
	for _, cl := range h.Clusters {
		if cl.Status != status.Active || !cl.HealthCheck {
			continue
		}
		for _, cmd := range h.Commands {
			if cmd.Status != status.Active {
				continue
			}
			if cl.Tags.Contains(cmd.ClusterTags) {
				probes = append(probes, &clusterProbe{
					cluster:    cl,
					handler:    h.commandHandlers[cmd.ID],
					pluginName: cmd.Plugin,
				})
				break
			}
		}
	}
	return probes
}

func (h *Heimdall) runHealthChecks(ctx context.Context, probes []*clusterProbe) []healthCheckResult {
	results := make([]healthCheckResult, len(probes))
	sem := make(chan struct{}, healthCheckConcurrency)
	var wg sync.WaitGroup
	for i, probe := range probes {
		wg.Add(1)
		go func(i int, probe *clusterProbe) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			results[i] = h.checkCluster(ctx, probe)
		}(i, probe)
	}
	wg.Wait()
	return results
}

func (h *Heimdall) checkCluster(ctx context.Context, probe *clusterProbe) healthCheckResult {
	start := time.Now()
	res := healthCheckResult{
		ClusterID:   probe.cluster.ID,
		ClusterName: probe.cluster.Name,
		Plugin:      probe.pluginName,
	}

	hc, ok := probe.handler.(plugin.HealthChecker)
	if !ok {
		res.Status = healthStatusUnchecked
		res.LatencyMs = time.Since(start).Milliseconds()
		return res
	}

	err := hc.HealthCheck(ctx, probe.cluster)
	res.LatencyMs = time.Since(start).Milliseconds()
	if err != nil {
		res.Status = healthStatusError
		res.Error = err.Error()
	} else {
		res.Status = healthStatusOK
	}
	return res
}
