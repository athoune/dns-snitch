package snitch

import (
	"net"
	"net/netip"
	"sync"
	"testing"

	"github.com/google/gopacket/layers"
	"github.com/miekg/dns"
)

// dnsUDP packs a DNS message into the payload of a fake UDP layer.
func dnsUDP(t *testing.T, msg *dns.Msg) *layers.UDP {
	t.Helper()
	data, err := msg.Pack()
	if err != nil {
		t.Fatalf("packing DNS message: %v", err)
	}
	return &layers.UDP{BaseLayer: layers.BaseLayer{Payload: data}}
}

func TestReadDNSARecord(t *testing.T) {
	s := New()
	msg := &dns.Msg{}
	msg.Response = true
	msg.Answer = append(msg.Answer, &dns.A{
		Hdr: dns.RR_Header{Name: "example.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
		A:   net.IPv4(93, 184, 216, 34).To4(),
	})
	if err := s.readDNS(dnsUDP(t, msg)); err != nil {
		t.Fatal(err)
	}
	if got := s.domain(net.IPv4(93, 184, 216, 34)); got != "example.com." {
		t.Errorf("got domain %q, want %q", got, "example.com.")
	}
}

func TestReadDNSAAAARecord(t *testing.T) {
	s := New()
	msg := &dns.Msg{}
	msg.Response = true
	msg.Answer = append(msg.Answer, &dns.AAAA{
		Hdr:  dns.RR_Header{Name: "ipv6.example.", Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 60},
		AAAA: net.ParseIP("2606:2800:220:1:248:1893:25c8:1946"),
	})
	if err := s.readDNS(dnsUDP(t, msg)); err != nil {
		t.Fatal(err)
	}
	if got := s.domain(net.ParseIP("2606:2800:220:1:248:1893:25c8:1946")); got != "ipv6.example." {
		t.Errorf("got domain %q, want %q", got, "ipv6.example.")
	}
}

func TestReadDNSIgnoresQueries(t *testing.T) {
	s := New()
	msg := &dns.Msg{} // not a response
	msg.Question = append(msg.Question, dns.Question{Name: "example.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET})
	if err := s.readDNS(dnsUDP(t, msg)); err != nil {
		t.Fatal(err)
	}
	if s.resolution.Len() != 0 {
		t.Errorf("a query must not populate the resolution cache, got %d entries", s.resolution.Len())
	}
}

func TestReadDNSMalformedPayload(t *testing.T) {
	s := New()
	udp := &layers.UDP{BaseLayer: layers.BaseLayer{Payload: []byte{0xde, 0xad, 0xbe}}}
	if err := s.readDNS(udp); err == nil {
		t.Error("expected an error for a truncated DNS payload")
	}
}

func TestReadDNSMultipleAnswers(t *testing.T) {
	s := New()
	msg := &dns.Msg{}
	msg.Response = true
	msg.Answer = append(msg.Answer,
		&dns.A{Hdr: dns.RR_Header{Name: "a.example.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60}, A: net.IPv4(1, 2, 3, 4).To4()},
		&dns.A{Hdr: dns.RR_Header{Name: "b.example.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60}, A: net.IPv4(5, 6, 7, 8).To4()},
	)
	if err := s.readDNS(dnsUDP(t, msg)); err != nil {
		t.Fatal(err)
	}
	if s.resolution.Len() != 2 {
		t.Errorf("got %d resolutions, want 2", s.resolution.Len())
	}
}

func TestDomainDeterministicPick(t *testing.T) {
	s := New()
	addr := netip.AddrFrom4([4]byte{93, 184, 216, 34})
	s.AddResolution(addr, "z.example.")
	s.AddResolution(addr, "a.example.")
	if got := s.domain(net.IPv4(93, 184, 216, 34)); got != "a.example." {
		t.Errorf("got %q, want %q (names are picked in sorted order)", got, "a.example.")
	}
}

func TestConcurrentResolution(t *testing.T) {
	s := New()
	addr := netip.AddrFrom4([4]byte{93, 184, 216, 34})
	w := &sync.WaitGroup{}
	for range 8 {
		w.Add(2)
		go func() {
			defer w.Done()
			for range 100 {
				s.AddResolution(addr, "example.com.")
			}
		}()
		go func() {
			defer w.Done()
			for range 100 {
				s.domain(net.IPv4(93, 184, 216, 34))
			}
		}()
	}
	w.Wait()
}
