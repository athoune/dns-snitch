package snitch

import (
	"testing"

	"github.com/google/gopacket/layers"
)

func TestBuildLine(t *testing.T) {
	s := New()
	tcp := &layers.TCP{}
	tcp.DstPort = 443
	line := s.buildLine([]byte{192, 168, 1, 1}, []byte{212, 27, 48, 10}, tcp)
	if line.From != "192.168.1.1" {
		t.Error("what the from", line.From)
	}
	if line.Target != "212.27.48.10" {
		t.Error("what the destination", line.Target)
	}
	if line.Direction != "UP" {
		t.Error("what the direction", line.Direction)
	}
	if line.Domain != "www.free.fr." {
		t.Error("What the domaine", line.Domain)
	}
	if line.Port != 443 {
		t.Error("what the port", line.Port)
	}

	tcp.SrcPort = 443
	line = s.buildLine([]byte{212, 27, 48, 10}, []byte{192, 168, 1, 1}, tcp)
	if line.From != "212.27.48.10" {
		t.Error("what the from", line.From)
	}
	if line.Target != "192.168.1.1" {
		t.Error("what the destination", line.Target)
	}
	if line.Direction != "DOWN" {
		t.Error("what the direction", line.Direction)
	}
	if line.Domain != "www.free.fr." {
		t.Error("What the domaine", line.Domain)
	}

}
