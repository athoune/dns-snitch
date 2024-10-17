package main

import (
	"net"
	"os"

	"github.com/athoune/dns-snitch/resolver"
)

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
