// Package eksjob runs a Heimdall job as a Kubernetes custom resource on an EKS cluster.
//
// The plugin is CRD-agnostic: the command context names the custom resource (group, version,
// kind, plural), the spec defaults, the size profiles and the cancel field; the job context
// supplies the size, labels and spec. The plugin creates the resource, polls a standard status
// contract (status.phase, status.failureReason), streams the logs of the resource's pods and
// cancels the run by patching a spec field.
package eksjob

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"github.com/patterninc/heimdall/internal/pkg/eksauth"
	heimdallContext "github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/duration"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	defaultPollInterval  = 15 * time.Second
	minPollInterval      = time.Second
	defaultSizeField     = `resources`
	defaultCancelValue   = `cancel`
	defaultNamespace     = `default`
	healthSessionName    = `heimdall-health`
	cancelTimeout        = 30 * time.Second
	logDiscoveryInterval = 5 * time.Second
	logDrainTimeout      = 30 * time.Second

	// outputs published for cluster attribute templates
	outputRunName      = `run_name`
	outputRunNamespace = `run_namespace`
)

var (
	ErrMissingCRDField     = fmt.Errorf("eksjob command context requires crd_group, crd_version, crd_kind and crd_plural")
	ErrInvalidSizeProfiles = fmt.Errorf("eksjob command context size_profiles must map each size name to an object with requests and/or limits")
	ErrMissingClusterName  = fmt.Errorf("eksjob cluster context requires cluster_name")
	ErrMissingCluster      = fmt.Errorf("eksjob requires a cluster")
)

// commandContext is the command's `context`.
type commandContext struct {
	CRDGroup   string `yaml:"crd_group,omitempty" json:"crd_group,omitempty"`
	CRDVersion string `yaml:"crd_version,omitempty" json:"crd_version,omitempty"`
	CRDKind    string `yaml:"crd_kind,omitempty" json:"crd_kind,omitempty"`
	CRDPlural  string `yaml:"crd_plural,omitempty" json:"crd_plural,omitempty"`

	PollInterval duration.Duration `yaml:"poll_interval,omitempty" json:"poll_interval,omitempty"`

	// SpecDefaults is merged first into the resource spec; job spec keys replace top-level keys.
	SpecDefaults map[string]any `yaml:"spec_defaults,omitempty" json:"spec_defaults,omitempty"`
	// RequiredSpecFields lists top-level spec keys that must be set after the merge.
	RequiredSpecFields []string `yaml:"required_spec_fields,omitempty" json:"required_spec_fields,omitempty"`

	// SizeField is the spec key that receives the job size's corev1.ResourceRequirements.
	SizeField string `yaml:"size_field,omitempty" json:"size_field,omitempty"`
	// SizeProfiles maps a job size to a corev1.ResourceRequirements object. The quantities are
	// parsed when a job uses the profile, so a bad value fails that job, not Heimdall startup.
	SizeProfiles map[string]json.RawMessage `yaml:"size_profiles,omitempty" json:"size_profiles,omitempty"`

	// CancelField is the spec key set to CancelValue to cancel a run. Empty disables cancellation.
	CancelField string `yaml:"cancel_field,omitempty" json:"cancel_field,omitempty"`
	CancelValue string `yaml:"cancel_value,omitempty" json:"cancel_value,omitempty"`

	// PodSelectorLabel selects the run's pods (<label>=<resource name>) for log streaming.
	PodSelectorLabel string `yaml:"pod_selector_label,omitempty" json:"pod_selector_label,omitempty"`

	newClients clientFactory
}

// clusterContext is the cluster's `context`.
type clusterContext struct {
	RoleARN     string `yaml:"role_arn,omitempty" json:"role_arn,omitempty"`
	Region      string `yaml:"region,omitempty" json:"region,omitempty"`
	ClusterName string `yaml:"cluster_name,omitempty" json:"cluster_name,omitempty"`
	Namespace   string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
}

// jobContext is the job's `context`.
type jobContext struct {
	Size   string         `yaml:"size,omitempty" json:"size,omitempty"`
	Labels map[string]any `yaml:"labels,omitempty" json:"labels,omitempty"`
	Spec   map[string]any `yaml:"spec,omitempty" json:"spec,omitempty"`
}

type kubeClients struct {
	dynamic dynamic.Interface
	kube    kubernetes.Interface
}

// clientFactory builds Kubernetes clients for a cluster; replaced in tests. With requireActive
// it fails unless the EKS cluster is ACTIVE (health checks); jobs also run while it is UPDATING.
type clientFactory func(ctx context.Context, cc *clusterContext, sessionName string, requireActive bool) (*kubeClients, error)

// New creates an eksjob plugin handler. It only parses and checks the shape of the command
// context and makes no network calls: Heimdall calls New for every command at startup and a
// New error stops Heimdall, so value checks (sizes, quantities, spec) happen per job.
func New(commandCtx *heimdallContext.Context) (plugin.Handler, error) {

	c := &commandContext{
		PollInterval: duration.Duration(defaultPollInterval),
		SizeField:    defaultSizeField,
		newClients:   eksClients,
	}

	if commandCtx != nil {
		if err := commandCtx.Unmarshal(c); err != nil {
			return nil, fmt.Errorf("failed to unmarshal eksjob command context: %w", err)
		}
	}

	if c.CRDGroup == `` || c.CRDVersion == `` || c.CRDKind == `` || c.CRDPlural == `` {
		return nil, ErrMissingCRDField
	}

	for _, profile := range c.SizeProfiles {
		var shape map[string]json.RawMessage
		if err := json.Unmarshal(profile, &shape); err != nil || shape == nil {
			return nil, ErrInvalidSizeProfiles
		}
	}

	if time.Duration(c.PollInterval) < minPollInterval {
		c.PollInterval = duration.Duration(minPollInterval)
	}
	if c.SizeField == `` {
		c.SizeField = defaultSizeField
	}
	if c.CancelField != `` && c.CancelValue == `` {
		c.CancelValue = defaultCancelValue
	}

	return c, nil

}

// Execute creates (or attaches to) the job's custom resource and waits for it to finish.
func (c *commandContext) Execute(ctx context.Context, r *plugin.Runtime, j *job.Job, cl *cluster.Cluster) error {

	stdout, stderr := newSyncWriter(r.Stdout), newSyncWriter(r.Stderr)

	jc, err := parseJobContext(j)
	if err != nil {
		return err
	}

	cc, err := parseClusterContext(cl)
	if err != nil {
		return err
	}

	name, err := runName(j.ID)
	if err != nil {
		return err
	}

	obj, err := c.buildRun(name, cc.Namespace, j.ID, jc, stderr)
	if err != nil {
		return err
	}

	clients, err := c.newClients(ctx, cc, j.ID, false)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes clients for EKS cluster %q: %w", cc.ClusterName, err)
	}

	j.SetOutput(outputRunName, name)
	j.SetOutput(outputRunNamespace, cc.Namespace)

	return c.newRunner(clients, cc.Namespace, name, stdout, stderr).run(ctx, obj)

}

// HealthCheck assumes the role, checks the EKS cluster is ACTIVE and lists the custom resource in the namespace.
func (c *commandContext) HealthCheck(ctx context.Context, cl *cluster.Cluster) error {

	cc, err := parseClusterContext(cl)
	if err != nil {
		return err
	}

	clients, err := c.newClients(ctx, cc, healthSessionName, true)
	if err != nil {
		return err
	}

	return c.checkResourceAccess(ctx, clients, cc.Namespace)

}

// Cleanup cancels the job's run (Heimdall job kill, or a stale job found by the janitor).
// The resource is not deleted; its operator removes finished runs.
func (c *commandContext) Cleanup(ctx context.Context, jobID string, cl *cluster.Cluster) error {

	if c.CancelField == `` {
		return nil
	}

	cc, err := parseClusterContext(cl)
	if err != nil {
		return err
	}

	name, err := runName(jobID)
	if err != nil {
		return err
	}

	clients, err := c.newClients(ctx, cc, jobID, false)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes clients for EKS cluster %q: %w", cc.ClusterName, err)
	}

	_, err = c.cancelRun(ctx, clients.dynamic, cc.Namespace, name)
	return err

}

func (c *commandContext) gvr() schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: c.CRDGroup, Version: c.CRDVersion, Resource: c.CRDPlural}
}

func parseJobContext(j *job.Job) (*jobContext, error) {
	jc := &jobContext{}
	if j != nil && j.Context != nil {
		if err := j.Context.Unmarshal(jc); err != nil {
			return nil, fmt.Errorf("failed to unmarshal eksjob job context: %w", err)
		}
	}
	return jc, nil
}

func parseClusterContext(cl *cluster.Cluster) (*clusterContext, error) {

	if cl == nil {
		return nil, ErrMissingCluster
	}

	cc := &clusterContext{}
	if cl.Context != nil {
		if err := cl.Context.Unmarshal(cc); err != nil {
			return nil, fmt.Errorf("failed to unmarshal eksjob cluster context: %w", err)
		}
	}

	if cc.ClusterName == `` {
		return nil, ErrMissingClusterName
	}
	if cc.Namespace == `` {
		cc.Namespace = defaultNamespace
	}

	return cc, nil

}

// eksClients assumes the cluster's role and builds clients authenticated with an EKS token.
func eksClients(ctx context.Context, cc *clusterContext, sessionName string, requireActive bool) (*kubeClients, error) {

	restConfig, err := eksauth.NewRESTConfig(ctx, eksauth.Config{
		RoleARN:     cc.RoleARN,
		Region:      cc.Region,
		ClusterName: cc.ClusterName,
		SessionName: sessionName,
	}, requireActive)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &kubeClients{dynamic: dynamicClient, kube: kubeClient}, nil

}

// syncWriter serialises writes from the log streams and the status loop.
type syncWriter struct {
	mu sync.Mutex
	w  io.Writer
}

func newSyncWriter(f *os.File) *syncWriter {
	if f == nil {
		return &syncWriter{w: io.Discard}
	}
	return &syncWriter{w: f}
}

func (s *syncWriter) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.w.Write(p)
}
