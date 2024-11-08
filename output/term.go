package output

import (
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
	"golang.org/x/term"
)

type Term struct {
	term int
}

func NewTerm() *Term {
	return &Term{
		term: int(os.Stdin.Fd()),
	}

}

func (t *Term) Write(k []*Line, v []int) error {
	width, height, err := term.GetSize(t.term)
	if err != nil {
		return err
	}
	pattern := fmt.Sprintf("|%%-%ds|%%6s|\n", width-9)
	fmt.Print("\033[H\033[2J") // clear screen
	lines := min(height, len(k)) - 2
	i := 0
	for range len(k) {
		if v[i] == 0 {
			continue
		}
		i++
		fmt.Printf(pattern, k[i].Domain, humanize.Bytes(uint64(v[i])))
		if i == lines {
			break
		}
	}
	return nil
}
