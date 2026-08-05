package counter

import (
	"log/slog"
	"sync"
	"time"
)

// Counters count stuff and trigger harvesters periodically.
type Counters[K comparable] struct {
	lock           *sync.Mutex
	counters       map[K]int
	batch_size     int
	cpt            int
	batch_duration time.Duration
	harvester      Harvester[K]
	timer          *time.Timer
}

type Harvester[K comparable] func(key []K, value []int) error

// New returns a new *Counters[K].
// If batch_size > 0, reaching it triggers an immediate harvest. The timer is
// a fallback that harvests whatever accumulated after batch_duration.
func New[K comparable](batch_size int, batch_duration time.Duration, h Harvester[K]) *Counters[K] {
	c := &Counters[K]{
		lock:           &sync.Mutex{},
		counters:       make(map[K]int),
		batch_size:     batch_size,
		batch_duration: batch_duration,
		harvester:      h,
		timer:          time.NewTimer(batch_duration),
	}
	if c.harvester != nil {
		go c.loopForHarvest()
	}
	return c
}

func (c *Counters[K]) loopForHarvest() {
	for {
		<-c.timer.C
		c.lock.Lock()
		if err := c.harvest(); err != nil {
			// Keep the counters: the next tick will retry.
			slog.Error("harvest failed", "err", err)
		}
		c.timer = time.NewTimer(c.batch_duration)
		c.lock.Unlock()
	}
}

// Add stores value under key. When the batch is full, the accumulated
// counters are harvested synchronously under the lock: batches are exact
// and the lock never travels across goroutines.
func (c *Counters[K]) Add(key K, value int) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.counters[key] += value
	c.cpt++
	if c.batch_size > 0 && c.cpt == c.batch_size {
		if err := c.harvest(); err != nil {
			return true, err
		}
		return true, nil
	}
	return false, nil
}

func (c *Counters[K]) harvest() error {
	if len(c.counters) == 0 {
		return nil
	}
	keys := make([]K, len(c.counters))
	values := make([]int, len(c.counters))
	i := 0
	for k, v := range c.counters {
		keys[i] = k
		values[i] = v
		i++
	}
	if err := c.harvester(keys, values); err != nil {
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
