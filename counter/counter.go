package counter

import (
	"sync"
	"time"
)

type Counters[K comparable] struct {
	lock           *sync.Mutex
	counters       map[K]int
	batch_size     int
	cpt            int
	batch_duration time.Duration
	harvester      Harvester[K]
	timer          *time.Timer
	batch_complete chan interface{}
}

type Harvester[K comparable] func(key []K, value []int) error

func New[K comparable](batch_size int, batch_duration time.Duration, h Harvester[K]) *Counters[K] {
	// New return a new *Counters[K]
	// if batch_size > 0, it's a trigger
	c := &Counters[K]{
		lock:           &sync.Mutex{},
		counters:       make(map[K]int),
		batch_size:     batch_size,
		batch_duration: batch_duration,
		harvester:      h,
		timer:          time.NewTimer(batch_duration),
		batch_complete: make(chan interface{}),
	}
	if c.harvester != nil {
		go c.loopForHarvest()
	}
	return c
}

func (c *Counters[K]) loopForHarvest() {
	for {
		select {
		case <-c.timer.C:
			c.lock.Lock()
		case <-c.batch_complete:
			// lock is Lock in the Add function
		}
		c.harvest()
		c.timer = time.NewTimer(c.batch_duration)
		c.lock.Unlock()
	}
}

func (c *Counters[K]) Add(key K, value int) (bool, error) {
	c.lock.Lock()
	r := false
	v, ok := c.counters[key]
	if !ok {
		c.counters[key] = value
	}
	c.counters[key] = v + value
	c.cpt++
	if c.cpt == c.batch_size { // increment is done before, for handling early return
		c.batch_complete <- new(interface{})
		r = true
	} else {
		c.lock.Unlock()
	}
	return r, nil
}

func (c *Counters[K]) harvest() error {
	var err error
	keys := make([]K, len(c.counters))
	values := make([]int, len(c.counters))
	i := 0
	for k, v := range c.counters {
		keys[i] = k
		values[i] = v
		i++
	}
	err = c.harvester(keys, values)
	if err != nil {
		return err
	}
	c.counters = make(map[K]int)
	c.cpt = 0
	return nil
}

func (c *Counters[K]) Harvest() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.harvest()
}
