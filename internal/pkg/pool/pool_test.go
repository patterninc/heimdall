package pool

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestStopClaimsNothingMoreAndWaitFinishesClaimedWork(t *testing.T) {

	var claims, finished atomic.Int32
	started := make(chan struct{}, 2)
	release := make(chan struct{})

	p := &Pool[int]{Size: 2}
	getWork := func(limit int) ([]int, error) {
		if claims.Add(1) == 1 {
			return []int{1, 2}, nil
		}
		return nil, nil
	}
	worker := func(_ context.Context, _ int) error {
		started <- struct{}{}
		<-release
		finished.Add(1)
		return nil
	}

	if err := p.Start(worker, getWork); err != nil {
		t.Fatal(err)
	}
	<-started
	<-started

	p.Stop()
	time.Sleep(50 * time.Millisecond)
	claimsAfterStop := claims.Load()
	time.Sleep(4 * minimumSleep * time.Millisecond)
	if got := claims.Load(); got != claimsAfterStop {
		t.Fatalf("pool kept claiming after Stop: %d claims, then %d", claimsAfterStop, got)
	}

	close(release)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := p.Wait(ctx); err != nil {
		t.Fatalf("Wait: %v", err)
	}
	if got := finished.Load(); got != 2 {
		t.Fatalf("finished %d jobs, want 2", got)
	}

}

func TestWaitStopsAtDeadline(t *testing.T) {

	started := make(chan struct{})
	release := make(chan struct{})
	defer close(release)

	p := &Pool[int]{Size: 1}
	var claimed atomic.Bool
	getWork := func(limit int) ([]int, error) {
		if claimed.CompareAndSwap(false, true) {
			return []int{1}, nil
		}
		return nil, nil
	}
	worker := func(_ context.Context, _ int) error {
		close(started)
		<-release
		return nil
	}

	if err := p.Start(worker, getWork); err != nil {
		t.Fatal(err)
	}
	<-started

	p.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := p.Wait(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Wait: got %v, want deadline exceeded", err)
	}

}
