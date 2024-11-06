package snitch

import (
	"net"
	"net/netip"
	"time"

	"github.com/athoune/dns-snitch/output"
	"github.com/google/gopacket/layers"
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

func (s *Snitch) readHTTP(src, dest net.IP, tcp *layers.TCP) error {
	size := len(tcp.Payload)
	var err error
	// fmt.Printf("src, dest : %s:%s %s:%s\n", src, tcp.SrcPort, dest, tcp.DstPort)
	l := &output.Line{
		From:   src.String(),
		Target: dest.String(),
		Domain: s.domain(dest),
		TS:     time.Now().UnixMicro(),
	}
	if tcp.DstPort == 443 || tcp.DstPort == 80 {
		l.Port = uint32(tcp.DstPort)
		l.Direction = "UP"
	}
	if tcp.SrcPort == 443 || tcp.SrcPort == 80 {
		l.Port = uint32(tcp.SrcPort)
		l.Direction = "DOWN"
	}
	if l.Direction != "" {
		_, err = s.counter.Add(l, size)
		if err != nil {
			return err
		}
	}
	return nil
}
