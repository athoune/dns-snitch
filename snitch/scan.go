package snitch

import (
	"errors"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/gookit/goutil/dump"
)

func (s *Snitch) Scan(ifaces []*net.Interface) error {
	oups := make(chan error)
	for _, iface := range ifaces {
		ips, err := InterfaceToIP(iface)
		if err != nil {
			return err
		}
		for _, addr := range ips {
			// Sanity-check that the interface has a good address.
			if addr == nil {
				return errors.New("no good IP network found")
			} else if addr.IP[0] == 127 {
				return errors.New("skipping localhost")
			} else if addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
				return errors.New("mask means network is too large")
			}
			//slog.Info("Using network", "range", addr, "interface", iface.Name)

			go func() {
				// Open up a pcap handle for packet reads/writes.
				handle, err := pcap.OpenLive(iface.Name, 65536, false, pcap.BlockForever)
				if err != nil {
					oups <- err
					return
				}
				defer handle.Close()
				err = handle.SetBPFFilter(`
			(
				(dst port 53 or src port 53)
				and proto UDP
			) or (
				proto TCP and (
					dst port 80 or dst port 443
					or src port 80 or src port 443
				)
			)`)
				if err != nil {
					oups <- err
					return
				}

				packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
				for packet := range packetSource.Packets() {
					err = s.read(packet)
					if err != nil {
						oups <- err
						return
					}
				}

			}()
		}
	}
	for {
		err := <-oups
		fmt.Println("pcap scan error :", err)
		return err
	}
	return nil
}

func (s *Snitch) read(packet gopacket.Packet) error {
	//fmt.Println(packet.Dump())
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		if udp.DstPort == 53 || udp.SrcPort == 53 {
			s.readDNS(udp)
		} else {
			dump.P(udp)
		}
		return nil
	}
	var src, dest net.IP
	ipLayer := packet.Layer(layers.LayerTypeIPv6)
	//dump.P("ip6", ipLayer)
	if ipLayer != nil {
		ip6 := ipLayer.(*layers.IPv6)
		src = ip6.SrcIP
		dest = ip6.DstIP
	} else {
		ipLayer = packet.Layer(layers.LayerTypeIPv4)
		ip4 := ipLayer.(*layers.IPv4)
		src = ip4.SrcIP
		dest = ip4.DstIP
		/*
			fmt.Println("\n\n\nip4")
			dump.P(src)
			dump.P(dest)
		*/
	}

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		err := s.readHTTPPacket(src, dest, tcp, packet)
		if err != nil {
			return errors.New(fmt.Sprint("HTTP read error", err))
		}

	}
	return nil
}
