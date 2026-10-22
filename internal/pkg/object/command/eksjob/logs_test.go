package eksjob

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	kubefake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

type buffer struct {
	mu sync.Mutex
	b  strings.Builder
}

func (b *buffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Write(p)
}

func (b *buffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.String()
}

const testSelector = testLabel + "=" + testRunName

func newTestStreamer(kube *kubefake.Clientset) (*logStreamer, *buffer, *buffer) {
	out, warn := &buffer{}, &buffer{}
	s := newLogStreamer(kube, testNamespace, testSelector, out, warn, time.Hour)
	return s, out, warn
}

// logRequests returns the container of every log request made (the fake does not record the
// pod name, so tests give containers distinct names).
func logRequests(kube *kubefake.Clientset) []string {
	var reqs []string
	for _, a := range kube.Actions() {
		if g, ok := a.(k8stesting.GenericAction); ok && a.GetSubresource() == "log" {
			opts := g.GetValue().(*corev1.PodLogOptions)
			if !opts.Follow || !opts.Timestamps {
				reqs = append(reqs, "unexpected-options")
			}
			reqs = append(reqs, opts.Container)
		}
	}
	return reqs
}

func TestLogStreamerSelectsRunPods(t *testing.T) {
	waiting := runPod("run-pending", testRunName, 0, false)
	waiting.Status.ContainerStatuses[0].State = corev1.ContainerState{Waiting: &corev1.ContainerStateWaiting{Reason: "ImagePullBackOff", Message: "image not found"}}
	creating := runPod("run-creating", testRunName, 0, false)
	creating.Status.ContainerStatuses[0].State = corev1.ContainerState{Waiting: &corev1.ContainerStateWaiting{Reason: "ContainerCreating"}}

	kube := kubefake.NewSimpleClientset(
		runPod("run-a", testRunName, 0, true),
		runPod("other-run", "heimdall-other", 0, true, "other"),
		waiting,
		creating,
	)
	s, out, warn := newTestStreamer(kube)
	s.start(context.Background())
	s.finish(time.Second)

	reqs := logRequests(kube)
	if len(reqs) != 1 || reqs[0] != "main" {
		t.Errorf("log requests = %v, want only run-a's main container", reqs)
	}
	// several run pods exist, so lines are prefixed with the pod name
	if got := out.String(); got != "[run-a] fake logs\n" {
		t.Errorf("stdout = %q", got)
	}
	if !strings.Contains(warn.String(), "pod run-pending container main is waiting: ImagePullBackOff: image not found") {
		t.Errorf("waiting reason not reported: %q", warn.String())
	}
	if strings.Contains(warn.String(), "ContainerCreating") || strings.Count(warn.String(), "ImagePullBackOff") != 1 {
		t.Errorf("unexpected warnings: %q", warn.String())
	}
}

func TestLogStreamerSinglePodIsNotPrefixed(t *testing.T) {
	kube := kubefake.NewSimpleClientset(runPod("run-a", testRunName, 0, true))
	s, out, _ := newTestStreamer(kube)
	s.start(context.Background())
	s.finish(time.Second)

	if got := out.String(); got != "fake logs\n" {
		t.Errorf("stdout = %q", got)
	}
}

func TestLogStreamerMultiContainerPrefix(t *testing.T) {
	kube := kubefake.NewSimpleClientset(runPod("run-a", testRunName, 0, true, "main", "sidecar"))
	s, out, _ := newTestStreamer(kube)
	s.start(context.Background())
	s.finish(time.Second)

	got := out.String()
	if !strings.Contains(got, "[run-a/main] fake logs\n") || !strings.Contains(got, "[run-a/sidecar] fake logs\n") {
		t.Errorf("stdout = %q", got)
	}
}

func TestLogStreamerStreamsEachContainerInstanceOnce(t *testing.T) {
	kube := kubefake.NewSimpleClientset(runPod("run-a", testRunName, 0, true))
	s, _, _ := newTestStreamer(kube)
	s.start(context.Background())
	s.stopDiscovery()
	s.wg.Wait()

	// a finished container that was fully streamed is not requested again
	s.sweep()
	s.wg.Wait()
	if reqs := logRequests(kube); len(reqs) != 1 {
		t.Fatalf("log requests = %v, want 1", reqs)
	}

	// a restart is a new container instance and is streamed again (pods appearing or
	// restarting late are picked up by later sweeps)
	restarted := runPod("run-a", testRunName, 1, false)
	if err := kube.Tracker().Update(corev1.SchemeGroupVersion.WithResource("pods"), restarted, testNamespace); err != nil {
		t.Fatal(err)
	}
	late := runPod("run-b", testRunName, 0, false, "late")
	if err := kube.Tracker().Add(late); err != nil {
		t.Fatal(err)
	}
	s.sweep()
	s.wg.Wait()
	s.cancel()

	reqs := logRequests(kube)
	if got := strings.Join(reqs, ","); len(reqs) != 3 || reqs[0] != "main" || (got != "main,main,late" && got != "main,late,main") {
		t.Errorf("log requests = %v, want run-a twice and run-b once", reqs)
	}
	if _, ok := s.streams["run-a/main/1"]; !ok {
		t.Errorf("restarted instance not tracked: %v", s.streams)
	}
}

func TestLogStreamerResumesFromLastTimestamp(t *testing.T) {
	s, out, _ := newTestStreamer(kubefake.NewSimpleClientset())
	state := &instanceState{}

	s.writeLine(state, "p", "c", false, "2026-01-01T00:00:01.000000001Z first\n")
	s.writeLine(state, "p", "c", false, "2026-01-01T00:00:02Z second\n")
	// a resumed stream (since=second granularity) replays lines already written
	s.writeLine(state, "p", "c", false, "2026-01-01T00:00:01.000000001Z first\n")
	s.writeLine(state, "p", "c", false, "2026-01-01T00:00:02Z second\n")
	s.writeLine(state, "p", "c", false, "2026-01-01T00:00:03Z third")
	s.writeLine(state, "p", "c", false, "no timestamp\n")

	if got, want := out.String(), "first\nsecond\nthird\nno timestamp\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
	if !state.last.Equal(time.Date(2026, 1, 1, 0, 0, 3, 0, time.UTC)) {
		t.Errorf("last = %v", state.last)
	}
}

func TestLogStreamerAbortStopsQuickly(t *testing.T) {
	kube := kubefake.NewSimpleClientset(runPod("run-a", testRunName, 0, false))
	s, _, _ := newTestStreamer(kube)
	s.start(context.Background())

	done := make(chan struct{})
	go func() {
		s.abort()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("abort did not return")
	}
}
