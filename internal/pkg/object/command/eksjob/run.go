package eksjob

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

// Status contract every custom resource run by eksjob must expose.
const (
	phasePending   = `Pending`
	phaseRunning   = `Running`
	phaseSucceeded = `Succeeded`
	phaseFailed    = `Failed`

	statusReportInterval     = time.Minute
	maxConsecutivePollErrors = 10
	maxReportedEvents        = 20
)

var (
	statusPhasePath          = []string{`status`, `phase`}
	statusFailureReasonPath  = []string{`status`, `failureReason`}
	statusStartTimePath      = []string{`status`, `startTime`}
	statusCompletionTimePath = []string{`status`, `completionTime`}
)

type runner struct {
	cmd       *commandContext
	clients   *kubeClients
	namespace string
	name      string
	stdout    io.Writer
	stderr    io.Writer

	pollInterval         time.Duration
	logDiscoveryInterval time.Duration
	logDrainTimeout      time.Duration
}

func (c *commandContext) newRunner(clients *kubeClients, namespace, name string, stdout, stderr io.Writer) *runner {
	return &runner{
		cmd:                  c,
		clients:              clients,
		namespace:            namespace,
		name:                 name,
		stdout:               stdout,
		stderr:               stderr,
		pollInterval:         time.Duration(c.PollInterval),
		logDiscoveryInterval: logDiscoveryInterval,
		logDrainTimeout:      logDrainTimeout,
	}
}

func (r *runner) resource() dynamic.ResourceInterface {
	return r.clients.dynamic.Resource(r.cmd.gvr()).Namespace(r.namespace)
}

func (r *runner) ref() string {
	return fmt.Sprintf("%s %s/%s", r.cmd.CRDKind, r.namespace, r.name)
}

// run submits the resource, streams its pod logs and waits for a terminal phase.
func (r *runner) run(ctx context.Context, obj *unstructured.Unstructured) error {

	if err := r.createOrAttach(ctx, obj); err != nil {
		return err
	}

	var logs *logStreamer
	if r.cmd.PodSelectorLabel != `` {
		logs = newLogStreamer(r.clients.kube, r.namespace, fmt.Sprintf("%s=%s", r.cmd.PodSelectorLabel, r.name), r.stdout, r.stderr, r.logDiscoveryInterval)
		logs.start(ctx)
	} else {
		fmt.Fprintf(r.stderr, "eksjob: pod_selector_label is not set; pod logs will not be streamed\n")
	}

	final, err := r.waitForTerminal(ctx)

	if ctx.Err() != nil {
		// Heimdall canceled the job: stop streaming and cancel the run with a fresh context,
		// since sync jobs never reach the janitor's Cleanup.
		if logs != nil {
			logs.abort()
		}
		if r.cmd.CancelField != `` {
			cancelCtx, cancel := context.WithTimeout(context.Background(), cancelTimeout)
			defer cancel()
			if patched, cancelErr := r.cmd.cancelRun(cancelCtx, r.clients.dynamic, r.namespace, r.name); cancelErr != nil {
				fmt.Fprintf(r.stderr, "eksjob: failed to cancel %s: %v\n", r.ref(), cancelErr)
			} else if patched {
				fmt.Fprintf(r.stdout, "eksjob: requested cancellation of %s\n", r.ref())
			}
		}
		return ctx.Err()
	}

	if logs != nil {
		logs.finish(r.logDrainTimeout)
	}

	if err != nil {
		return err
	}

	r.printTimes(final)

	if phaseOf(final) == phaseSucceeded {
		fmt.Fprintf(r.stdout, "eksjob: %s succeeded\n", r.ref())
		return nil
	}

	reason, _, _ := unstructured.NestedString(final.Object, statusFailureReasonPath...)
	if reason == `` {
		reason = `no failureReason reported`
	}
	r.reportWarningEvents(ctx)
	fmt.Fprintf(r.stderr, "eksjob: %s failed: %s\n", r.ref(), reason)

	return fmt.Errorf("%s failed: %s", r.ref(), reason)

}

// createOrAttach creates the resource; if a resource with the same name already exists (a retry
// of the same Heimdall job) the plugin attaches to it instead of failing.
func (r *runner) createOrAttach(ctx context.Context, obj *unstructured.Unstructured) error {

	fmt.Fprintf(r.stdout, "eksjob: creating %s (apiVersion %s, spec fields: %s)\n", r.ref(), obj.GetAPIVersion(), strings.Join(specKeys(obj), `, `))

	_, err := r.resource().Create(ctx, obj, metav1.CreateOptions{})
	if err == nil {
		fmt.Fprintf(r.stdout, "eksjob: created %s\n", r.ref())
		return nil
	}

	if !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create %s: %w", r.ref(), err)
	}

	existing, err := r.resource().Get(ctx, r.name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("%s already exists but could not be read: %w", r.ref(), err)
	}

	phase := phaseOf(existing)
	if phase == `` {
		phase = `<none>`
	}
	fmt.Fprintf(r.stdout, "eksjob: %s already exists (phase %s); attaching to it\n", r.ref(), phase)

	return nil

}

// waitForTerminal polls the resource until status.phase is Succeeded or Failed.
func (r *runner) waitForTerminal(ctx context.Context) (*unstructured.Unstructured, error) {

	ticker := time.NewTicker(r.pollInterval)
	defer ticker.Stop()

	lastPhase, lastReport, consecutiveErrors := `-`, time.Time{}, 0

	for {

		obj, err := r.resource().Get(ctx, r.name, metav1.GetOptions{})
		switch {
		case ctx.Err() != nil:
			return nil, ctx.Err()
		case apierrors.IsNotFound(err):
			return nil, fmt.Errorf("%s was deleted before it finished", r.ref())
		case err != nil:
			consecutiveErrors++
			fmt.Fprintf(r.stderr, "eksjob: failed to get %s (%d/%d): %v\n", r.ref(), consecutiveErrors, maxConsecutivePollErrors, err)
			if consecutiveErrors >= maxConsecutivePollErrors {
				return nil, fmt.Errorf("failed to get %s: %w", r.ref(), err)
			}
		default:
			consecutiveErrors = 0
			phase := phaseOf(obj)
			if phase != lastPhase || time.Since(lastReport) >= statusReportInterval {
				shown := phase
				if shown == `` {
					shown = `<none>`
				}
				fmt.Fprintf(r.stdout, "%s - eksjob: %s phase: %s\n", time.Now().UTC().Format(time.RFC3339), r.ref(), shown)
				lastPhase, lastReport = phase, time.Now()
			}
			if phase == phaseSucceeded || phase == phaseFailed {
				return obj, nil
			}
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}

	}

}

// cancelRun merge-patches spec.<cancel_field> = cancel_value. It returns false without patching
// when the run does not exist or has already finished.
func (c *commandContext) cancelRun(ctx context.Context, client dynamic.Interface, namespace, name string) (bool, error) {

	resource := client.Resource(c.gvr()).Namespace(namespace)

	obj, err := resource.Get(ctx, name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get %s %s/%s: %w", c.CRDKind, namespace, name, err)
	}

	if phase := phaseOf(obj); phase == phaseSucceeded || phase == phaseFailed {
		return false, nil
	}

	patch, err := json.Marshal(map[string]any{`spec`: map[string]any{c.CancelField: c.CancelValue}})
	if err != nil {
		return false, err
	}

	if _, err := resource.Patch(ctx, name, types.MergePatchType, patch, metav1.PatchOptions{}); err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to cancel %s %s/%s: %w", c.CRDKind, namespace, name, err)
	}

	return true, nil

}

// checkResourceAccess lists the custom resource in the namespace (the health check's Kubernetes probe).
func (c *commandContext) checkResourceAccess(ctx context.Context, clients *kubeClients, namespace string) error {
	if _, err := clients.dynamic.Resource(c.gvr()).Namespace(namespace).List(ctx, metav1.ListOptions{Limit: 1}); err != nil {
		return fmt.Errorf("failed to list %s in namespace %q: %w", c.CRDPlural, namespace, err)
	}
	return nil
}

// reportWarningEvents writes recent Warning events for the run and its pods to stderr.
func (r *runner) reportWarningEvents(ctx context.Context) {

	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), cancelTimeout)
	defer cancel()

	names := map[string]bool{r.name: true}
	if r.cmd.PodSelectorLabel != `` {
		pods, err := r.clients.kube.CoreV1().Pods(r.namespace).List(ctx, metav1.ListOptions{LabelSelector: fmt.Sprintf("%s=%s", r.cmd.PodSelectorLabel, r.name)})
		if err == nil {
			for _, pod := range pods.Items {
				names[pod.Name] = true
			}
		}
	}

	var warnings []corev1.Event
	for name := range names {
		events, err := r.clients.kube.CoreV1().Events(r.namespace).List(ctx, metav1.ListOptions{
			FieldSelector: fmt.Sprintf("involvedObject.name=%s,type=%s", name, corev1.EventTypeWarning),
		})
		if err != nil {
			fmt.Fprintf(r.stderr, "eksjob: failed to list events for %s: %v\n", name, err)
			continue
		}
		for _, event := range events.Items {
			// field selectors are re-checked client-side
			if event.InvolvedObject.Name == name && event.Type == corev1.EventTypeWarning {
				warnings = append(warnings, event)
			}
		}
	}

	if len(warnings) == 0 {
		return
	}

	sort.Slice(warnings, func(i, j int) bool { return eventTime(warnings[i]).Before(eventTime(warnings[j])) })
	if len(warnings) > maxReportedEvents {
		warnings = warnings[len(warnings)-maxReportedEvents:]
	}

	fmt.Fprintf(r.stderr, "eksjob: warning events for %s:\n", r.ref())
	for _, event := range warnings {
		fmt.Fprintf(r.stderr, "  %s %s/%s %s: %s\n", eventTime(event).UTC().Format(time.RFC3339), event.InvolvedObject.Kind, event.InvolvedObject.Name, event.Reason, event.Message)
	}

}

func (r *runner) printTimes(obj *unstructured.Unstructured) {
	start, _, _ := unstructured.NestedString(obj.Object, statusStartTimePath...)
	end, _, _ := unstructured.NestedString(obj.Object, statusCompletionTimePath...)
	if start == `` && end == `` {
		return
	}
	line := fmt.Sprintf("eksjob: %s startTime=%s completionTime=%s", r.ref(), start, end)
	if s, err := time.Parse(time.RFC3339, start); err == nil {
		if e, err := time.Parse(time.RFC3339, end); err == nil {
			line += fmt.Sprintf(" duration=%s", e.Sub(s))
		}
	}
	fmt.Fprintln(r.stdout, line)
}

func phaseOf(obj *unstructured.Unstructured) string {
	if obj == nil {
		return ``
	}
	phase, _, _ := unstructured.NestedString(obj.Object, statusPhasePath...)
	return phase
}

func specKeys(obj *unstructured.Unstructured) []string {
	spec, _, _ := unstructured.NestedMap(obj.Object, `spec`)
	keys := make([]string, 0, len(spec))
	for k := range spec {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func eventTime(e corev1.Event) time.Time {
	switch {
	case !e.LastTimestamp.IsZero():
		return e.LastTimestamp.Time
	case !e.EventTime.IsZero():
		return e.EventTime.Time
	default:
		return e.FirstTimestamp.Time
	}
}
