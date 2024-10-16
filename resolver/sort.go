package resolver

import "sort"

type kv struct {
	k string
	v int
}

type byKV []kv

func (k byKV) Len() int           { return len(k) }
func (k byKV) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }
func (k byKV) Less(i, j int) bool { return k[i].v < k[j].v }

func sortKV(in map[string]int) []kv {
	out := make([]kv, len(in))
	i := 0
	for k, v := range in {
		out[i] = kv{k, v}
		i += 1
	}
	sort.Sort(byKV(out))
	return out
}
