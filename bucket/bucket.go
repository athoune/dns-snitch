package bucket

import (
	"fmt"
	"io"
)

type LeakyBucket[K comparable] struct {
	datas    map[K]*BucketValues
	capacity int
}

func NewLeakyBucket[K comparable](capacity int) *LeakyBucket[K] {
	return &LeakyBucket[K]{
		datas:    make(map[K]*BucketValues),
		capacity: capacity,
	}
}

func (l *LeakyBucket[K]) Get(line K) *BucketValues {
	v, ok := l.datas[line]
	if ok {
		return v
	}
	return nil
}

func (l *LeakyBucket[K]) Dump(out io.Writer) {
	for k, v := range l.datas {
		fmt.Fprint(out, k)
		fmt.Fprintln(out, " => ", v.Sum(), "\n", v)
	}
}

func (l *LeakyBucket[K]) LeaksAll() {
	for _, v := range l.datas {
		v.Leak()
	}
}

func (l *LeakyBucket[K]) Add(line K, value int) {
	current, ok := l.datas[line]
	if ok {
		current.Add(value)
	} else {
		current = NewBucketValues(l.capacity, value)
		l.datas[line] = current
	}
}

func (l *LeakyBucket[K]) Values() ([]K, []int) {
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
