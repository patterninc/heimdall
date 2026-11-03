package heimdall

import (
	"context"
	"sync"
)

// background tracks goroutines that outlive the request or job that started them,
// so shutdown can wait for them. The zero value is ready to use.
type background struct {
	mu      sync.Mutex
	running int
	idle    chan struct{}
}

func (b *background) Go(f func()) {

	b.mu.Lock()
	b.running++
	b.mu.Unlock()

	go func() {
		defer b.done()
		f()
	}()

}

func (b *background) done() {

	b.mu.Lock()
	defer b.mu.Unlock()

	b.running--
	if b.running == 0 && b.idle != nil {
		close(b.idle)
		b.idle = nil
	}

}

// Wait blocks until no tracked goroutines are running, or ctx is done.
func (b *background) Wait(ctx context.Context) error {

	b.mu.Lock()
	if b.running == 0 {
		b.mu.Unlock()
		return nil
	}
	if b.idle == nil {
		b.idle = make(chan struct{})
	}
	idle := b.idle
	b.mu.Unlock()

	select {
	case <-idle:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

}
