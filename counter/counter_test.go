package counter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBatchCounter(t *testing.T) {
	n := 100
	done := 0
	total := 0
	c := New[string](10, 10*time.Second, func(k []string, v []int) error {
		if len(k) == 0 {
			return nil
		}
		if v[0] != n/10 {
			t.Error("cb, v!=n", v, n/10)
		}
		if k[0] != "pim" {
			t.Error("Bad value", v[0])
		}
		done += len(k)
		for _, i := range v {
			total += i
		}
		return nil
	})
	w := &sync.WaitGroup{}
	var cpt atomic.Int32
	w.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			if ok, _ := c.Add("pim", 1); ok {
				cpt.Add(1)
			}
			w.Done()
		}()
	}
	w.Wait()
	c.Harvest()
	if done != n/10 {
		t.Error("Not enough harvester loop", done, "!=", n/10)
	}
	if cpt.Load() != int32(n/10) {
		t.Error("Bad cpt", cpt.Load())
	}
	if total != n {
		t.Error("Wrong total", total, "!=", n)
	}
}

// TestTimeCounter checks that the timer triggers harvests even when the batch
// size is never reached. The batch is larger than the number of adds, so only
// the timer can flush the counters.
func TestTimeCounter(t *testing.T) {
	n := 100
	var cpt atomic.Int32
	harvested := make(chan struct{}, 1)
	c := New[string](n+1, 50*time.Millisecond, func(k []string, v []int) error {
		for _, i := range v {
			cpt.Add(int32(i))
		}
		select {
		case harvested <- struct{}{}:
		default:
		}
		return nil
	})
	w := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			if _, err := c.Add("pim", 1); err != nil {
				t.Error("Add loop error :", err)
			}
		}()
	}
	w.Wait()

	// The timer must fire at least once.
	select {
	case <-harvested:
	case <-time.After(2 * time.Second):
		t.Fatal("timer never harvested")
	}

	// Drain whatever is left, then the total must be exact.
	if err := c.Harvest(); err != nil {
		t.Fatal(err)
	}
	if cpt.Load() != int32(n) {
		t.Error("Not enough", cpt.Load())
	}
}
