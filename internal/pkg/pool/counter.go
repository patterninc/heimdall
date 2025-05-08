package pool

import (
	"sync"
)

type counter struct {
	c int
	sync.Mutex
}

func (c *counter) Add(i int) {

	c.Lock()
	defer c.Unlock()

	c.c += i

}

func (c *counter) Get() int {

	c.Lock()
	defer c.Unlock()

	return c.c

}
