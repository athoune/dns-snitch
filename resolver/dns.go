package resolver

import (
	"net/netip"

	"github.com/google/gopacket/layers"
	"github.com/miekg/dns"
)

func (r *Resolver) readDNS(udp *layers.UDP) error {
	msg := &dns.Msg{}
	err := msg.Unpack(udp.Payload)
	if err != nil {
		return err
	}
	if msg.Response {
		for _, answer := range msg.Answer {
			switch answer.(type) {
			case *dns.A:
				a, _ := answer.(*dns.A)
				addr := netip.AddrFrom4([4]byte(a.A))
				r.Append(addr, a.Hdr.Name)
			case *dns.AAAA:
				a, _ := answer.(*dns.AAAA)
				addr := netip.AddrFrom16([16]byte(a.AAAA))
				r.Append(addr, a.Hdr.Name)
			default:
				//dump.P(msg)
			}
		}
	}
	return nil
}
