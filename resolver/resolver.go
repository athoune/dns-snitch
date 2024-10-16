package resolver

import (
	"fmt"
	"io"
	"net/netip"
	"sync"

	"github.com/hashicorp/go-set"
	lru "github.com/hashicorp/golang-lru/v2"
)

type Resolver struct {
	resolution *lru.Cache[netip.Addr, *set.Set[string]]
	upload     map[string]int
	download   map[string]int
	mutex      *sync.RWMutex
	bpMutex    *sync.RWMutex
}

func New() *Resolver {
	l, _ := lru.New[netip.Addr, *set.Set[string]](256)
	return &Resolver{
		resolution: l,
		mutex:      &sync.RWMutex{},
		bpMutex:    &sync.RWMutex{},
		upload:     make(map[string]int),
		download:   make(map[string]int),
	}
}

func (r *Resolver) Dump(out io.Writer) {
	r.mutex.RLock()
	ips := r.resolution.Keys()
	values := r.resolution.Values()
	r.mutex.RUnlock()
	fmt.Fprint(out, "\n===\n")
	for i, ip := range ips {
		fmt.Fprintf(out, "%s %s\n", ip, values[i])
	}
	fmt.Fprint(out, "\n---\n")
	r.bpMutex.RLock()
	for k, v := range r.download {
		fmt.Fprintf(out, "v %s %d\n", k, v)
	}
	for k, v := range r.upload {
		fmt.Fprintf(out, "^ %s %d\n", k, v)
	}
	r.bpMutex.RUnlock()
	fmt.Fprint(out, "\n")
}
