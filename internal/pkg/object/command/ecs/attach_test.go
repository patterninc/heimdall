package ecs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"

	"github.com/patterninc/heimdall/pkg/duration"
	"github.com/patterninc/heimdall/pkg/plugin"
)

// fakeTask is a task the fake ECS endpoint knows about; a STOPPED task's container exited with exitCode
type fakeTask struct {
	arn, startedBy, status string
	createdAt              float64
	exitCode               int
}

// fakeCalls records RunTask and StopTask calls; tasks started by RunTask stop at once with the next runExitCodes entry
type fakeCalls struct {
	mu           sync.Mutex
	runExitCodes []int
	clientTokens []string
	stopped      []string
}

// fakeECS answers ListTasks (filtered by startedBy and desired status), DescribeTasks, RunTask and StopTask
func fakeECS(t *testing.T, tasks []fakeTask) (*ecs.Client, *fakeCalls) {

	calls := &fakeCalls{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.mu.Lock()
		defer calls.mu.Unlock()

		var input map[string]any
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			t.Errorf("decode request: %v", err)
		}

		var output map[string]any
		switch target := r.Header.Get("X-Amz-Target"); {
		case strings.HasSuffix(target, ".ListTasks"):
			arns := []string{}
			for _, task := range tasks {
				if task.startedBy == input["startedBy"] && task.status == input["desiredStatus"] {
					arns = append(arns, task.arn)
				}
			}
			output = map[string]any{"taskArns": arns}
		case strings.HasSuffix(target, ".DescribeTasks"):
			described := []map[string]any{}
			for _, arn := range input["tasks"].([]any) {
				for _, task := range tasks {
					if task.arn == arn {
						described = append(described, map[string]any{
							"taskArn": task.arn, "startedBy": task.startedBy, "lastStatus": task.status, "createdAt": task.createdAt,
							"containers": []map[string]any{{"name": "main", "exitCode": task.exitCode}},
						})
					}
				}
			}
			output = map[string]any{"tasks": described}
		case strings.HasSuffix(target, ".RunTask"):
			exitCode := 0
			if len(calls.runExitCodes) > 0 {
				exitCode, calls.runExitCodes = calls.runExitCodes[0], calls.runExitCodes[1:]
			}
			calls.clientTokens = append(calls.clientTokens, input["clientToken"].(string))
			task := fakeTask{arn: "run-" + input["clientToken"].(string), startedBy: input["startedBy"].(string), status: "STOPPED", createdAt: float64(1000 + len(tasks)), exitCode: exitCode}
			tasks = append(tasks, task)
			output = map[string]any{"tasks": []map[string]any{{"taskArn": task.arn}}}
		case strings.HasSuffix(target, ".StopTask"):
			calls.stopped = append(calls.stopped, input["task"].(string))
			output = map[string]any{}
		default:
			t.Errorf("unexpected ECS call %s", target)
		}

		w.Header().Set("Content-Type", "application/x-amz-json-1.1")
		json.NewEncoder(w).Encode(output)
	}))
	t.Cleanup(server.Close)

	return ecs.New(ecs.Options{
		Region:       "us-west-2",
		BaseEndpoint: aws.String(server.URL),
		Credentials:  aws.AnonymousCredentials{},
	}), calls

}

func TestAttachTasks(t *testing.T) {

	const job = "job1"
	name := func(n string) string { return startedByPrefix + job + "-" + n }

	tests := []struct {
		name        string
		taskCount   int
		tasks       []fakeTask
		wantARNs    map[string]string
		wantRetries map[string]int
		wantErr     bool
	}{
		{
			name:      "running task is attached even when a newer stopped one exists",
			taskCount: 1,
			tasks: []fakeTask{
				{arn: "failed", startedBy: name("0"), status: "STOPPED", createdAt: 300},
				{arn: "running", startedBy: name("0"), status: "RUNNING", createdAt: 200},
			},
			wantARNs:    map[string]string{name("0"): "running"},
			wantRetries: map[string]int{name("0"): 1},
		},
		{
			name:      "only stopped tasks: the newest is attached and earlier ones count as retries",
			taskCount: 1,
			tasks: []fakeTask{
				{arn: "attempt0", startedBy: name("0"), status: "STOPPED", createdAt: 100},
				{arn: "attempt2", startedBy: name("0"), status: "STOPPED", createdAt: 300},
				{arn: "attempt1", startedBy: name("0"), status: "STOPPED", createdAt: 200},
			},
			wantARNs:    map[string]string{name("0"): "attempt2"},
			wantRetries: map[string]int{name("0"): 2},
		},
		{
			name:      "each task number is attached on its own",
			taskCount: 2,
			tasks: []fakeTask{
				{arn: "zero", startedBy: name("0"), status: "RUNNING", createdAt: 100},
				{arn: "one-failed", startedBy: name("1"), status: "STOPPED", createdAt: 100},
				{arn: "one", startedBy: name("1"), status: "RUNNING", createdAt: 200},
			},
			wantARNs:    map[string]string{name("0"): "zero", name("1"): "one"},
			wantRetries: map[string]int{name("0"): 0, name("1"): 1},
		},
		{
			name:      "a missing task number fails the resume",
			taskCount: 2,
			tasks: []fakeTask{
				{arn: "zero", startedBy: name("0"), status: "RUNNING", createdAt: 100},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, err := os.CreateTemp(t.TempDir(), "stdout")
			if err != nil {
				t.Fatal(err)
			}
			client, _ := fakeECS(t, tt.tasks)
			execCtx := &executionContext{
				TaskCount:     tt.taskCount,
				ClusterConfig: &clusterContext{ClusterName: "cluster"},
				ecsClient:     client,
				taskDefARN:    aws.String("task-definition"),
				tasks:         make(map[string]*taskTracker),
				runtime:       &plugin.Runtime{Stdout: stdout},
			}

			err = execCtx.attachTasks(context.Background(), job)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}

			if len(execCtx.tasks) != len(tt.wantARNs) {
				t.Fatalf("tracked %d tasks, want %d", len(execCtx.tasks), len(tt.wantARNs))
			}
			for taskName, wantARN := range tt.wantARNs {
				tracker := execCtx.tasks[taskName]
				if tracker == nil {
					t.Fatalf("no tracker for %s", taskName)
				}
				if tracker.ActiveARN != wantARN {
					t.Errorf("%s: attached %s, want %s", taskName, tracker.ActiveARN, wantARN)
				}
				if tracker.Retries != tt.wantRetries[taskName] {
					t.Errorf("%s: retries %d, want %d", taskName, tracker.Retries, tt.wantRetries[taskName])
				}
			}
		})
	}

}

// TestResumeThenPoll covers a task that stopped while no agent was watching: it is attached, then polled
// as Execute would, so a failure is restarted or gives up at max_fail_count across agents
func TestResumeThenPoll(t *testing.T) {

	const job = "job1"
	name := startedByPrefix + job + "-0"
	other := startedByPrefix + job + "-1"

	tests := []struct {
		name         string
		taskCount    int
		tasks        []fakeTask
		runExitCodes []int
		wantTokens   []string
		wantFailed   bool
		wantStopped  []string
	}{
		{
			name:  "task that succeeded unwatched completes without a restart",
			tasks: []fakeTask{{arn: "done", startedBy: name, status: "STOPPED", createdAt: 100, exitCode: 0}},
		},
		{
			name:       "task that failed unwatched is restarted with the next unused client token",
			tasks:      []fakeTask{{arn: "attempt0", startedBy: name, status: "STOPPED", createdAt: 100, exitCode: 1}, {arn: "attempt1", startedBy: name, status: "STOPPED", createdAt: 200, exitCode: 1}},
			wantTokens: []string{name + "-2"},
		},
		{
			name:         "restarts keep failing until max_fail_count",
			tasks:        []fakeTask{{arn: "attempt0", startedBy: name, status: "STOPPED", createdAt: 100, exitCode: 1}},
			runExitCodes: []int{1, 1},
			wantTokens:   []string{name + "-1", name + "-2"},
			wantFailed:   true,
		},
		{
			name:       "failures by the earlier agent count toward max_fail_count",
			tasks:      []fakeTask{{arn: "attempt0", startedBy: name, status: "STOPPED", createdAt: 100, exitCode: 1}, {arn: "attempt1", startedBy: name, status: "STOPPED", createdAt: 200, exitCode: 1}, {arn: "attempt2", startedBy: name, status: "STOPPED", createdAt: 300, exitCode: 1}},
			wantFailed: true,
		},
		{
			name:        "a task giving up after resume stops the other task numbers",
			taskCount:   2,
			tasks:       []fakeTask{{arn: "attempt0", startedBy: name, status: "STOPPED", createdAt: 100, exitCode: 1}, {arn: "attempt1", startedBy: name, status: "STOPPED", createdAt: 200, exitCode: 1}, {arn: "attempt2", startedBy: name, status: "STOPPED", createdAt: 300, exitCode: 1}, {arn: "still-running", startedBy: other, status: "RUNNING", createdAt: 100}},
			wantFailed:  true,
			wantStopped: []string{"still-running"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, err := os.CreateTemp(t.TempDir(), "stdout")
			if err != nil {
				t.Fatal(err)
			}
			if tt.taskCount == 0 {
				tt.taskCount = 1
			}
			client, calls := fakeECS(t, tt.tasks)
			calls.runExitCodes = tt.runExitCodes
			execCtx := &executionContext{
				TaskCount:             tt.taskCount,
				ClusterConfig:         &clusterContext{ClusterName: "cluster", LaunchType: "FARGATE", VPCConfig: vpcConfig{Subnets: []string{"subnet"}}},
				TaskDefinitionWrapper: &taskDefinitionWrapper{EssentialContainers: map[string]bool{"main": true}},
				PollingInterval:       duration.Duration(time.Millisecond),
				Timeout:               duration.Duration(time.Minute),
				MaxFailCount:          3,
				ecsClient:             client,
				taskDefARN:            aws.String("task-definition"),
				tasks:                 make(map[string]*taskTracker),
				runtime:               &plugin.Runtime{Stdout: stdout},
			}

			if err := execCtx.attachTasks(context.Background(), job); err != nil {
				t.Fatal(err)
			}
			if err := execCtx.pollForCompletion(context.Background()); err != nil {
				t.Fatal(err)
			}

			if strings.Join(calls.clientTokens, ",") != strings.Join(tt.wantTokens, ",") {
				t.Errorf("RunTask client tokens %v, want %v", calls.clientTokens, tt.wantTokens)
			}
			for _, arn := range tt.wantStopped {
				if !strings.Contains(strings.Join(calls.stopped, ","), arn) {
					t.Errorf("StopTask calls %v, want %s stopped", calls.stopped, arn)
				}
			}
			if failed := execCtx.failureError != nil; failed != tt.wantFailed {
				t.Errorf("failed %v (%v), want %v", failed, execCtx.failureError, tt.wantFailed)
			}
			if completed := execCtx.tasks[name].Completed; completed == tt.wantFailed {
				t.Errorf("completed %v, want %v", completed, !tt.wantFailed)
			}
		})
	}

}
