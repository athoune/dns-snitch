package main

import (
	"log"
	"net"
	"os"
	"sync"

	"github.com/athoune/dns-snitch/resolver"
)

func listen(cb func(*net.Interface) error) {
	// Get a list of all interfaces.
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, iface := range ifaces {
		wg.Add(1)
		// Start up a scan on each interface.
		go func(iface net.Interface) {
			defer wg.Done()
			//dump.P(iface)
			if err := cb(&iface); err != nil {
				log.Printf("interface %v: %v", iface.Name, err)
			}
		}(iface)
	}
	// Wait for all interfaces' scans to complete.  They'll try to run
	// forever, but will stop on an error, so if we get past this Wait
	// it means all attempts to write have failed.
	wg.Wait()
}

func main() {
	ifaces := make([]*net.Interface, 0)
	for i := 1; i < len(os.Args); i++ {
		iface, err := net.InterfaceByName(os.Args[1])
		if err != nil {
			panic(err)
		}
		ifaces = append(ifaces, iface)
	}
	resolver := resolver.New()
	err := resolver.Scan(ifaces)
	if err != nil {
		panic(err)
	}

}
