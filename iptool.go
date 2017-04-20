// Package iptool provides tools for working with IP addresses.
package iptool

import (
	"net"

	"github.com/getlantern/errors"
	"github.com/getlantern/golog"
)

var (
	log = golog.LoggerFor("landetect")

	globalPrivateNets []*net.IPNet

	globalPrivateCIDRs = []string{
		"10.0.0.0/8",     // reserved private
		"172.16.0.0/12",  // reserved private
		"192.168.0.0/16", // reserved private
		"127.0.0.1/32",   // loopback
		"169.254.0.0/16", // link-local
		"fc00::/7",       // reserved private
		"fe80::/10",      // link-local
	}
)

func init() {
	// initialize reserved private network ranges
	for _, cidr := range globalPrivateCIDRs {
		_, privateNet, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Fatalf("Unable to parse CIDR %v: %v", cidr, err)
		}
		globalPrivateNets = append(globalPrivateNets, privateNet)
	}
}

type Tool interface {
	// IsPrivate checks whether the given IP address is private, meaning it's
	// using one of the commonly reserved private address ranges or points
	// specifically at the address of one of the interfaces on this device.
	IsPrivate(addr *net.IPAddr) bool
}

type tool struct {
	privateNets []*net.IPNet
}

func New() (Tool, error) {
	// Build comprehensive list of private networks by combining interfaces with
	// global private networks.
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, errors.New("Unable to determine interface addresses: %v", err)
	}
	privateNets := make([]*net.IPNet, len(globalPrivateNets), len(globalPrivateNets)+len(addrs))
	copy(privateNets, globalPrivateNets)
	for _, addr := range addrs {
		switch t := addr.(type) {
		case *net.IPNet:
			privateNets = append(privateNets, t)
		}
	}
	return &tool{
		privateNets: privateNets,
	}, nil
}

func (t *tool) IsPrivate(addr *net.IPAddr) bool {
	for _, privateNet := range t.privateNets {
		if privateNet.Contains(addr.IP) {
			return true
		}
	}
	return false
}
