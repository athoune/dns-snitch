package bucket

import (
	"fmt"
	"io"
	"log/slog"
	"sync"
)

type LeakyBucket[K comparable] struct {
	datas    map[K]*BucketValues
	capacity int
	lock     *sync.RWMutex
}

func NewLeakyBucket[K comparable](capacity int) *LeakyBucket[K] {
	return &LeakyBucket[K]{
		datas:    make(map[K]*BucketValues),
		capacity: capacity,
		lock:     &sync.RWMutex{},
	}
}

func (l *LeakyBucket[K]) Get(line K) *BucketValues {
	l.lock.RLock()
	defer l.lock.RUnlock()
	v, ok := l.datas[line]
	if ok {
		return v
	}
	return nil
}

func (l *LeakyBucket[K]) Dump(out io.Writer) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	for k, v := range l.datas {
		fmt.Fprint(out, k)
		fmt.Fprintln(out, " => ", v.Sum(), "\n", v)
	}
}

func (l *LeakyBucket[K]) LeaksAll() {
	l.lock.Lock()
	defer l.lock.Unlock()
	before := len(l.datas)
	olds := make([]K, 0)
	for k, v := range l.datas {
		v.Leak()
		if v.Sum() == 0 {
			olds = append(olds, k)
		}
	}
	for _, old := range olds {
		delete(l.datas, old)
	}
	slog.Info("LeaksAll", "before", before, "after", len(l.datas), "olds", olds)
}

func (l *LeakyBucket[K]) Add(line K, value int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	current, ok := l.datas[line]
	if ok {
		current.Add(value)
	} else {
		current = NewBucketValues(l.capacity, value)
		l.datas[line] = current
	}
}

func (l *LeakyBucket[K]) Values() ([]K, []int) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	kk := make([]K, len(l.datas))
	vv := make([]int, len(l.datas))
	i := 0
	for k, v := range l.datas {
		kk[i] = k
		vv[i] = v.Sum()
		i++
	}
	return kk, vv
}

func (l *LeakyBucket[K]) Length() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return len(l.datas)
}
