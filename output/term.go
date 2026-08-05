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
	slog.Info("Term.Write", "lines", len(k))
	// buckets management
	t.buckets.LeaksAll()
	for i, line := range k {
		t.buckets.Add(line, v[i])
	}

	// terminal sizes
	width, height, err := term.GetSize(t.term)
	if err != nil {
		return err
	}
	pattern := fmt.Sprintf("|%%-%ds :%%-4d %%-4s|%%7s|%%2d\n", width-23)

	// Values
	ll, vv := t.buckets.Values()
	lines := min(height-2, len(ll))
	lv := Lines2LineValues(ll, vv)
	sort.Sort(LineValueBySize(lv))

	// Print to terminal
	fmt.Print("\033[H\033[2J") // clear screen
	fmt.Println(len(k), "fresh lines", t.buckets.Length(), "current elements")
	for i, line := range lv {
		fmt.Printf(pattern, line.Domain+" ["+line.Server()+"]", line.Port, line.Direction,
			humanize.Bytes(uint64(line.Weight)), t.buckets.Get(line.Line).Length())
		if i == lines {
			h := len(ll) - i
			if h > 0 {
				fmt.Print(h, "…")
			}
			break
		}
	}

	return nil
}
