package snitch

import (
	"io"
	"net/netip"
	"sync"
	"time"

	"github.com/athoune/dns-snitch/counter"
	"github.com/athoune/dns-snitch/output"
	"github.com/hashicorp/go-set"
	lru "github.com/hashicorp/golang-lru/v2"
)

type Snitch struct {
	resolution *lru.Cache[netip.Addr, *set.Set[string]]
	counter    *counter.Counters[*output.Line]
	writer     *output.Writer
	mutex      *sync.RWMutex
}

func New(batch_size int, batch_duration time.Duration, path string) (*Snitch, error) {
	l, _ := lru.New[netip.Addr, *set.Set[string]](256)
	v, err := output.New(path)
	if err != nil {
		return nil, err
	}

	return &Snitch{
		resolution: l,
		counter:    counter.New[*output.Line](batch_size, batch_duration, v.Write),
		writer:     v,
		mutex:      &sync.RWMutex{},
	}, nil
}

func (r *Snitch) Dump(output io.Writer) {
}
