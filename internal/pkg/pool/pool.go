package pool

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/patterninc/heimdall/internal/pkg/database"
	"github.com/patterninc/heimdall/pkg/object/job"
)

const (
	defaultPoolSize = 10
	minimumSleep    = 250 // ms
)

type Pool[T any] struct {
	Size  int `yaml:"size,omitempty" json:"size,omitempty"`
	Sleep int `yaml:"sleep,omitempty" json:"sleep,omitempty"`
	queue chan *job.Job

	runningJobs    map[string]context.CancelFunc
	runningJobsMux sync.RWMutex
	db             *database.Database
}

func (p *Pool[T]) Start(worker func(context.Context, *job.Job) error, getWork func(int) ([]*job.Job, error), database *database.Database) error {

	// record database context
	p.db = database

	// do we have the size set?
	if p.Size <= 0 {
		p.Size = defaultPoolSize
	}

	// what about sleep setting?
	if p.Sleep <= minimumSleep {
		p.Sleep = minimumSleep
	}

	// set the queue of the size of double the pool size
	p.queue = make(chan *job.Job, p.Size*2)

	// let's set the counter
	tokens := &counter{}

	// Initialize tracking
	p.runningJobs = make(map[string]context.CancelFunc)

	// Start cancellation polling loop
	go p.pollForCancellations()

	// let's provision workers
	for i := 0; i < p.Size; i++ {

		go func(_ int, c *counter) {

			for {

				// we're ready to work on a single work item
				c.Add(1)

				// let's wait and get one item to work on...
				w, ok := <-p.queue

				// if we're not having the active queue, quit...
				if !ok {
					break
				}

				ctx, cancel := context.WithCancel(context.Background())

				// register job as running
				p.registerRunningJob(w.ID, cancel)

				// do the work....
				err := worker(ctx, w)

				// remove job from running jobs
				p.unregisterRunningJob(w.ID)

				if err != nil {
					// TODO: implement proper error logging
					fmt.Println(`worker:`, err)
				}

			}

		}(i, tokens)

	}

	// let's provision process that will add work to the queue
	go func(c *counter) {

		for {

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
				p.queue <- item
			}

			// substruct the number of the items from our counter
			tokens.Add(-itemsCount)

		}

	}(tokens)

	return nil
}

func (p *Pool[T]) pollForCancellations() {
	// let's poll for cancellations every 15 seconds
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {

		// Get jobs in CANCELLING state from database
		cancellingJobs := p.getCancellingJobs()

		// Check each cancelling job
		for _, jobID := range cancellingJobs {
			if cancelFunc, isLocal := p.isJobRunningLocally(jobID); isLocal {
				cancelFunc() // Trigger context cancellation

				// Update job status to CANCELLED in database
				p.updateJobStatusToCancelled(jobID)
			}
		}
	}
}

func (p *Pool[T]) registerRunningJob(jobID string, cancel context.CancelFunc) {

	p.runningJobsMux.Lock()
	defer p.runningJobsMux.Unlock()

	p.runningJobs[jobID] = cancel

}

func (p *Pool[T]) unregisterRunningJob(jobID string) {

	p.runningJobsMux.Lock()
	defer p.runningJobsMux.Unlock()

	delete(p.runningJobs, jobID)
}

// Check if a job is running locally
func (p *Pool[T]) isJobRunningLocally(jobID string) (context.CancelFunc, bool) {

	p.runningJobsMux.RLock()
	defer p.runningJobsMux.RUnlock()

	cancelFunc, exists := p.runningJobs[jobID]
	return cancelFunc, exists
}
