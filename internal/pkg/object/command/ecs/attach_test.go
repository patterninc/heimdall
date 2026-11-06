package ecs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"

	"github.com/patterninc/heimdall/pkg/plugin"
)

// fakeTask is a task the fake ECS endpoint knows about
type fakeTask struct {
	arn, startedBy, status string
	createdAt              float64
}

// fakeECS answers ListTasks (filtered by startedBy and desired status) and DescribeTasks for the given tasks
func fakeECS(t *testing.T, tasks []fakeTask) *ecs.Client {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
						})
					}
				}
			}
			output = map[string]any{"tasks": described}
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
	})

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
			execCtx := &executionContext{
				TaskCount:     tt.taskCount,
				ClusterConfig: &clusterContext{ClusterName: "cluster"},
				ecsClient:     fakeECS(t, tt.tasks),
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
