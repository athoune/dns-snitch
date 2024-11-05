package output

import (
	"os"
	"testing"

	"time"

	"github.com/athoune/dns-snitch/counter"
)

func TestParquet(t *testing.T) {
	f, err := os.CreateTemp("", "parquet")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())
	w := NewFromWriter(f)
	c := counter.New[*Line](10, 10*time.Second, w.Write)
	c.Add(&Line{
		From:      "192.168.1.1",
		Target:    "127.0,.0,.1",
		TS:        int64(time.Now().Nanosecond()),
		Domain:    "localhost",
		Port:      80,
		Direction: "up",
	}, 42)
	err = c.Harvest()
	if err != nil {
		t.Error("Can't write parquet", err)
	}
	s, _ := f.Stat()
	if s.Size() == 0 {
		t.Error("Empty parquet file", s.Size())
	}
}
