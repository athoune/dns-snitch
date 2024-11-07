package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/athoune/dns-snitch/snitch"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Some interfaces are mandatory")
		return
	}
	ifaces := make([]*net.Interface, 0)

	for i := 1; i < len(os.Args); i++ {
		iface, err := net.InterfaceByName(os.Args[1])
		if err != nil {
			panic(err)
		}
		ifaces = append(ifaces, iface)
	}
	resolver, err := snitch.New(100, 10*time.Second, "./snitch.parquet")
	if err != nil {
		panic(err)
	}
	err = resolver.Scan(ifaces)
	if err != nil {
		panic(err)
	}

}
