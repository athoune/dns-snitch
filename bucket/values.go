package bucket

import (
	"fmt"
	"strings"
)

type BucketValues struct {
	values      []int
	current_pos int
	length      int
}

func NewBucketValues(capacity, value int) *BucketValues {
	b := &BucketValues{
		values:      make([]int, capacity),
		current_pos: 0,
		length:      1,
	}
	b.values[0] = value
	return b
}

func (b *BucketValues) Sum() int {
	v := 0
	for _, i := range b.values {
		v += i
	}
	return v
}

// Leak the last value
func (b *BucketValues) Leak() {
	if b.length == 0 { // Can't remove thing from an empty collection
		return
	}
	capacity := len(b.values)
	n := b.current_pos + 1
	if n > capacity {
		n -= capacity
	}
	b.values[n] = 0
	b.length--
}

func (b *BucketValues) Add(value int) {
	b.values[b.next()] = value
	if b.length < len(b.values) {
		b.length++
	}
}

func (b *BucketValues) next() int {
	b.current_pos++
	if b.current_pos == len(b.values) {
		b.current_pos = 0
	}
	return b.current_pos
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
