package ecs

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	smithy "github.com/aws/smithy-go"
)

func TestPartialStoppedDescribeKeepsPolling(t *testing.T) {
	// Two active tasks, only one STOPPED task returned. The omitted task must keep polling.
	if describedTaskCountMatches(2, 1) {
		t.Fatal("a partial DescribeTasks response must not finish polling")
	}
	if !describedTaskCountMatches(2, 2) {
		t.Fatal("polling can finish once every active task is returned")
	}
}

func TestIsTaskSuccessful(t *testing.T) {
	execCtx := &executionContext{TaskDefinitionWrapper: &taskDefinitionWrapper{
		EssentialContainers: map[string]bool{"main": true},
	}}
	zero, failed, oom := int32(0), int32(1), int32(137)

	if isTaskSuccessful(types.Task{Containers: []types.Container{{Name: aws.String("main")}}}, execCtx) {
		t.Fatal("a container that never exited is not success")
	}
	if isTaskSuccessful(types.Task{Containers: []types.Container{{Name: aws.String("main"), ExitCode: &failed}}}, execCtx) {
		t.Fatal("nonzero exit is not success")
	}
	if isTaskSuccessful(types.Task{Containers: []types.Container{{Name: aws.String("main"), ExitCode: &oom}}}, execCtx) {
		t.Fatal("exit 137 is OOM and is not success")
	}
	if !isTaskSuccessful(types.Task{Containers: []types.Container{{Name: aws.String("main"), ExitCode: &zero}}}, execCtx) {
		t.Fatal("zero exit is success")
	}
}

func TestRetryWithBackoff_Throttling(t *testing.T) {
	calls := 0
	throttleErr := &smithy.GenericAPIError{Code: "ThrottlingException", Message: "Rate exceeded"}

	err := retryWithBackoff(context.Background(), 3, func() error {
		calls++
		return throttleErr
	})

	if err != throttleErr {
		t.Fatalf("expected throttle error, got %v", err)
	}
	if calls != 4 { // 1 initial + 3 retries
		t.Fatalf("expected 4 calls, got %d", calls)
	}
}

func TestRetryWithBackoff_NonThrottling(t *testing.T) {
	calls := 0
	otherErr := &smithy.GenericAPIError{Code: "AccessDeniedException", Message: "denied"}

	retryWithBackoff(context.Background(), 3, func() error {
		calls++
		return otherErr
	})

	if calls != 1 { // should not retry
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetryWithBackoff_SuccessAfterThrottle(t *testing.T) {
	calls := 0
	throttleErr := &smithy.GenericAPIError{Code: "ThrottlingException", Message: "Rate exceeded"}

	err := retryWithBackoff(context.Background(), 3, func() error {
		calls++
		if calls < 3 {
			return throttleErr
		}
		return nil
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRetryWithBackoff_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	calls := 0
	throttleErr := &smithy.GenericAPIError{Code: "ThrottlingException", Message: "Rate exceeded"}

	err := retryWithBackoff(ctx, 3, func() error {
		calls++
		return throttleErr
	})

	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}
