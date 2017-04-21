package iptool

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPrivate(t *testing.T) {
	tool, err := New()
	if !assert.NoError(t, err) {
		return
	}

	var addrs []string
	var expected []bool

	ifAddrs, err := net.InterfaceAddrs()
	if !assert.NoError(t, err) {
		return
	}

	for _, ifAddr := range ifAddrs {
		switch t := ifAddr.(type) {
		case *net.IPNet:
			addrs = append(addrs, t.IP.String())
			expected = append(expected, true)
		}
	}

	for _, addr := range globalPrivateUseCIDRs {
		ip, _, err := net.ParseCIDR(addr)
		if !assert.NoError(t, err) {
			return
		}
		addrs = append(addrs, ip.String())
		expected = append(expected, true)
	}

	addrs = append(addrs, "67.205.132.40")
	expected = append(expected, false)
	addrs = append(addrs, "www.google.com")
	expected = append(expected, false)

	for i, addr := range addrs {
		e := expected[i]
		ipaddr, err := net.ResolveIPAddr("ip", addr)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, e, tool.IsPrivate(ipaddr), addr)
	}
}
