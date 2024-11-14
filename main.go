package main

import (
	"fmt"
	"log/slog"
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
	f, err := os.OpenFile("snitch.log", os.O_CREATE+os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	//logger := slog.New(slog.NewJSONHandler(f, nil))
	//logger := slog.Default()
	slog.SetLogLoggerLevel(slog.LevelDebug)
	//slog.SetDefault(logger)

	ifaces := make([]*net.Interface, 0)

	for i := 1; i < len(os.Args); i++ {
		iface, err := net.InterfaceByName(os.Args[i])
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
	Snitch.AddCounter(counter.New[output.Line](100, 30*time.Second, v.Write))

	term := output.NewTerm(60, 10*time.Second)
	Snitch.AddCounter(counter.New[output.Line](1000, 10*time.Second, term.Write))
	err = Snitch.Scan(ifaces)
	if err != nil {
		panic(err)
	}

}
