package output

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/athoune/dns-snitch/counter"
	"github.com/parquet-go/parquet-go"
)

func TestParquet(t *testing.T) {
	f, err := os.CreateTemp("", "parquet")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	w := NewFromWriter(f)
	c := counter.New[Line](10, 10*time.Second, w.Write)
	c.Add(Line{
		From:      "192.168.1.1",
		Target:    "127.0.0.1",
		TS:        time.Now().UnixMicro(),
		Domain:    "localhost",
		Port:      80,
		Direction: "UP",
	}, 42)
	if err := c.Harvest(); err != nil {
		t.Fatal(err)
	}
	if info, _ := f.Stat(); info.Size() == 0 {
		t.Error("empty parquet file")
	}
}

// TestParquetHonorsPath is a regression test: New used to ignore its path
// argument and always wrote to ./snitch.parquet.
func TestParquetHonorsPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom.parquet")
	w, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	line := Line{
		From:      "192.168.1.1",
		Target:    "93.184.216.34",
		TS:        time.Now().UnixMicro(),
		Domain:    "example.com",
		Port:      443,
		Direction: "UP",
	}
	if err := w.Write([]Line{line}, []int{42}); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(path); err != nil {
		t.Fatalf("parquet file not written at the requested path: %v", err)
	} else if info.Size() == 0 {
		t.Fatal("empty parquet file")
	}

	// Round-trip: read the file back and check the row content.
	rows, err := parquet.ReadFile[*LineValue](path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(rows))
	}
	got := rows[0]
	if got.Domain != "example.com" || got.Direction != "UP" || got.Weight != 42 {
		t.Errorf("bad row: %+v", got)
	}
	if got.From != "192.168.1.1" || got.Target != "93.184.216.34" || got.Port != 443 {
		t.Errorf("bad row addresses: %+v", got)
	}
}
