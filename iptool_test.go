package iptool

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsPrivate(t *testing.T) {
	tool, includesLocalInterfaces := New()
	require.True(t, includesLocalInterfaces)

	var addrs []string
	var expected []bool

	ifAddrs, err := net.InterfaceAddrs()
	require.NoError(t, err)

	for _, ifAddr := range ifAddrs {
		switch t := ifAddr.(type) {
		case *net.IPNet:
			addrs = append(addrs, t.IP.String())
			expected = append(expected, true)
		}
	}

	for _, addr := range globalPrivateUseCIDRs {
		ip, _, err := net.ParseCIDR(addr)
		require.NoError(t, err)
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
		require.NoError(t, err)
		assert.Equal(t, e, tool.IsPrivate(ipaddr), addr)
	}
}
