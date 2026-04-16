package ecs

import (
	"context"
	"testing"

	smithy "github.com/aws/smithy-go"
)

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
