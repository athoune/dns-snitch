package bucket

import "testing"

func TestValues(t *testing.T) {
	b := NewBucketValues(3, 42)
	if b.Length() != 1 {
		t.Fatalf("bad length: got %d, want 1", b.Length())
	}
	b.Add(12)
	b.Add(3)
	if s := b.Sum(); s != 57 {
		t.Fatalf("bad sum: got %d, want 57", s)
	}
	b.Add(1) // wraps around, overwriting the first slot
	if s := b.Sum(); s != 16 {
		t.Fatalf("bad sum after wrap: got %d, want 16", s)
	}
	b.Leak()
	if s := b.Sum(); s != 4 {
		t.Fatalf("bad sum after leak: got %d, want 4", s)
	}
}

func TestLeakOnEmpty(t *testing.T) {
	b := NewBucketValues(3, 42)
	if got := b.Leak(); got != 0 {
		t.Fatalf("bad leak result: got %d, want 0", got)
	}
	if !b.Recyclable() {
		t.Fatal("expected recyclable after draining all values")
	}
	// Leaking an empty collection is a no-op.
	if got := b.Leak(); got != -1 {
		t.Fatalf("bad leak on empty: got %d, want -1", got)
	}
}
