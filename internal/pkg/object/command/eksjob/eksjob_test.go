package eksjob

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	kubefake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"

	heimdallContext "github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/duration"
	"github.com/patterninc/heimdall/pkg/object"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	testNamespace = "runs"
	testJobID     = "2B1C0C5E-8F0A-4B5E-9D1E-0E0C9B9A1F00"
	testRunName   = "heimdall-2b1c0c5e-8f0a-4b5e-9d1e-0e0c9b9a1f00"
	testLabel     = "example.com/run"
)

var testGVR = schema.GroupVersionResource{Group: "example.com", Version: "v1alpha1", Resource: "exampleruns"}

// contractCommandContext mirrors the command context in the shared contract, with example
// CRD names and the placeholder image the first config ships with.
func contractCommandContext() map[string]any {
	return map[string]any{
		"crd_group":     "example.com",
		"crd_version":   "v1alpha1",
		"crd_kind":      "ExampleRun",
		"crd_plural":    "exampleruns",
		"poll_interval": "15s",
		"spec_defaults": map[string]any{
			"image":              "PLACEHOLDER-caterpillar-eks-image-tbd",
			"serviceAccountName": "runner",
			"ttl":                "1h",
			"nodeSelector":       map[string]any{"workload_type": "runs"},
			"tolerations": []any{
				map[string]any{"key": "dedicated", "operator": "Equal", "value": "runs", "effect": "NoSchedule"},
			},
		},
		"size_field": "resources",
		"size_profiles": map[string]any{
			"xs":      map[string]any{"requests": map[string]any{"cpu": "500m", "memory": "1Gi"}, "limits": map[string]any{"memory": "1Gi"}},
			"s":       map[string]any{"requests": map[string]any{"cpu": "1", "memory": "2Gi"}, "limits": map[string]any{"memory": "2Gi"}},
			"m":       map[string]any{"requests": map[string]any{"cpu": "2", "memory": "4Gi"}, "limits": map[string]any{"memory": "4Gi"}},
			"l":       map[string]any{"requests": map[string]any{"cpu": "4", "memory": "8Gi"}, "limits": map[string]any{"memory": "8Gi"}},
			"l_mem_m": map[string]any{"requests": map[string]any{"cpu": "4", "memory": "16Gi"}, "limits": map[string]any{"memory": "16Gi"}},
			"l_mem_l": map[string]any{"requests": map[string]any{"cpu": "4", "memory": "20Gi"}, "limits": map[string]any{"memory": "20Gi"}},
		},
		"cancel_field":       "setState",
		"cancel_value":       "cancel",
		"pod_selector_label": testLabel,
	}
}

func newCommand(t *testing.T, ctx map[string]any) *commandContext {
	t.Helper()
	h, err := New(heimdallContext.New(ctx))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return h.(*commandContext)
}

func testCluster() *cluster.Cluster {
	return &cluster.Cluster{Object: object.Object{Name: "eks", Context: heimdallContext.New(map[string]any{
		"role_arn":     "arn:aws:iam::123456789012:role/example",
		"region":       "us-west-2",
		"cluster_name": "example",
		"namespace":    testNamespace,
	})}}
}

func testJob(ctx map[string]any) *job.Job {
	j := &job.Job{}
	j.ID = testJobID
	j.Context = heimdallContext.New(ctx)
	return j
}

func defaultJobContext() map[string]any {
	return map[string]any{
		"size":   "s",
		"labels": map[string]any{"owner": "jane.doe@example.com", "dag_id": "ingest_orders", "task_id": "load"},
		"spec":   map[string]any{"pipeline": "tasks: []\n", "taskCount": 1, "timeout": "7200s", "retries": 0},
	}
}

type fakes struct {
	dynamic *dynamicfake.FakeDynamicClient
	kube    *kubefake.Clientset
}

func newFakes(objects ...runtime.Object) *fakes {
	var dyn []runtime.Object
	var kube []runtime.Object
	for _, o := range objects {
		if _, ok := o.(*unstructured.Unstructured); ok {
			dyn = append(dyn, o)
		} else {
			kube = append(kube, o)
		}
	}
	return &fakes{
		dynamic: dynamicfake.NewSimpleDynamicClientWithCustomListKinds(runtime.NewScheme(), map[schema.GroupVersionResource]string{testGVR: "ExampleRunList"}, dyn...),
		kube:    kubefake.NewSimpleClientset(kube...),
	}
}

// install points the command at the fakes and speeds up polling.
func (f *fakes) install(c *commandContext) *int {
	calls := 0
	c.newClients = func(context.Context, *clusterContext, string, bool) (*kubeClients, error) {
		calls++
		return &kubeClients{dynamic: f.dynamic, kube: f.kube}, nil
	}
	c.PollInterval = duration.Duration(time.Millisecond)
	return &calls
}

// phases makes successive GETs of the run report the given status phases (the last one repeats),
// standing in for the operator.
func (f *fakes) phases(failureReason string, phases ...string) {
	var mu sync.Mutex
	n := 0
	f.dynamic.PrependReactor("get", testGVR.Resource, func(action k8stesting.Action) (bool, runtime.Object, error) {
		mu.Lock()
		phase := phases[min(n, len(phases)-1)]
		n++
		mu.Unlock()
		get := action.(k8stesting.GetAction)
		obj, err := f.dynamic.Tracker().Get(testGVR, get.GetNamespace(), get.GetName())
		if err != nil {
			return true, nil, err
		}
		u := obj.(*unstructured.Unstructured).DeepCopy()
		status := map[string]any{"phase": phase, "startTime": "2026-01-01T00:00:00Z"}
		if phase == phaseSucceeded || phase == phaseFailed {
			status["completionTime"] = "2026-01-01T00:05:00Z"
		}
		if phase == phaseFailed {
			status["failureReason"] = failureReason
		}
		u.Object["status"] = status
		return true, u, nil
	})
}

func (f *fakes) run(t *testing.T) *unstructured.Unstructured {
	t.Helper()
	obj, err := f.dynamic.Tracker().Get(testGVR, testNamespace, testRunName)
	if err != nil {
		t.Fatalf("run not found: %v", err)
	}
	return obj.(*unstructured.Unstructured)
}

type testRuntime struct {
	*plugin.Runtime
}

func newRuntime(t *testing.T) *testRuntime {
	t.Helper()
	dir := t.TempDir()
	stdout, err := os.Create(filepath.Join(dir, "stdout"))
	if err != nil {
		t.Fatal(err)
	}
	stderr, err := os.Create(filepath.Join(dir, "stderr"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { stdout.Close(); stderr.Close() })
	return &testRuntime{&plugin.Runtime{WorkingDirectory: dir, Stdout: stdout, Stderr: stderr}}
}

func (r *testRuntime) stdout(t *testing.T) string {
	b, err := os.ReadFile(r.Stdout.Name())
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func (r *testRuntime) stderr(t *testing.T) string {
	b, err := os.ReadFile(r.Stderr.Name())
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func runPod(name, runName string, restarts int32, terminated bool, containers ...string) *corev1.Pod {
	if len(containers) == 0 {
		containers = []string{"main"}
	}
	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: testNamespace, Labels: map[string]string{testLabel: runName}}}
	for _, c := range containers {
		pod.Spec.Containers = append(pod.Spec.Containers, corev1.Container{Name: c})
		st := corev1.ContainerStatus{Name: c, RestartCount: restarts}
		if terminated {
			st.State.Terminated = &corev1.ContainerStateTerminated{ExitCode: 0}
		} else {
			st.State.Running = &corev1.ContainerStateRunning{}
		}
		pod.Status.ContainerStatuses = append(pod.Status.ContainerStatuses, st)
	}
	return pod
}

// --- New / context parsing ---

func TestNewAcceptsContractCommandContextWithPlaceholderImage(t *testing.T) {
	c := newCommand(t, contractCommandContext())

	if time.Duration(c.PollInterval) != 15*time.Second {
		t.Errorf("poll_interval = %v", time.Duration(c.PollInterval))
	}
	if c.SpecDefaults["image"] != "PLACEHOLDER-caterpillar-eks-image-tbd" {
		t.Errorf("spec_defaults.image = %v", c.SpecDefaults["image"])
	}
	if len(c.SizeProfiles) != 6 || c.CancelField != "setState" || c.CancelValue != "cancel" || c.PodSelectorLabel != testLabel {
		t.Errorf("unexpected command context: %+v", c)
	}
	if c.newClients == nil {
		t.Error("default client factory not set")
	}
}

func TestNewDefaults(t *testing.T) {
	ctx := map[string]any{"crd_group": "example.com", "crd_version": "v1", "crd_kind": "ExampleRun", "crd_plural": "exampleruns", "cancel_field": "stop"}
	c := newCommand(t, ctx)
	if time.Duration(c.PollInterval) != defaultPollInterval || c.SizeField != defaultSizeField || c.CancelValue != defaultCancelValue {
		t.Errorf("defaults not applied: %+v", c)
	}
}

func TestNewDefersValueChecksToExecute(t *testing.T) {
	ctx := contractCommandContext()
	ctx["size_profiles"].(map[string]any)["bad"] = map[string]any{"requests": map[string]any{"cpu": "lots"}}
	ctx["unknown_field"] = "ignored"
	c := newCommand(t, ctx)

	_, err := c.buildRun(testRunName, testNamespace, testJobID, &jobContext{Size: "bad"}, &strings.Builder{})
	if err == nil || !strings.Contains(err.Error(), `invalid size profile "bad"`) {
		t.Fatalf("expected invalid size profile error, got %v", err)
	}
}

func TestNewRejectsMalformedCommandContext(t *testing.T) {
	cases := map[string]func(map[string]any){
		"missing crd_group":     func(m map[string]any) { delete(m, "crd_group") },
		"missing crd_version":   func(m map[string]any) { delete(m, "crd_version") },
		"missing crd_kind":      func(m map[string]any) { delete(m, "crd_kind") },
		"missing crd_plural":    func(m map[string]any) { delete(m, "crd_plural") },
		"profile not an object": func(m map[string]any) { m["size_profiles"] = map[string]any{"s": "big"} },
		"spec_defaults a list":  func(m map[string]any) { m["spec_defaults"] = []any{"x"} },
		"bad poll_interval":     func(m map[string]any) { m["poll_interval"] = "soon" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := contractCommandContext()
			mutate(ctx)
			if _, err := New(heimdallContext.New(ctx)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
	if _, err := New(nil); !errors.Is(err, ErrMissingCRDField) {
		t.Fatalf("New(nil) = %v", err)
	}
}

func TestParseClusterContext(t *testing.T) {
	cc, err := parseClusterContext(testCluster())
	if err != nil {
		t.Fatal(err)
	}
	if cc.RoleARN != "arn:aws:iam::123456789012:role/example" || cc.Region != "us-west-2" || cc.ClusterName != "example" || cc.Namespace != testNamespace {
		t.Errorf("unexpected cluster context: %+v", cc)
	}

	cc, err = parseClusterContext(&cluster.Cluster{Object: object.Object{Context: heimdallContext.New(map[string]any{"cluster_name": "c", "extra": 1})}})
	if err != nil || cc.Namespace != defaultNamespace {
		t.Errorf("expected default namespace, got %+v, %v", cc, err)
	}

	if _, err := parseClusterContext(&cluster.Cluster{}); !errors.Is(err, ErrMissingClusterName) {
		t.Errorf("expected ErrMissingClusterName, got %v", err)
	}
	if _, err := parseClusterContext(nil); !errors.Is(err, ErrMissingCluster) {
		t.Errorf("expected ErrMissingCluster, got %v", err)
	}
}

func TestParseJobContext(t *testing.T) {
	ctx := defaultJobContext()
	ctx["ignored"] = true
	jc, err := parseJobContext(testJob(ctx))
	if err != nil {
		t.Fatal(err)
	}
	if jc.Size != "s" || jc.Labels["dag_id"] != "ingest_orders" || jc.Spec["pipeline"] != "tasks: []\n" {
		t.Errorf("unexpected job context: %+v", jc)
	}

	if _, err := parseJobContext(testJob(map[string]any{"spec": "not-a-map"})); err == nil {
		t.Error("expected error for a non-object spec")
	}

	jc, err = parseJobContext(&job.Job{})
	if err != nil || jc.Size != "" {
		t.Errorf("empty job context: %+v, %v", jc, err)
	}
}

// --- naming and spec building ---

func TestRunName(t *testing.T) {
	name, err := runName(testJobID)
	if err != nil || name != testRunName {
		t.Fatalf("runName = %q, %v", name, err)
	}

	name, err = runName("My_Job.ID")
	if err != nil || name != "heimdall-my-job-id" {
		t.Errorf("runName = %q, %v", name, err)
	}

	name, err = runName(strings.Repeat("a", 30) + "-" + strings.Repeat("b", 40))
	if err != nil || len(name) > 63 || strings.HasSuffix(name, "-") {
		t.Errorf("long runName = %q (%d), %v", name, len(name), err)
	}

	if _, err := runName("__"); err == nil {
		t.Error("expected error for an unusable job ID")
	}
}

func TestBuildRunMergesDefaultsJobSpecAndSize(t *testing.T) {
	c := newCommand(t, contractCommandContext())
	jc := &jobContext{
		Size:   "m",
		Labels: map[string]any{"owner": "jane.doe@example.com", "dag_id": "ingest_orders"},
		Spec: map[string]any{
			"pipeline":     "tasks: []\n",
			"image":        "registry.example.com/runner:2",
			"nodeSelector": map[string]any{"pool": "other"},
			"resources":    map[string]any{"requests": map[string]any{"cpu": "64"}},
		},
	}

	obj, err := c.buildRun(testRunName, testNamespace, testJobID, jc, &strings.Builder{})
	if err != nil {
		t.Fatal(err)
	}

	if obj.GetAPIVersion() != "example.com/v1alpha1" || obj.GetKind() != "ExampleRun" || obj.GetName() != testRunName || obj.GetNamespace() != testNamespace {
		t.Errorf("unexpected object meta: %s %s %s/%s", obj.GetAPIVersion(), obj.GetKind(), obj.GetNamespace(), obj.GetName())
	}

	spec := obj.Object["spec"].(map[string]any)
	if spec["image"] != "registry.example.com/runner:2" {
		t.Errorf("job spec image should override default, got %v", spec["image"])
	}
	if spec["serviceAccountName"] != "runner" || spec["ttl"] != "1h" {
		t.Errorf("defaults not merged: %v", spec)
	}
	if ns := spec["nodeSelector"].(map[string]any); ns["pool"] != "other" || ns["workload_type"] != nil {
		t.Errorf("top-level job spec keys should replace defaults, got %v", ns)
	}
	if _, ok := spec["tolerations"].([]any); !ok {
		t.Errorf("tolerations default missing: %v", spec["tolerations"])
	}

	// the size profile is authoritative for spec.<size_field>
	cpu, _, _ := unstructured.NestedString(obj.Object, "spec", "resources", "requests", "cpu")
	memory, _, _ := unstructured.NestedString(obj.Object, "spec", "resources", "requests", "memory")
	limit, _, _ := unstructured.NestedString(obj.Object, "spec", "resources", "limits", "memory")
	_, hasCPULimit, _ := unstructured.NestedFieldNoCopy(obj.Object, "spec", "resources", "limits", "cpu")
	if cpu != "2" || memory != "4Gi" || limit != "4Gi" || hasCPULimit {
		t.Errorf("unexpected resources: cpu=%q memory=%q limit=%q cpuLimit=%v", cpu, memory, limit, hasCPULimit)
	}

	labels := obj.GetLabels()
	want := map[string]string{
		"owner":        "jane.doe_example.com",
		"dag_id":       "ingest_orders",
		labelManagedBy: labelManagedByValue,
		labelJobID:     testJobID,
	}
	for k, v := range want {
		if labels[k] != v {
			t.Errorf("label %s = %q, want %q", k, labels[k], v)
		}
	}

	// the shared command defaults must not be mutated by a job
	spec["serviceAccountName"] = "changed"
	spec["tolerations"].([]any)[0].(map[string]any)["value"] = "changed"
	if c.SpecDefaults["serviceAccountName"] != "runner" || c.SpecDefaults["tolerations"].([]any)[0].(map[string]any)["value"] != "runs" {
		t.Error("buildRun mutated spec_defaults")
	}
}

func TestBuildRunSizeErrors(t *testing.T) {
	c := newCommand(t, contractCommandContext())

	_, err := c.buildRun(testRunName, testNamespace, testJobID, &jobContext{Size: "xxl"}, &strings.Builder{})
	if err == nil || !strings.Contains(err.Error(), `unknown job size "xxl"`) || !strings.Contains(err.Error(), "l, l_mem_l, l_mem_m, m, s, xs") {
		t.Errorf("unknown size error = %v", err)
	}

	_, err = c.buildRun(testRunName, testNamespace, testJobID, &jobContext{}, &strings.Builder{})
	if err == nil || !strings.Contains(err.Error(), "size is required") {
		t.Errorf("missing size error = %v", err)
	}

	// without size profiles, size is optional and no resources are written
	ctx := contractCommandContext()
	delete(ctx, "size_profiles")
	obj, err := newCommand(t, ctx).buildRun(testRunName, testNamespace, testJobID, &jobContext{}, &strings.Builder{})
	if err != nil {
		t.Fatal(err)
	}
	if _, found, _ := unstructured.NestedFieldNoCopy(obj.Object, "spec", "resources"); found {
		t.Error("resources set without a size")
	}
}

func TestBuildRunRequiredSpecFields(t *testing.T) {
	ctx := contractCommandContext()
	ctx["required_spec_fields"] = []any{"pipeline", "image"}
	c := newCommand(t, ctx)

	_, err := c.buildRun(testRunName, testNamespace, testJobID, &jobContext{Size: "s"}, &strings.Builder{})
	if err == nil || !strings.Contains(err.Error(), "pipeline") || strings.Contains(err.Error(), "image") {
		t.Errorf("required field error = %v", err)
	}

	if _, err := c.buildRun(testRunName, testNamespace, testJobID, &jobContext{Size: "s", Spec: map[string]any{"pipeline": "x"}}, &strings.Builder{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSanitizeLabels(t *testing.T) {
	in := map[string]any{
		"owner":                 "jane.doe@example.com",
		"dag_id":                "ingest_orders",
		"Task ID":               "load data",
		"example.COM/team":      "data",
		"count":                 3,
		"enabled":               true,
		"empty":                 nil,
		"long":                  strings.Repeat("v", 70),
		"edge":                  "-_value_-",
		"!!!":                   "dropped",
		"bad prefix!/name":      "x",
		"-/name":                "dropped",
		strings.Repeat("k", 70): "truncated",
	}
	labels, dropped := sanitizeLabels(in)

	want := map[string]string{
		"owner":                 "jane.doe_example.com",
		"dag_id":                "ingest_orders",
		"Task_ID":               "load_data",
		"example.com/team":      "data",
		"count":                 "3",
		"enabled":               "true",
		"empty":                 "",
		"long":                  strings.Repeat("v", 63),
		"edge":                  "value",
		"bad-prefix/name":       "x",
		strings.Repeat("k", 63): "truncated",
	}
	for k, v := range want {
		if got, ok := labels[k]; !ok || got != v {
			t.Errorf("label %q = %q (present %v), want %q", k, got, ok, v)
		}
	}
	if len(labels) != len(want) {
		t.Errorf("unexpected labels: %v", labels)
	}
	if strings.Join(dropped, ",") != "!!!,-/name" {
		t.Errorf("dropped = %v", dropped)
	}
}

func TestBuildRunWarnsAboutDroppedLabels(t *testing.T) {
	c := newCommand(t, contractCommandContext())
	var warnings strings.Builder
	obj, err := c.buildRun(testRunName, testNamespace, testJobID, &jobContext{Size: "s", Labels: map[string]any{"???": "x", labelManagedBy: "someone-else"}}, &warnings)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(warnings.String(), `dropping label "???"`) {
		t.Errorf("warnings = %q", warnings.String())
	}
	if obj.GetLabels()[labelManagedBy] != labelManagedByValue {
		t.Error("job labels must not override managed labels")
	}
}

// --- Execute ---

func TestExecuteSucceedsAndStreamsLogs(t *testing.T) {
	c := newCommand(t, contractCommandContext())
	f := newFakes(
		runPod("run-a", testRunName, 0, true),
		runPod("run-b", testRunName, 0, true),
		runPod("other", "heimdall-someone-else", 0, true),
	)
	f.install(c)
	f.phases("", "", phasePending, phaseRunning, phaseSucceeded)

	r := newRuntime(t)
	j := testJob(defaultJobContext())
	if err := c.Execute(context.Background(), r.Runtime, j, testCluster()); err != nil {
		t.Fatalf("Execute: %v\nstderr: %s", err, r.stderr(t))
	}

	out := r.stdout(t)
	for _, want := range []string{"created ExampleRun runs/" + testRunName, "phase: Running", "succeeded", "[run-a] fake logs", "[run-b] fake logs", "duration=5m0s"} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "[other]") {
		t.Errorf("logs of an unrelated pod were streamed:\n%s", out)
	}

	run := f.run(t)
	if pipeline, _, _ := unstructured.NestedString(run.Object, "spec", "pipeline"); pipeline != "tasks: []\n" {
		t.Errorf("spec.pipeline = %q", pipeline)
	}
	if run.GetLabels()["dag_id"] != "ingest_orders" {
		t.Errorf("labels = %v", run.GetLabels())
	}
	if j.Outputs()[outputRunName] != testRunName || j.Outputs()[outputRunNamespace] != testNamespace {
		t.Errorf("outputs = %v", j.Outputs())
	}
}

func TestExecuteFailedSurfacesFailureReason(t *testing.T) {
	c := newCommand(t, contractCommandContext())
	event := &corev1.Event{
		ObjectMeta:     metav1.ObjectMeta{Name: "ev1", Namespace: testNamespace},
		InvolvedObject: corev1.ObjectReference{Kind: "Pod", Name: "run-a"},
		Type:           corev1.EventTypeWarning,
		Reason:         "BackOff",
		Message:        "Back-off restarting failed container",
		LastTimestamp:  metav1.NewTime(time.Unix(1700000000, 0)),
	}
	normal := event.DeepCopy()
	normal.Name, normal.Type, normal.Message = "ev2", corev1.EventTypeNormal, "Pulled image"
	f := newFakes(runPod("run-a", testRunName, 0, true), event, normal)
	f.install(c)
	f.phases("task 0 exited with code 2", phaseRunning, phaseFailed)

	r := newRuntime(t)
	err := c.Execute(context.Background(), r.Runtime, testJob(defaultJobContext()), testCluster())
	if err == nil || !strings.Contains(err.Error(), "task 0 exited with code 2") || !strings.Contains(err.Error(), testRunName) {
		t.Fatalf("expected failureReason in error, got %v", err)
	}

	stderr := r.stderr(t)
	if !strings.Contains(stderr, "Back-off restarting failed container") || strings.Contains(stderr, "Pulled image") {
		t.Errorf("stderr should list warning events only:\n%s", stderr)
	}
}

func TestExecuteFailedWithoutReason(t *testing.T) {
	c := newCommand(t, contractCommandContext())
	f := newFakes()
	f.install(c)
	f.phases("", phaseFailed)

	err := c.Execute(context.Background(), newRuntime(t).Runtime, testJob(defaultJobContext()), testCluster())
	if err == nil || !strings.Contains(err.Error(), "no failureReason reported") {
		t.Fatalf("got %v", err)
	}
}

func TestExecuteAttachesToExistingRun(t *testing.T) {
	c := newCommand(t, contractCommandContext())
	existing := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "example.com/v1alpha1",
		"kind":       "ExampleRun",
		"metadata":   map[string]any{"name": testRunName, "namespace": testNamespace},
		"spec":       map[string]any{"pipeline": "original"},
		"status":     map[string]any{"phase": phaseRunning},
	}}
	f := newFakes(existing)
	f.install(c)
	f.phases("", phaseRunning, phaseSucceeded)

	r := newRuntime(t)
	if err := c.Execute(context.Background(), r.Runtime, testJob(defaultJobContext()), testCluster()); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(r.stdout(t), "already exists (phase Running); attaching to it") {
		t.Errorf("stdout:\n%s", r.stdout(t))
	}
	if pipeline, _, _ := unstructured.NestedString(f.run(t).Object, "spec", "pipeline"); pipeline != "original" {
		t.Errorf("existing run was modified: spec.pipeline = %q", pipeline)
	}
	creates := 0
	for _, a := range f.dynamic.Actions() {
		if a.GetVerb() == "create" {
			creates++
		}
	}
	if creates != 1 {
		t.Errorf("expected exactly one create attempt, got %d", creates)
	}
}

func TestExecuteValidatesBeforeContactingCluster(t *testing.T) {
	c := newCommand(t, contractCommandContext())
	calls := newFakes().install(c)

	err := c.Execute(context.Background(), newRuntime(t).Runtime, testJob(map[string]any{"size": "huge"}), testCluster())
	if err == nil || !strings.Contains(err.Error(), "unknown job size") {
		t.Fatalf("got %v", err)
	}
	if *calls != 0 {
		t.Error("clients were created for an invalid job")
	}
}

func TestExecuteCanceledPatchesCancelField(t *testing.T) {
	c := newCommand(t, contractCommandContext())
	f := newFakes()
	f.install(c)
	f.phases("", phaseRunning)

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(100*time.Millisecond, cancel)

	r := newRuntime(t)
	err := c.Execute(ctx, r.Runtime, testJob(defaultJobContext()), testCluster())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}

	state, _, _ := unstructured.NestedString(f.run(t).Object, "spec", "setState")
	if state != "cancel" {
		t.Errorf("spec.setState = %q, want cancel", state)
	}
	if !strings.Contains(r.stdout(t), "requested cancellation") {
		t.Errorf("stdout:\n%s", r.stdout(t))
	}
}

func TestExecutePollFailureCancelsRun(t *testing.T) {
	c := newCommand(t, contractCommandContext())
	f := newFakes()
	f.install(c)

	// the API fails until the poll gives up, then recovers for the cancel
	var mu sync.Mutex
	gets := 0
	f.dynamic.PrependReactor("get", testGVR.Resource, func(k8stesting.Action) (bool, runtime.Object, error) {
		mu.Lock()
		defer mu.Unlock()
		gets++
		if gets <= maxConsecutivePollErrors {
			return true, nil, errors.New("connection refused")
		}
		return false, nil, nil
	})

	r := newRuntime(t)
	err := c.Execute(context.Background(), r.Runtime, testJob(defaultJobContext()), testCluster())
	if err == nil || !strings.Contains(err.Error(), "connection refused") {
		t.Fatalf("expected the poll error, got %v", err)
	}

	if state, _, _ := unstructured.NestedString(f.run(t).Object, "spec", "setState"); state != "cancel" {
		t.Errorf("spec.setState = %q, want cancel", state)
	}
	if !strings.Contains(r.stdout(t), "requested cancellation") {
		t.Errorf("stdout:\n%s", r.stdout(t))
	}
}

func TestExecuteDeletedRunIsNotPatched(t *testing.T) {
	c := newCommand(t, contractCommandContext())
	f := newFakes()
	f.install(c)
	f.dynamic.PrependReactor("get", testGVR.Resource, func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, apierrors.NewNotFound(testGVR.GroupResource(), testRunName)
	})

	err := c.Execute(context.Background(), newRuntime(t).Runtime, testJob(defaultJobContext()), testCluster())
	if err == nil || !strings.Contains(err.Error(), "deleted before it finished") {
		t.Fatalf("got %v", err)
	}
	for _, a := range f.dynamic.Actions() {
		if a.GetVerb() == "patch" {
			t.Error("a deleted run was patched")
		}
	}
}

// --- Cleanup / HealthCheck ---

func existingRun(phase string) *unstructured.Unstructured {
	u := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "example.com/v1alpha1",
		"kind":       "ExampleRun",
		"metadata":   map[string]any{"name": testRunName, "namespace": testNamespace},
		"spec":       map[string]any{"pipeline": "p"},
	}}
	if phase != "" {
		u.Object["status"] = map[string]any{"phase": phase}
	}
	return u
}

func TestCleanupPatchesCancelField(t *testing.T) {
	c := newCommand(t, contractCommandContext())
	f := newFakes(existingRun(phaseRunning))
	f.install(c)

	if err := c.Cleanup(context.Background(), testJobID, testCluster()); err != nil {
		t.Fatal(err)
	}

	run := f.run(t)
	if state, _, _ := unstructured.NestedString(run.Object, "spec", "setState"); state != "cancel" {
		t.Errorf("spec.setState = %q", state)
	}
	if pipeline, _, _ := unstructured.NestedString(run.Object, "spec", "pipeline"); pipeline != "p" {
		t.Error("merge patch must keep the rest of the spec")
	}
	for _, a := range f.dynamic.Actions() {
		if a.GetVerb() == "delete" {
			t.Error("Cleanup must not delete the run")
		}
		if p, ok := a.(k8stesting.PatchAction); ok && p.GetPatchType() != "application/merge-patch+json" {
			t.Errorf("patch type = %s", p.GetPatchType())
		}
	}
}

func TestCleanupSkipsFinishedOrMissingRuns(t *testing.T) {
	for _, tc := range []struct {
		name    string
		objects []runtime.Object
	}{
		{"succeeded", []runtime.Object{existingRun(phaseSucceeded)}},
		{"failed", []runtime.Object{existingRun(phaseFailed)}},
		{"missing", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := newCommand(t, contractCommandContext())
			f := newFakes(tc.objects...)
			f.install(c)
			if err := c.Cleanup(context.Background(), testJobID, testCluster()); err != nil {
				t.Fatal(err)
			}
			for _, a := range f.dynamic.Actions() {
				if a.GetVerb() == "patch" {
					t.Error("finished or missing run was patched")
				}
			}
		})
	}
}

func TestCleanupWithoutCancelFieldIsNoop(t *testing.T) {
	ctx := contractCommandContext()
	delete(ctx, "cancel_field")
	delete(ctx, "cancel_value")
	c := newCommand(t, ctx)
	calls := newFakes(existingRun(phaseRunning)).install(c)

	if err := c.Cleanup(context.Background(), testJobID, testCluster()); err != nil {
		t.Fatal(err)
	}
	if *calls != 0 {
		t.Error("clients created although cancellation is disabled")
	}
}

func TestHealthCheck(t *testing.T) {
	c := newCommand(t, contractCommandContext())
	f := newFakes(existingRun(phaseRunning))
	var session string
	var active bool
	c.newClients = func(_ context.Context, cc *clusterContext, s string, requireActive bool) (*kubeClients, error) {
		session, active = s, requireActive
		return &kubeClients{dynamic: f.dynamic, kube: f.kube}, nil
	}

	if err := c.HealthCheck(context.Background(), testCluster()); err != nil {
		t.Fatal(err)
	}
	if session != healthSessionName || !active {
		t.Errorf("session name = %q, requireActive = %v", session, active)
	}
	listed := false
	for _, a := range f.dynamic.Actions() {
		if a.GetVerb() == "list" && a.GetNamespace() == testNamespace && a.GetResource() == testGVR {
			listed = true
		}
	}
	if !listed {
		t.Error("health check did not list the custom resource")
	}

	f.dynamic.PrependReactor("list", testGVR.Resource, func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("forbidden")
	})
	if err := c.HealthCheck(context.Background(), testCluster()); err == nil || !strings.Contains(err.Error(), "forbidden") {
		t.Errorf("expected list error, got %v", err)
	}

	c.newClients = func(context.Context, *clusterContext, string, bool) (*kubeClients, error) {
		return nil, errors.New(`EKS cluster "example" is not active: CREATING`)
	}
	if err := c.HealthCheck(context.Background(), testCluster()); err == nil || !strings.Contains(err.Error(), "not active") {
		t.Errorf("expected cluster error, got %v", err)
	}
}
