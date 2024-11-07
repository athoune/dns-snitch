package snitch

import (
	"io"
	"net/netip"
	"sync"

	"github.com/athoune/dns-snitch/counter"
	"github.com/athoune/dns-snitch/output"
	"github.com/hashicorp/go-set"
	lru "github.com/hashicorp/golang-lru/v2"
)

type Snitch struct {
	resolution *lru.Cache[netip.Addr, *set.Set[string]]
	counters   []*counter.Counters[*output.Line]
	mutex      *sync.RWMutex
}

func New() *Snitch {
	l, _ := lru.New[netip.Addr, *set.Set[string]](256)

	return &Snitch{
		resolution: l,
		counters:   make([]*counter.Counters[*output.Line], 0),
		mutex:      &sync.RWMutex{},
	}
}

func (s *Snitch) AddCounter(c *counter.Counters[*output.Line]) {
	s.counters = append(s.counters, c)
}

func (r *Snitch) Dump(output io.Writer) {
}
