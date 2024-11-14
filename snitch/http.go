package snitch

import (
	"io"
	"log"
	"net"

	"github.com/athoune/dns-snitch/output"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type tcpStreamFactory struct{}

// tcpStream will handle the actual decoding of http requests.
type tcpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

func (t *tcpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	tstream := &tcpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	go tstream.run() // Important... we must guarantee that data from the reader stream is read.

	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return &tstream.r
}

func (t *tcpStream) run() {
	const SIZE = 1024 * 1024
	buff := make([]byte, SIZE)
	for {
		cpt := 0
		for {
			i, err := t.r.Read(buff)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
				break
			}
			cpt += i
		}
	}
}

func (s *Snitch) buildLine(src, dest net.IP, tcp *layers.TCP) *output.Line {
	//ts := time.Now().Truncate(s.truncated)
	l := &output.Line{
		From:   src.String(),
		Target: dest.String(),
		Domain: s.domain(dest),
	}
	if tcp.DstPort == 443 || tcp.DstPort == 80 {
		l.Port = uint32(tcp.DstPort)
		l.Direction = "UP"
	}
	if tcp.SrcPort == 443 || tcp.SrcPort == 80 {
		l.Port = uint32(tcp.SrcPort)
		l.Direction = "DOWN"
	}
	return l
}

func (s *Snitch) readHTTPPacket(src, dest net.IP, tcp *layers.TCP, packet gopacket.Packet) error {
	var err error
	var size int
	l := s.buildLine(src, dest, tcp)
	if l.Direction != "" {
		ap := packet.ApplicationLayer()
		if ap == nil {
			return nil
		}
		size = len(ap.Payload())
		for _, counter := range s.counters {
			_, err = counter.Add(*l, size)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
