package output

import (
	"fmt"
	"strings"
)

type BucketValues struct {
	values      []int
	current_pos int
	length      int
}

type LeakyBucket struct {
	datas    map[Line]*BucketValues
	capacity int
}

func NewLeakyBucket(capacity int) *LeakyBucket {
	return &LeakyBucket{
		datas:    make(map[Line]*BucketValues),
		capacity: capacity,
	}
}

func (l *LeakyBucket) Get(line *Line) *BucketValues {
	v, ok := l.datas[*line]
	if ok {
		return v
	}
	return nil
}

func (l *LeakyBucket) LineValues() []*LineValue {
	values := make([]*LineValue, len(l.datas))
	i := 0
	for k, v := range l.datas {
		values[i] = &LineValue{
			Line: k,
			Size: int32(v.Sum()),
		}
		i++
	}
	return values
}

func (b *BucketValues) Sum() int {
	v := 0
	for _, i := range b.values {
		v += i
	}
	return v
}

func (b *BucketValues) Leak() {
	if b.length == 0 { // Can't remove thing from an empty collection
		return
	}
	capacity := len(b.values)
	n := b.current_pos - b.length
	if n < 0 {
		n += capacity
	}
	b.values[n] = 0
	b.length--
}

func (b *BucketValues) next() int {
	b.current_pos++
	if b.current_pos == len(b.values) {
		b.current_pos = 0
	}
	return b.current_pos
}

func (l *LeakyBucket) Add(line *Line, value int) {
	current, ok := l.datas[*line]
	if ok {
		current.values[current.next()] = value
	} else {
		current = &BucketValues{
			values:      make([]int, l.capacity),
			current_pos: 0,
		}
		current.values[0] = value
		l.datas[*line] = current
	}
	current.length = min(l.capacity, current.length+1)
}

func (b *BucketValues) String() string {
	sizes := make([]int, len(b.values))
	var buff strings.Builder
	fmt.Fprint(&buff, "|")
	for i, v := range b.values {
		sizes[i] = len(fmt.Sprintf("%d", v))
		fmt.Fprintf(&buff, " %d |", v)
	}
	fmt.Fprintln(&buff)
	for i, j := range sizes {
		if i == b.current_pos {
			fmt.Fprint(&buff, " ⬆")
			break
		}
		for range j + 3 {
			fmt.Fprint(&buff, " ")
		}
	}
	return buff.String()
}
