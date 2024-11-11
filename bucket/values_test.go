package bucket

import (
	"fmt"
	"testing"
)

func TestValues(t *testing.T) {
	b := NewBucketValues(3, 42)
	if b.length != 1 {
		t.Error(b.length, "!=", 1)
	}
	b.Add(12)
	b.Add(3)
	sum := b.Sum()
	if sum != 57 {
		t.Error(sum, "!=", 57)
	}
	b.Add(1)
	sum = b.Sum()
	if sum != 16 {
		t.Error(sum, "!=", 16)
	}
	fmt.Println(b)
	b.Leak()
	fmt.Println(b)
	sum = b.Sum()
	if sum != 4 {
		fmt.Println("current", b.current_pos)
		t.Error(sum, "!=", 4)
	}
}
