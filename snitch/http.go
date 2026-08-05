package snitch

import (
	"net"

	"github.com/athoune/dns-snitch/output"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func (s *Snitch) buildLine(src, dest net.IP, tcp *layers.TCP) *output.Line {
	l := &output.Line{
		From:   src.String(),
		Target: dest.String(),
	}
	if tcp.DstPort == 443 || tcp.DstPort == 80 {
		l.Port = uint32(tcp.DstPort)
		l.Direction = "UP"
		l.Domain = s.domain(dest)
	}
	if tcp.SrcPort == 443 || tcp.SrcPort == 80 {
		l.Port = uint32(tcp.SrcPort)
		l.Direction = "DOWN"
		l.Domain = s.domain(src)
	}
	return l
}

func (s *Snitch) readHTTPPacket(src, dest net.IP, tcp *layers.TCP, packet gopacket.Packet) error {
	l := s.buildLine(src, dest, tcp)
	if l.Direction == "" {
		return nil
	}
	ap := packet.ApplicationLayer()
	if ap == nil {
		return nil
	}
	size := len(ap.Payload())
	for _, counter := range s.counters {
		if _, err := counter.Add(*l, size); err != nil {
			return err
		}
	}
	return nil
}
