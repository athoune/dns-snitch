package bucket

import (
	"fmt"
	"testing"
)

func TestBucket(t *testing.T) {
	bucket := NewLeakyBucket[string](6)
	line := "popo.com"
	bucket.Add(line, 1)
	bucket.Add(line, 2)
	v := bucket.Get(line)
	fmt.Println(v)
	if v.length != 2 {
		t.Error("Bad length 2 !=", v.length)
	}
	fmt.Println(v.String())
	if v == nil {
		t.Error("Unknown line", line, bucket.datas)
	}
	s := v.Sum()
	if s != 3 {
		t.Error("Sum error", s, "!= 3")
	}
	bucket.Add(line, 3)
	bucket.Add(line, 4)
	bucket.Add(line, 5)
	bucket.Add(line, 6)
	fmt.Println(v.String())
	v = bucket.Get(line)
	s = v.Sum()
	if s != 21 {
		t.Error("Sum error", s, "!= 21", v.values, v.current_pos)
	}
	bucket.Add(line, 7) // 7, it loops
	fmt.Println(v.String())
	v = bucket.Get(line)
	s = v.Sum()
	if s != 27 {
		t.Error("Sum error", s, "!= 27", v.values, v.current_pos)
	}
	if v.current_pos != 0 {
		t.Error("The pointer is lost 0 !=", v.current_pos)
	}
	v.Leak()
	if v.current_pos != 0 {
		t.Error("The pointer is lost 0 !=", v.current_pos)
	}
	if v.length != 5 {
		fmt.Println(v)
		t.Error("Bad length 5 !=", v.length)
	}
	fmt.Println(v.String())
}
