package output

import (
	"fmt"
	"testing"
)

func TestBucket(t *testing.T) {
	bucket := NewLeakyBucket(6)
	line := &Line{
		Domain: "popo.com",
	}
	bucket.Add(line, 1)
	bucket.Add(line, 2)
	v := bucket.Get(line)
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
	bucket.Add(line, 2)
	bucket.Add(line, 2)
	bucket.Add(line, 2)
	bucket.Add(line, 2)
	fmt.Println(v.String())
	v = bucket.Get(line)
	s = v.Sum()
	if s != 11 {
		t.Error("Sum error", s, "!= 11", v.values, v.current_pos)
	}
	bucket.Add(line, 2) // 7, it loops
	fmt.Println(v.String())
	v = bucket.Get(line)
	s = v.Sum()
	if s != 12 {
		t.Error("Sum error", s, "!= 12", v.values, v.current_pos)
	}
	if v.current_pos != 0 {
		t.Error("The pointer is lost 0 !=", v.current_pos)
	}
	v.Leak()
	if v.current_pos != 0 {
		t.Error("The pointer is lost 0 !=", v.current_pos)
	}
	if v.length != 5 {
		t.Error("Bad length 5 !=", v.length)
	}
	fmt.Println(v.String())
}
