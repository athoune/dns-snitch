package bucket

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

func TestBucket(t *testing.T) {
	bucket := NewLeakyBucket[string](6)
	line := "popo.com"
	bucket.Add(line, 1)
	bucket.Add(line, 2)
	if bucket.Length() != 1 {
		t.Fatalf("bad bucket length: got %d, want 1", bucket.Length())
	}

	// Dump renders one "key => sum" line per key.
	buff := &bytes.Buffer{}
	if err := bucket.Dump(buff); err != nil {
		t.Fatalf("dump error: %v", err)
	}
	if dump := buff.String(); dump != "popo.com => 3\n" {
		t.Fatalf("bad dump: %q", dump)
	}

	v := bucket.Get(line)
	if v == nil {
		t.Fatalf("unknown line %q", line)
	}
	if v.Length() != 2 {
		t.Errorf("bad length: got %d, want 2", v.Length())
	}
	if s := v.Sum(); s != 3 {
		t.Errorf("bad sum: got %d, want 3", s)
	}

	// Fill the bucket up to its capacity.
	bucket.Add(line, 3)
	bucket.Add(line, 4)
	bucket.Add(line, 5)
	bucket.Add(line, 6)
	if s := bucket.Get(line).Sum(); s != 21 {
		t.Errorf("bad sum: got %d, want 21", s)
	}

	// One more add wraps around and overwrites the oldest slot.
	bucket.Add(line, 7)
	v = bucket.Get(line)
	if s := v.Sum(); s != 27 {
		t.Errorf("bad sum after wrap: got %d, want 27", s)
	}
	if v.current_pos != 0 {
		t.Errorf("pointer lost: got %d, want 0", v.current_pos)
	}
	if v.Length() != 6 {
		t.Errorf("bad length after wrap: got %d, want 6", v.Length())
	}

	// Leaking zeroes the oldest slot and shrinks the window.
	v.Leak()
	if v.current_pos != 1 {
		t.Errorf("pointer lost: got %d, want 1", v.current_pos)
	}
	if v.Length() != 5 {
		t.Errorf("bad length after leak: got %d, want 5", v.Length())
	}
}

func TestBucketValues(t *testing.T) {
	bucket := NewLeakyBucket[string](3)
	type Data struct {
		domain string
		size   int
	}
	datas := []Data{
		{"popo.com", 42},
		{"popo.com", 2},
		{"popo.com", 3},
	}
	for _, data := range datas {
		bucket.Add(data.domain, data.size)
	}
	keys, values := bucket.Values()
	if len(keys) != 1 {
		t.Fatalf("got %d keys, want 1", len(keys))
	}
	if keys[0] != "popo.com" {
		t.Errorf("got key %q, want %q", keys[0], "popo.com")
	}
	if values[0] != 47 {
		t.Errorf("got value %d, want 47", values[0])
	}
}

func TestLeak(t *testing.T) {
	bucket := NewLeakyBucket[string](3)
	type Data struct {
		domain string
		size   int
	}
	for _, data := range []Data{
		{"popo.com", 42},
		{"popo.com", 2},
		{"popo.com", 3},
	} {
		bucket.Add(data.domain, data.size)
	}
	line := bucket.Get("popo.com")
	if line.Length() != 3 {
		t.Fatalf("bad length: got %d, want 3", line.Length())
	}
	if line.Sum() != 47 {
		t.Fatalf("bad sum: got %d, want 47", line.Sum())
	}

	bucket.LeaksAll()
	line = bucket.Get("popo.com")
	if line.Length() != 2 {
		t.Errorf("bad length after leak: got %d, want 2", line.Length())
	}
	if line.Sum() != 5 {
		t.Errorf("bad sum after leak: got %d, want 5", line.Sum())
	}

	bucket.LeaksAll()
	line = bucket.Get("popo.com")
	if line.Length() != 1 {
		t.Errorf("bad length after second leak: got %d, want 1", line.Length())
	}
	if line.Sum() != 3 {
		t.Errorf("bad sum after second leak: got %d, want 3", line.Sum())
	}
	if line.Recyclable() {
		t.Error("not recyclable yet")
	}

	olds := bucket.LeaksAll()
	if olds != 1 {
		t.Errorf("bad LeaksAll result: got %d, want 1", olds)
	}
	if line = bucket.Get("popo.com"); line != nil {
		t.Errorf("line is not nil after full leak: %v", line)
	}
	if bucket.Length() > 0 {
		t.Errorf("bucket is not empty: %v", bucket.data)
	}
	if olds = bucket.LeaksAll(); olds != 0 {
		t.Errorf("bad LeaksAll result on empty bucket: got %d, want 0", olds)
	}
}

// TestConcurrentBucket exercises the whole API from many goroutines.
// Run with -race to detect data races.
func TestConcurrentBucket(t *testing.T) {
	b := NewLeakyBucket[string](4)
	const goroutines = 32
	const adds = 200
	w := &sync.WaitGroup{}
	w.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(n int) {
			defer w.Done()
			key := fmt.Sprintf("domain-%d.com", n%8)
			for j := 0; j < adds; j++ {
				b.Add(key, j)
				if j%10 == 0 {
					b.LeaksAll()
					b.Values()
					b.Length()
					b.Get(key)
				}
			}
		}(i)
	}
	w.Wait()
	if b.Length() > 8 {
		t.Errorf("unexpected number of keys: %d", b.Length())
	}
}
