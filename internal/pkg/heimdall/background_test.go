package heimdall

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestBackgroundWaitReturnsWhenIdle(t *testing.T) {

	var b background
	if err := b.Wait(context.Background()); err != nil {
		t.Fatalf("Wait: %v", err)
	}

}

func TestBackgroundWaitWaitsForRunningWork(t *testing.T) {

	var b background
	var done atomic.Int32
	for range 3 {
		b.Go(func() {
			time.Sleep(50 * time.Millisecond)
			done.Add(1)
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := b.Wait(ctx); err != nil {
		t.Fatalf("Wait: %v", err)
	}
	if got := done.Load(); got != 3 {
		t.Fatalf("Wait returned with %d of 3 done", got)
	}

}

func TestBackgroundWaitStopsAtDeadline(t *testing.T) {

	var b background
	release := make(chan struct{})
	defer close(release)
	b.Go(func() { <-release })

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if err := b.Wait(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Wait: got %v, want deadline exceeded", err)
	}

}
