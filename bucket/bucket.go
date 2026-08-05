package bucket

import (
	"fmt"
	"io"
	"log/slog"
	"sync"
)

type LeakyBucket[K comparable] struct {
	data     map[K]*BucketValues
	capacity int
	lock     *sync.RWMutex
}

func NewLeakyBucket[K comparable](capacity int) *LeakyBucket[K] {
	return &LeakyBucket[K]{
		data:     make(map[K]*BucketValues),
		capacity: capacity,
		lock:     &sync.RWMutex{},
	}
}

func (l *LeakyBucket[K]) Get(line K) *BucketValues {
	l.lock.RLock()
	defer l.lock.RUnlock()
	v, ok := l.data[line]
	if ok {
		return v
	}
	return nil
}

func (l *LeakyBucket[K]) Dump(out io.Writer) error {
	l.lock.RLock()
	defer l.lock.RUnlock()
	for k, v := range l.data {
		_, err := fmt.Fprintf(out, "%v => %d\n", k, v.Sum())
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *LeakyBucket[K]) LeaksAll() int {
	l.lock.Lock()
	defer l.lock.Unlock()
	before := len(l.data)
	olds := make([]K, 0)
	for k, v := range l.data {
		if v.Leak() == 0 {
			olds = append(olds, k)
		}
	}
	for _, old := range olds {
		delete(l.data, old)
	}
	slog.Info("LeaksAll", "before", before, "after", len(l.data), "olds", olds)
	return len(olds)
}

func (l *LeakyBucket[K]) Add(line K, value int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	current, ok := l.data[line]
	if ok {
		current.Add(value)
	} else {
		current = NewBucketValues(l.capacity, value)
		l.data[line] = current
	}
}

func (l *LeakyBucket[K]) Values() ([]K, []int) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	kk := make([]K, len(l.data))
	vv := make([]int, len(l.data))
	i := 0
	for k, v := range l.data {
		kk[i] = k
		vv[i] = v.Sum()
		i++
	}
	return kk, vv
}

func (l *LeakyBucket[K]) Length() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return len(l.data)
}
