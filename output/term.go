package output

import (
	"fmt"
	"os"
	"sort"

	"github.com/dustin/go-humanize"
	"golang.org/x/term"
)

type Term struct {
	term    int
	buckets LeakyBucket
}

func NewTerm(capacity int) *Term {
	return &Term{
		term:    int(os.Stdin.Fd()),
		buckets: *NewLeakyBucket(capacity),
	}

}

func (t *Term) Write(k []*Line, v []int) error {
	t.buckets.LeaksAll()
	for i, line := range k {
		t.buckets.Add(line, v[i])
	}
	width, height, err := term.GetSize(t.term)
	if err != nil {
		return err
	}
	pattern := fmt.Sprintf("|%%-%ds %%4s|%%6s|\n", width-14)
	fmt.Print("\033[H\033[2J") // clear screen
	lines := min(height, len(k)) - 2
	lv := Lines2LineValues(k, v)
	sort.Sort(LineValueBySize(lv))

	for i, line := range lv {
		fmt.Printf(pattern, line.Domain, line.Direction, humanize.Bytes(uint64(line.Size)))
		if i == lines {
			break
		}
	}
	return nil
}
