package main

import (
	"fmt"
	"github.com/getlantern/golog"
	"github.com/getlantern/iptool"
	"net"
	"os"
)

var (
	log = golog.LoggerFor("ipdetect")
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify at least one ip address to check")
	}
	tool, err := iptool.New()
	if err != nil {
		log.Fatal(err)
	}

	for _, addr := range os.Args[1:] {
		ipAddr, err := net.ResolveIPAddr("ip", addr)
		if err != nil {
			fmt.Printf("%v: Unable to resolve IP addr: %v\n", addr, err)
			continue
		}
		fmt.Printf("%v: %v -> %v\n", addr, ipAddr, tool.IsPrivate(ipAddr))
	}

}
