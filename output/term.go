package output

import (
	"fmt"
	"log/slog"
	"os"
	"sort"
	"time"

	"github.com/athoune/dns-snitch/bucket"
	"github.com/dustin/go-humanize"
	"golang.org/x/term"
)

type Term struct {
	term     int // term file descriptor
	buckets  bucket.LeakyBucket[Line]
	truncate time.Duration
}

func NewTerm(capacity int, truncate time.Duration) *Term {
	return &Term{
		term:     int(os.Stdout.Fd()),
		buckets:  *bucket.NewLeakyBucket[Line](capacity),
		truncate: truncate,
	}

}

func (t *Term) Write(k []Line, v []int) error {
	t.buckets.LeaksAll()
	fmt.Print("\033[H\033[2J") // clear screen
	fmt.Println(len(k), "fresh lines", t.buckets.Length(), "current elements")
	slog.Info("Term.Write", "lines", len(k))
	for i, line := range k {
		t.buckets.Add(line, v[i])
	}
	width, height, err := term.GetSize(t.term)
	if err != nil {
		return err
	}
	pattern := fmt.Sprintf("|%%-%ds :%%-4d %%-4s|%%7s|\n", width-21)
	lines := min(height, t.buckets.Length()) - 2
	ll, vv := t.buckets.Values()
	lv := Lines2LineValues(ll, vv)
	sort.Sort(LineValueBySize(lv))

	for i, line := range lv {
		fmt.Printf(pattern, line.Domain, line.Port, line.Direction, humanize.Bytes(uint64(line.Weight)))
		if i == lines-2 {
			break
		}
	}
	return nil
}
