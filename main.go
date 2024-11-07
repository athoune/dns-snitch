package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/athoune/dns-snitch/counter"
	"github.com/athoune/dns-snitch/output"
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
	Snitch := snitch.New()
	v, err := output.New("./snitch.parquet")
	if err != nil {
		panic(err)
	}
	Snitch.AddCounter(counter.New[*output.Line](100, 10*time.Second, v.Write))
	err = Snitch.Scan(ifaces)
	if err != nil {
		panic(err)
	}

}
