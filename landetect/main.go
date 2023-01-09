package main

import (
	"fmt"
	"net"
	"os"

	"github.com/getlantern/golog"
	"github.com/getlantern/iptool"
)

var (
	log = golog.LoggerFor("ipdetect")
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify at least one ip address to check")
	}
	tool, ok := iptool.New()
	if !ok {
		log.Fatal("Unable to build iptool")
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
