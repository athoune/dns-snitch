package output

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/athoune/dns-snitch/bucket"
	"github.com/dustin/go-humanize"
	"golang.org/x/term"
)

type Term struct {
	term     int
	buckets  bucket.LeakyBucket[Line]
	truncate time.Duration
}

func NewTerm(capacity int, truncate time.Duration) *Term {
	return &Term{
		term:     int(os.Stdin.Fd()),
		buckets:  *bucket.NewLeakyBucket[Line](capacity),
		truncate: truncate,
	}

}

func (t *Term) Write(k []Line, v []int) error {
	fmt.Print("\033[H\033[2J") // clear screen
	if len(k) == 0 {
		return nil
	}
	truncate := t.truncate.Microseconds()
	normalizedTS := (k[0].TS / truncate) * truncate
	for i, line := range k {
		line.TS = normalizedTS
		t.buckets.Add(line, v[i])
	}
	width, height, err := term.GetSize(t.term)
	if err != nil {
		return err
	}
	pattern := fmt.Sprintf("|%%-%ds :%%-4d %%-4s|%%7s|\n", width-21)
	lines := min(height, len(k)) - 2
	ll, vv := t.buckets.Values()
	lv := Lines2LineValues(ll, vv)
	sort.Sort(LineValueBySize(lv))

	for i, line := range lv {
		fmt.Printf(pattern, line.Domain, line.Port, line.Direction, humanize.Bytes(uint64(line.Size)))
		if i == lines {
			break
		}
	}
	return nil
}
