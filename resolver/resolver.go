package resolver

import (
	"fmt"
	"io"
	"net/netip"
	"sync"

	"github.com/dustin/go-humanize"
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
	/*
			r.mutex.RLock()
			ips := r.resolution.Keys()
			values := r.resolution.Values()
			r.mutex.RUnlock()
		fmt.Fprint(out, "\n===\n")
		for i, ip := range ips {
			fmt.Fprintf(out, "%s %s\n", ip, values[i])
		}
	*/
	fmt.Fprint(out, "\n---\n")
	r.bpMutex.RLock()
	down := sortKV(r.download)
	for _, kv := range down {
		fmt.Fprintf(out, "%-40s |▼ %-20s", kv.k, humanize.Bytes(uint64(kv.v)))
		up, ok := r.upload[kv.k]
		if ok {
			fmt.Fprintf(out, " |▲ %-20s", humanize.Bytes(uint64(up)))
		}
		fmt.Fprint(out, "\n")
	}
	r.bpMutex.RUnlock()
	fmt.Fprint(out, "\n")
}
