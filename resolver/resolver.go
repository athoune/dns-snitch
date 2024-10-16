package resolver

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/netip"
	"os"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/gookit/goutil/dump"
	"github.com/hashicorp/go-set"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/miekg/dns"
)

type Resolver struct {
	resolution *lru.Cache[netip.Addr, *set.Set[string]]
	mutex      *sync.RWMutex
}

func New() *Resolver {
	l, _ := lru.New[netip.Addr, *set.Set[string]](256)
	return &Resolver{
		resolution: l,
		mutex:      &sync.RWMutex{},
	}
}

func InterfaceToIP(iface *net.Interface) ([]*net.IPNet, error) {
	ips := make([]*net.IPNet, 0)

	var addr *net.IPNet
	if addrs, err := iface.Addrs(); err != nil {
		return nil, err
	} else {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					addr = &net.IPNet{
						IP:   ip4,
						Mask: ipnet.Mask[len(ipnet.Mask)-4:],
					}
					ips = append(ips, addr)
					continue
				}
				if ip6 := ipnet.IP.To16(); ip6 != nil {
					addr = &net.IPNet{
						IP:   ip6,
						Mask: ipnet.Mask[len(ipnet.Mask)-16:],
					}
					ips = append(ips, addr)
				}
			}
		}
	}
	return ips, nil
}

func (r *Resolver) Scan(iface *net.Interface) error {
	ips, err := InterfaceToIP(iface)
	if err != nil {
		return err
	}
	oups := make(chan error)
	for _, addr := range ips {
		// Sanity-check that the interface has a good address.
		if addr == nil {
			return errors.New("no good IP network found")
		} else if addr.IP[0] == 127 {
			return errors.New("skipping localhost")
		} else if addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
			return errors.New("mask means network is too large")
		}
		log.Printf("Using network range %v for interface %v", addr, iface.Name)

		go func() {
			// Open up a pcap handle for packet reads/writes.
			handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
			if err != nil {
				oups <- err
				return
			}
			err = handle.SetBPFFilter("((dst port 53 or src port 53) and proto UDP) or (proto TCP and ( dst port 80 or dst port 443))")
			if err != nil {
				oups <- err
				return
			}
			defer handle.Close()

			for {
				data, _, err := handle.ReadPacketData()
				if err != nil {
					fmt.Println("Read packet error", err)
					continue
				}
				r.read(data)
			}
		}()
	}
	go func() {
		for {
			time.Sleep(5 * time.Second)
			r.Dump(os.Stdout)
		}
	}()
	for {
		err = <-oups
		fmt.Println(err)
	}
	return nil
}

func (r *Resolver) Dump(out io.Writer) {
	r.mutex.RLock()
	ips := r.resolution.Keys()
	values := r.resolution.Values()
	r.mutex.RUnlock()
	fmt.Fprint(out, "\n---\n")
	for i, ip := range ips {
		fmt.Fprintf(out, "%s %s\n", ip, values[i])
	}
}

func (r *Resolver) read(myPacketData []byte) {
	packet := gopacket.NewPacket(myPacketData, layers.LayerTypeEthernet, gopacket.Default)
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		if udp.DstPort == 53 || udp.SrcPort == 53 {
			r.readDNS(udp)
		} else {
			dump.P(udp)
		}
		return
	}
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
	}
}

func (r *Resolver) Append(addr netip.Addr, domain string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	value, ok := r.resolution.Get(addr)
	if !ok {
		value = set.New[string](3)
	}
	value.Insert(domain)
	r.resolution.Add(addr, value)
}

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
