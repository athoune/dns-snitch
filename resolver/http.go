package resolver

import (
	"net"
	"net/netip"
	"time"

	"github.com/google/gopacket/layers"
	"github.com/hashicorp/go-set"
)

// domain return the domain of an IP
func (r *Resolver) domain(ip net.IP) string {
	r.mutex.RLock()
	var domains *set.Set[string]
	var ok bool
	if len(ip) == 4 {
		domains, ok = r.resolution.Get(netip.AddrFrom4([4]byte(ip)))
	} else {
		domains, ok = r.resolution.Get(netip.AddrFrom16([16]byte(ip)))
	}
	r.mutex.RUnlock()
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

func (r *Resolver) readHTTP(src, dest net.IP, tcp *layers.TCP) error {
	size := len(tcp.Payload)
	// fmt.Printf("src, dest : %s:%s %s:%s\n", src, tcp.SrcPort, dest, tcp.DstPort)
	if tcp.DstPort == 443 || tcp.DstPort == 80 {
		domain := r.domain(dest)
		r.bpMutex.Lock()
		s, ok := r.upload[domain]
		if !ok {
			s = 0
		}
		r.upload[domain] = s + size
		r.lines <- &Line{
			Port:   uint32(tcp.DstPort),
			Domain: domain,
			Target: dest,
			From:   src,
			Size:   uint32(size),
			TS:     time.Now().Local().UnixMicro(),
		}
		r.bpMutex.Unlock()
		return nil
	}
	if tcp.SrcPort == 443 || tcp.SrcPort == 80 {
		domain := r.domain(src)
		r.bpMutex.Lock()
		s, ok := r.download[domain]
		if !ok {
			s = 0
		}
		r.download[domain] = s + size
		r.bpMutex.Unlock()
		return nil
	}
	return nil
}
