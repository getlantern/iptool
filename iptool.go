// Package iptool provides tools for working with IP addresses.
package iptool

import (
	"net"

	"github.com/getlantern/errors"
	"github.com/getlantern/golog"
)

var (
	log = golog.LoggerFor("landetect")

	globalPrivateUseNets []*net.IPNet

	globalPrivateUseCIDRs = []string{
		// IPv4 see https://tools.ietf.org/html/rfc5735#section-3
		"0.0.0.0/8",          // "This" Network             RFC 1122, Section 3.2.1.3
		"10.0.0.0/8",         // Private-Use Networks       RFC 1918
		"127.0.0.0/8",        // Loopback                   RFC 1122, Section 3.2.1.3
		"169.254.0.0/16",     // Link Local                 RFC 3927
		"172.16.0.0/12",      // Private-Use Networks       RFC 1918
		"192.0.0.0/24",       // IETF Protocol Assignments  RFC 5736
		"192.0.2.0/24",       // TEST-NET-1                 RFC 5737
		"192.88.99.0/24",     // 6to4 Relay Anycast         RFC 3068
		"192.168.0.0/16",     // Private-Use Networks       RFC 1918
		"198.18.0.0/15",      // Network Interconnect Device Benchmark Testing   RFC 2544
		"198.51.100.0/24",    // TEST-NET-2                 RFC 5737
		"203.0.113.0/24",     // TEST-NET-3                 RFC 5737
		"224.0.0.0/4",        // Multicast                  RFC 3171
		"240.0.0.0/4",        // Reserved for Future Use    RFC 1112, Section 4
		"255.255.255.255/32", // Limited Broadcast          RFC 919, Section 7

		// IPv6 see https://tools.ietf.org/html/rfc5156
		"::1/128", // node-scoped unicast
		"::/128",  // node-scoped unicast
		// "::FFFF:0:0/96", // IPv4 mapped addresses
		"fe80::/10",     // link-local unicast
		"fc00::/7",      // unique local
		"2001:db8::/32", // documentation
		"2001:10::/28",  // ORCHID addresses
		"::/0",          // default route
		"ff00::/8",      // multicast
	}
)

func init() {
	// initialize reserved private network ranges
	for _, cidr := range globalPrivateUseCIDRs {
		_, privateNet, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Fatalf("Unable to parse CIDR %v: %v", cidr, err)
		}
		globalPrivateUseNets = append(globalPrivateUseNets, privateNet)
	}
}

type Tool interface {
	// IsPrivate checks whether the given IP address is private, meaning it's
	// using one of addresses designed by IANA as not routable on the Internet or
	// the address of one of the interfaces on the current host.
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
	privateNets := make([]*net.IPNet, len(globalPrivateUseNets), len(globalPrivateUseNets)+len(addrs))
	copy(privateNets, globalPrivateUseNets)
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
			log.Debugf("%v contains %v", privateNet, addr.IP)
			return true
		}
	}
	return false
}
