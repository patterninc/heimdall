package pool

import (
	"context"
	"fmt"
	"time"
)

const (
	defaultPoolSize = 10
	minimumSleep    = 250 // ms
)

type Pool[T any] struct {
	Size  int `yaml:"size,omitempty" json:"size,omitempty"`
	Sleep int `yaml:"sleep,omitempty" json:"sleep,omitempty"`
	queue chan T
}

func (p *Pool[T]) Start(ctx context.Context, worker func(context.Context, T) error, getWork func(int) ([]T, error)) error {

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

	// let's set the counter
	tokens := &counter{}

	// let's provision workers
	for i := 0; i < p.Size; i++ {

		go func(_ int, c *counter) {

			for {

				// we're ready to work on a single work item
				c.Add(1)

				// let's wait and get one item to work on...
				select {
				case <-ctx.Done():
					return
				case w, ok := <-p.queue:
					if !ok {
						return
					}

					// do the work....
					if err := worker(ctx, w); err != nil {
						// TODO: implement proper error logging
						fmt.Println(`worker:`, err)
					}
				}

			}

		}(i, tokens)

	}

	// let's provision process that will add work to the queue
	go func(c *counter) {
		defer close(p.queue)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			// how much work are we ready to request for our pool?
			limit := c.Get()

			// if all our workers are busy, we're waiting...
			if limit == 0 {
				time.Sleep(time.Duration(p.Sleep) * time.Millisecond)
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
				time.Sleep(time.Duration(p.Sleep) * time.Millisecond)
				continue
			}

			// let's add the work we got to the queue...
			for _, item := range items {
				select {
				case <-ctx.Done():
					return
				case p.queue <- item:
				}
			}

			// substruct the number of the items from our counter
			tokens.Add(-itemsCount)

		}

	}(tokens)

	return nil
}
