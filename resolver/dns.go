package resolver

import (
	"net/netip"

	"github.com/google/gopacket/layers"
	"github.com/hashicorp/go-set"
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
				r.AddResolution(addr, a.Hdr.Name)
			case *dns.AAAA:
				a, _ := answer.(*dns.AAAA)
				addr := netip.AddrFrom16([16]byte(a.AAAA))
				r.AddResolution(addr, a.Hdr.Name)
			default:
				//dump.P(msg)
			}
		}
	}
	return nil
}

func (r *Resolver) AddResolution(addr netip.Addr, domain string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	value, ok := r.resolution.Get(addr)
	if !ok {
		value = set.New[string](3)
	}
	value.Insert(domain)
	r.resolution.Add(addr, value)
}
