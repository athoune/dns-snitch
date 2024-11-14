package snitch

import (
	"net"
	"net/netip"

	"github.com/hashicorp/go-set"
)

// domain return the domain of an IP
func (s *Snitch) domain(ip net.IP) string {
	s.mutex.RLock()
	var domains *set.Set[string]
	var ok bool
	if len(ip) == 4 {
		domains, ok = s.resolution.Get(netip.AddrFrom4([4]byte(ip)))
	} else {
		domains, ok = s.resolution.Get(netip.AddrFrom16([16]byte(ip)))
	}
	s.mutex.RUnlock()
	if ok {
		return domains.Slice()[0] // FIXME
	} else {
		addr, err := net.LookupAddr(ip.String())
		if err != nil {
			return ip.String()
		}
		return addr[0]
	}
}
