package snitch

import (
	"net"
	"net/netip"
	"sort"
)

// domain return the domain of an IP
func (s *Snitch) domain(ip net.IP) string {
	s.mutex.RLock()
	var names []string
	found := false
	if ip4 := ip.To4(); ip4 != nil {
		if domains, ok := s.resolution.Get(netip.AddrFrom4([4]byte(ip4))); ok {
			// Copy the names while holding the lock: AddResolution may
			// mutate the set concurrently.
			names, found = domains.Slice(), true
		}
	} else if domains, ok := s.resolution.Get(netip.AddrFrom16([16]byte(ip))); ok {
		names, found = domains.Slice(), true
	}
	s.mutex.RUnlock()
	if found {
		// An IP can be resolved by several names, pick a deterministic one.
		sort.Strings(names)
		return names[0]
	}
	addr, err := net.LookupAddr(ip.String())
	if err != nil {
		return ip.String()
	}
	return addr[0]
}
