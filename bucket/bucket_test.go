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
	if v.current_pos != 1 {
		t.Error("The pointer is lost 5 !=", v.current_pos)
	}
	if v.length != 5 {
		fmt.Println(v)
		t.Error("Bad length 5 !=", v.length)
	}
	fmt.Println(v.String())
}

func TestBucketValues(t *testing.T) {
	bucket := NewLeakyBucket[string](3)
	type Data struct {
		domain string
		size   int
	}
	datas := []Data{
		{
			"popo.com",
			42,
		},
		{
			"popo.com",
			2,
		},
		{
			"popo.com",
			3,
		},
	}
	for _, data := range datas {
		bucket.Add(data.domain, data.size)
	}
	keys, values := bucket.Values()
	if len(keys) != 1 {
		t.Error("?! keys", len(keys))
	}
	if keys[0] != "popo.com" {
		t.Error("?! key name", keys[0])
	}
	if values[0] != 47 {
		t.Error("?! values", values)
	}

}

func TestLeak(t *testing.T) {
	bucket := NewLeakyBucket[string](3)
	type Data struct {
		domain string
		size   int
	}
	datas := []Data{
		{
			"popo.com",
			42,
		},
		{
			"popo.com",
			2,
		},
		{
			"popo.com",
			3,
		},
	}
	for _, data := range datas {
		bucket.Add(data.domain, data.size)
	}
	line := bucket.Get("popo.com")
	if line.length != 3 {
		t.Error("oups length", line.length)
	}
	if line.Sum() != 47 {
		t.Error("Oups sum", line.Sum())
	}
	bucket.LeaksAll()
	line = bucket.Get("popo.com")
	if line.length != 2 {
		t.Error("oups length", line.length)
	}
	if line.Sum() != 5 {
		t.Error("Oups sum", line.Sum())
	}
	bucket.LeaksAll()
	line = bucket.Get("popo.com")
	if line.length != 1 {
		t.Error("oups length", line.length)
	}
	if line.Sum() != 3 {
		t.Error("Oups sum", line.Sum())
	}
	bucket.LeaksAll()
	line = bucket.Get("popo.com")
	if line != nil {
		t.Error("Line is not nil")
	}
	if len(bucket.datas) > 0 {
		t.Error("datas is not empty", bucket.datas)
	}
}
