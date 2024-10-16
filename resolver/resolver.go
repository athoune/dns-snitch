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
	mutex      *sync.RWMutex
}

func New() *Resolver {
	l, _ := lru.New[netip.Addr, *set.Set[string]](256)
	return &Resolver{
		resolution: l,
		mutex:      &sync.RWMutex{},
	}
}

func (r *Resolver) Dump(out io.Writer) {
	r.mutex.RLock()
	ips := r.resolution.Keys()
	values := r.resolution.Values()
	r.mutex.RUnlock()
	fmt.Fprint(out, "\n---\n")
	for i, ip := range ips {
		fmt.Fprintf(out, "%s %s\n", ip, values[i])
	}
}

func (r *Resolver) Append(addr netip.Addr, domain string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	value, ok := r.resolution.Get(addr)
	if !ok {
		value = set.New[string](3)
	}
	value.Insert(domain)
	r.resolution.Add(addr, value)
}
