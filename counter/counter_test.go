package counter

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBatchCounter(t *testing.T) {
	n := 100
	done := 0
	total := 0
	c := New[string](10, 10*time.Second, func(k []string, v []int) error {
		fmt.Println("v", v)
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
	cpt := 0
	w.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			if ok, _ := c.Add("pim", 1); ok {
				cpt++
			}
			w.Done()
		}()
	}
	w.Wait()
	c.Harvest()
	if done != n/10 {
		t.Error("Not enough harvester loop", done, "!=", n/10)
	}
	if cpt != n/10 {
		t.Error("Bad cpt", cpt)
	}
	if total != n {
		t.Error("Wrong total", total, "!=", n)
	}
}

func TestTimeCounter(t *testing.T) {
	n := 100
	cpt := 0
	c := New[string](10, 100*time.Millisecond, func(k []string, v []int) error {
		fmt.Println(v)
		if len(v) == 0 {
			return nil
		}
		cpt += v[0]
		return nil
	})
	for i := 0; i < n; i++ {
		go func() {
			_, err := c.Add("pim", 1)
			if err != nil {
				t.Error("Add loop error :", err)
			}
		}()
	}
	time.Sleep(200 * time.Millisecond)
	if cpt != n {
		t.Error("Not enough", cpt)
	}
}
