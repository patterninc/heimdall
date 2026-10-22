package eksjob

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// waiting reasons that are part of a normal start and not worth reporting
var benignWaitingReasons = map[string]bool{
	``:                  true,
	`ContainerCreating`: true,
	`PodInitializing`:   true,
}

// logStreamer follows the logs of every container of every pod matching a label selector.
//
// Pods are discovered by listing on an interval, so pods that appear late are picked up. Each
// container instance (pod, container, restart count) is streamed once with follow=true; a
// restarted container is a new instance and is streamed again. A stream that drops while its
// container is still running is resumed from the last timestamp seen, skipping lines already
// written.
type logStreamer struct {
	kube      kubernetes.Interface
	namespace string
	selector  string
	out       io.Writer
	warn      io.Writer
	interval  time.Duration

	ctx      context.Context
	cancel   context.CancelFunc
	stopLoop chan struct{}
	stopOnce sync.Once
	loopDone chan struct{}
	wg       sync.WaitGroup

	mu       sync.Mutex
	pods     map[string]bool
	streams  map[string]*instanceState
	reported map[string]bool
}

type instanceState struct {
	active bool      // a stream is open
	ended  bool      // a stream reached EOF
	last   time.Time // timestamp of the last line written
	warned bool      // a stream error was reported
}

func newLogStreamer(kube kubernetes.Interface, namespace, selector string, out, warn io.Writer, interval time.Duration) *logStreamer {
	return &logStreamer{
		kube:      kube,
		namespace: namespace,
		selector:  selector,
		out:       out,
		warn:      warn,
		interval:  interval,
		stopLoop:  make(chan struct{}),
		loopDone:  make(chan struct{}),
		pods:      map[string]bool{},
		streams:   map[string]*instanceState{},
		reported:  map[string]bool{},
	}
}

// start begins pod discovery in the background.
func (s *logStreamer) start(ctx context.Context) {
	s.ctx, s.cancel = context.WithCancel(ctx)
	go s.loop()
}

func (s *logStreamer) loop() {
	defer close(s.loopDone)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		s.sweep()
		select {
		case <-s.ctx.Done():
			return
		case <-s.stopLoop:
			return
		case <-ticker.C:
		}
	}
}

// finish stops discovery, picks up containers that finished since the last sweep and waits up
// to grace for open streams to reach the end of their logs.
func (s *logStreamer) finish(grace time.Duration) {

	s.stopDiscovery()
	s.sweep()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(grace):
		fmt.Fprintf(s.warn, "eksjob: log streams still open after %s; closing them\n", grace)
	}

	s.cancel()
	s.wg.Wait()

}

// abort stops discovery and closes all streams immediately.
func (s *logStreamer) abort() {
	s.cancel()
	s.stopDiscovery()
	s.wg.Wait()
}

func (s *logStreamer) stopDiscovery() {
	s.stopOnce.Do(func() { close(s.stopLoop) })
	<-s.loopDone
}

// sweep lists the selected pods and starts a stream for each container instance that has
// started and is not already streaming or fully streamed.
func (s *logStreamer) sweep() {

	if s.ctx.Err() != nil {
		return
	}

	pods, err := s.kube.CoreV1().Pods(s.namespace).List(s.ctx, metav1.ListOptions{LabelSelector: s.selector})
	if err != nil {
		if s.ctx.Err() == nil {
			fmt.Fprintf(s.warn, "eksjob: failed to list pods %s: %v\n", s.selector, err)
		}
		return
	}

	// record every pod before starting streams, so prefixing does not depend on stream timing
	s.mu.Lock()
	for _, pod := range pods.Items {
		s.pods[pod.Name] = true
	}
	s.mu.Unlock()

	for i := range pods.Items {
		pod := &pods.Items[i]

		statuses := make(map[string]corev1.ContainerStatus, len(pod.Status.ContainerStatuses))
		for _, st := range pod.Status.ContainerStatuses {
			statuses[st.Name] = st
		}

		for _, container := range pod.Spec.Containers {
			st, found := statuses[container.Name]
			if !found {
				continue
			}
			if waiting := st.State.Waiting; waiting != nil {
				s.reportWaiting(pod.Name, container.Name, waiting)
				continue
			}
			if st.State.Running == nil && st.State.Terminated == nil {
				continue
			}

			key := fmt.Sprintf("%s/%s/%d", pod.Name, container.Name, st.RestartCount)

			s.mu.Lock()
			state, found := s.streams[key]
			if !found {
				state = &instanceState{}
				s.streams[key] = state
			}
			skip := state.active || (state.ended && st.State.Terminated != nil)
			if !skip {
				state.active = true
			}
			s.mu.Unlock()

			if skip {
				continue
			}

			s.wg.Add(1)
			go s.stream(pod.Name, container.Name, len(pod.Spec.Containers) > 1, key)
		}
	}

}

// stream follows one container instance's log until EOF or cancellation.
func (s *logStreamer) stream(podName, containerName string, multiContainer bool, key string) {

	defer s.wg.Done()

	s.mu.Lock()
	state := s.streams[key]
	since := state.last
	s.mu.Unlock()

	opts := &corev1.PodLogOptions{Container: containerName, Follow: true, Timestamps: true}
	if !since.IsZero() {
		t := metav1.NewTime(since)
		opts.SinceTime = &t
	}

	rc, err := s.kube.CoreV1().Pods(s.namespace).GetLogs(podName, opts).Stream(s.ctx)
	if err != nil {
		s.mu.Lock()
		state.active = false
		warn := !state.warned && s.ctx.Err() == nil
		state.warned = true
		s.mu.Unlock()
		if warn {
			fmt.Fprintf(s.warn, "eksjob: failed to stream logs of pod %s container %s (will retry): %v\n", podName, containerName, err)
		}
		return
	}
	defer rc.Close()

	reader := bufio.NewReader(rc)
	for {
		line, err := reader.ReadString('\n')
		if line != `` {
			s.writeLine(state, podName, containerName, multiContainer, line)
		}
		if err != nil {
			s.mu.Lock()
			state.active = false
			state.ended = errors.Is(err, io.EOF)
			s.mu.Unlock()
			return
		}
	}

}

// writeLine strips the kubelet timestamp, drops lines already written by an earlier stream of
// the same container instance and writes the line, prefixed with the pod (and container) name
// once more than one pod has been seen.
func (s *logStreamer) writeLine(state *instanceState, podName, containerName string, multiContainer bool, line string) {

	text := line
	if i := strings.IndexByte(line, ' '); i > 0 {
		if ts, err := time.Parse(time.RFC3339Nano, line[:i]); err == nil {
			text = line[i+1:]
			s.mu.Lock()
			duplicate := !state.last.IsZero() && !ts.After(state.last)
			if !duplicate {
				state.last = ts
			}
			s.mu.Unlock()
			if duplicate {
				return
			}
		}
	}

	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}

	s.mu.Lock()
	multiPod := len(s.pods) > 1
	s.mu.Unlock()

	switch {
	case multiContainer:
		text = fmt.Sprintf("[%s/%s] %s", podName, containerName, text)
	case multiPod:
		text = fmt.Sprintf("[%s] %s", podName, text)
	}

	io.WriteString(s.out, text)

}

// reportWaiting reports, once per reason, a container that cannot start (for example ImagePullBackOff).
func (s *logStreamer) reportWaiting(podName, containerName string, waiting *corev1.ContainerStateWaiting) {

	if benignWaitingReasons[waiting.Reason] {
		return
	}

	key := podName + `/` + containerName + `/` + waiting.Reason
	s.mu.Lock()
	seen := s.reported[key]
	s.reported[key] = true
	s.mu.Unlock()

	if !seen {
		fmt.Fprintf(s.warn, "eksjob: pod %s container %s is waiting: %s: %s\n", podName, containerName, waiting.Reason, waiting.Message)
	}

}
