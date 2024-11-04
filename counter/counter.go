package counter

import "sync"

type Counters[K comparable] struct {
	lock       *sync.Mutex
	counters   map[K]int
	batch_size int
	cpt        int
	harvester  Harvester[K]
}

type Harvester[K comparable] func(key []K, value []int) error

func New[K comparable](batch_size int, h Harvester[K]) *Counters[K] {
	// New return a new *Counters[K]
	// if batch_size > 0, it's a trigger
	return &Counters[K]{
		lock:       &sync.Mutex{},
		counters:   make(map[K]int),
		batch_size: batch_size,
		harvester:  h,
	}
}

func (c *Counters[K]) Add(key K, value int) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	r := false
	v, ok := c.counters[key]
	if !ok {
		c.counters[key] = value
		return r, nil
	}
	c.counters[key] = v + value
	c.cpt++
	if c.cpt == c.batch_size-1 { // increment is done before, for handling early return
		r = true
		c.cpt = 0
		if c.harvester != nil {
			err := c.harvest()
			if err != nil {
				return false, err
			}
		}
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
	return nil
}

func (c *Counters[K]) Harvest() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.harvest()
}
