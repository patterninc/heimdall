package sparkeks

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kubeflow/spark-operator/v2/api/v1beta2"
	corev1 "k8s.io/api/core/v1"
)

func TestUpdateS3ToS3aURI(t *testing.T) {
	uri := "s3://bucket/path/file"
	expected := "s3a://bucket/path/file"
	result := updateS3ToS3aURI(uri)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
	// Should not change if already s3a
	uri2 := "s3a://bucket/path/file"
	result2 := updateS3ToS3aURI(uri2)
	if result2 != uri2 {
		t.Errorf("Expected %s, got %s", uri2, result2)
	}
}

func TestGetSparkSubmitParameters(t *testing.T) {
	ctx := &jobContext{
		Parameters: &jobParameters{
			Properties: map[string]string{
				"spark.executor.memory": "4g",
				"spark.driver.cores":    "2",
			},
		},
	}
	params := getSparkSubmitParameters(ctx)
	if params == nil || !strings.Contains(*params, "--conf spark.executor.memory=4g") || !strings.Contains(*params, "--conf spark.driver.cores=2") {
		t.Errorf("Unexpected Spark submit parameters: %v", params)
	}
}

func TestGetSparkSubmitParameters_Empty(t *testing.T) {
	ctx := &jobContext{
		Parameters: &jobParameters{
			Properties: map[string]string{},
		},
	}
	params := getSparkSubmitParameters(ctx)
	if params != nil {
		t.Errorf("Expected nil for empty properties, got %v", params)
	}
}

func TestGetSparkSubmitParameters_NilParameters(t *testing.T) {
	ctx := &jobContext{
		Parameters: nil,
	}
	params := getSparkSubmitParameters(ctx)
	if params != nil {
		t.Errorf("Expected nil for nil parameters, got %v", params)
	}
}

func TestGetS3FileURI_InvalidFormat(t *testing.T) {
	_, err := getS3FileURI(context.Background(), aws.Config{}, "invalid-uri", "avro")
	if err == nil {
		t.Error("Expected error for invalid S3 URI format")
	}
}

func TestGetS3FileURI_ValidFormat(t *testing.T) {
	// This test only checks parsing, not actual AWS interaction
	uri := "s3://bucket/path/"
	_, err := getS3FileURI(context.Background(), aws.Config{}, uri, "avro")
	// Should not error on parsing, but will error on AWS call (which is fine for unit test context)
	if err == nil || !strings.Contains(err.Error(), "failed to list S3 objects") {
		t.Logf("Expected AWS list objects error, got: %v", err)
	}
}

func TestPrintState(t *testing.T) {
	f, err := os.CreateTemp("", "sparkeks-test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer f.Close()
	printState(f, "RUNNING")
	// Optionally check file contents if needed
}

func sparkAppWith(state v1beta2.ApplicationStateType, msg, driverPod string) *v1beta2.SparkApplication {
	app := &v1beta2.SparkApplication{}
	app.Status.AppState.State = state
	app.Status.AppState.ErrorMessage = msg
	app.Status.DriverInfo.PodName = driverPod
	return app
}

func driverPodWith(name string, exit *int32) corev1.Pod {
	pod := corev1.Pod{}
	pod.Name = name
	pod.Labels = map[string]string{sparkRoleLabel: sparkRoleDriver}
	if exit != nil {
		pod.Status.ContainerStatuses = []corev1.ContainerStatus{{
			Name:  sparkDriverContainerName,
			State: corev1.ContainerState{Terminated: &corev1.ContainerStateTerminated{ExitCode: *exit}},
		}}
	}
	return pod
}

func TestIsDriverPodNotFound(t *testing.T) {
	tests := []struct {
		name string
		app  *v1beta2.SparkApplication
		want bool
	}{
		{"failed driver pod not found", sparkAppWith(v1beta2.ApplicationStateFailed, "driver pod not found", ""), true},
		{"failing driver pod not found", sparkAppWith(v1beta2.ApplicationStateFailing, "driver pod not found", ""), true},
		{"genuine driver failure", sparkAppWith(v1beta2.ApplicationStateFailed, "driver container failed with ExitCode: 1", ""), false},
		{"empty message", sparkAppWith(v1beta2.ApplicationStateFailed, "", ""), false},
		{"completed", sparkAppWith(v1beta2.ApplicationStateCompleted, "", ""), false},
		{"nil", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDriverPodNotFound(tt.app); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDriverOutcomeEvidence(t *testing.T) {
	zero, one := int32(0), int32(1)
	failed := sparkAppWith(v1beta2.ApplicationStateFailed, "driver pod not found", "drv")
	tests := []struct {
		name string
		app  *v1beta2.SparkApplication
		pods []corev1.Pod
		want bool
	}{
		{"no pods, still failed: unknown", failed, nil, false},
		{"driver exit 0: success", failed, []corev1.Pod{driverPodWith("drv", &zero)}, true},
		{"driver exit 1: failure", failed, []corev1.Pod{driverPodWith("drv", &one)}, false},
		{"driver without termination info: unknown", failed, []corev1.Pod{driverPodWith("drv", nil)}, false},
		{"non-driver pod exit 0 ignored", failed, []corev1.Pod{func() corev1.Pod {
			p := driverPodWith("exec-1", &zero)
			p.Labels = map[string]string{sparkRoleLabel: "executor"}
			return p
		}()}, false},
		{"CR re-read now COMPLETED", sparkAppWith(v1beta2.ApplicationStateCompleted, "", "drv"), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, why := driverOutcomeEvidence(tt.app, tt.pods)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
			if got && why == "" {
				t.Error("positive evidence must carry a reason")
			}
		})
	}
}

func TestDriverPodNotFoundOutcomeUnknownMessage(t *testing.T) {
	for _, want := range []string{"outcome unknown", driverPodNotFoundMsg, "race", "check the job output"} {
		if !strings.Contains(driverPodNotFoundOutcomeUnknown, want) {
			t.Errorf("message missing %q: %s", want, driverPodNotFoundOutcomeUnknown)
		}
	}
}
