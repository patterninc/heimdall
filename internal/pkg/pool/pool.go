package pool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	defaultPoolSize = 10
	minimumSleep    = 250 // ms
)

type Pool[T any] struct {
	Size       int `yaml:"size,omitempty" json:"size,omitempty"`
	Sleep      int `yaml:"sleep,omitempty" json:"sleep,omitempty"`
	queue      chan T
	stop       chan struct{}
	feederDone chan struct{}
	workers    sync.WaitGroup
}

func (p *Pool[T]) Start(worker func(context.Context, T) error, getWork func(int) ([]T, error)) error {

	// do we have the size set?
	if p.Size <= 0 {
		p.Size = defaultPoolSize
	}

	// what about sleep setting?
	if p.Sleep <= minimumSleep {
		p.Sleep = minimumSleep
	}

	// set the queue of the size of double the pool size
	p.queue = make(chan T, p.Size*2)
	p.stop = make(chan struct{})
	p.feederDone = make(chan struct{})

	// let's set the counter
	tokens := &counter{}

	// let's provision workers
	p.workers.Add(p.Size)
	for i := 0; i < p.Size; i++ {

		go func(_ int, c *counter) {

			defer p.workers.Done()

			for {

				// we're ready to work on a single work item
				c.Add(1)

				// let's wait and get one item to work on...
				w, ok := <-p.queue

				// if we're not having the active queue, quit...
				if !ok {
					break
				}

				// do the work....
				err := worker(context.Background(), w)

				if err != nil {
					// TODO: implement proper error logging
					fmt.Println(`worker:`, err)
				}

			}

		}(i, tokens)

	}

	// let's provision process that will add work to the queue
	go func(c *counter) {

		defer close(p.feederDone)

		for {

			// once stopped, claim nothing more
			select {
			case <-p.stop:
				return
			default:
			}

			// how much work are we ready to request for our pool?
			limit := c.Get()

			// if all our workers are busy, we're waiting...
			if limit == 0 {
				p.sleep()
				continue
			}

			// let's get the work up to the limit...
			items, err := getWork(limit)
			itemsCount := len(items)

			if err != nil {
				// TODO: implement proper error logging
				fmt.Println(`getWork:`, err)
			}

			// if we did not get any work, keep waiting...
			if itemsCount == 0 {
				p.sleep()
				continue
			}

			// let's add the work we got to the queue...
			for _, item := range items {
				p.queue <- item
			}

			// substruct the number of the items from our counter
			tokens.Add(-itemsCount)

		}

	}(tokens)

	return nil
}

// Stop makes the pool claim no more work. Work already claimed keeps running; use Wait to wait for it.
func (p *Pool[T]) Stop() {

	close(p.stop)

}

// Wait blocks until the workers have finished all claimed work, or ctx is done. Call Stop first.
func (p *Pool[T]) Wait(ctx context.Context) error {

	// the queue can only be closed once the feeder can no longer send to it
	select {
	case <-p.feederDone:
	case <-ctx.Done():
		return ctx.Err()
	}
	close(p.queue)

	done := make(chan struct{})
	go func() {
		p.workers.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

}

func (p *Pool[T]) sleep() {

	select {
	case <-p.stop:
	case <-time.After(time.Duration(p.Sleep) * time.Millisecond):
	}

}
