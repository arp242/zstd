// Package znet contains utilities for network communication.
package znet

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"zgo.at/zstd/zstring"
)

var (
	privateCIDR     []*net.IPNet
	privateCIDROnce sync.Once
)

func setupPrivateCIDR() {
	// https://en.wikipedia.org/wiki/Reserved_IP_addresses
	blocks := []string{
		"0.0.0.0/8",       // Current network; RFC6890
		"10.0.0.0/8",      // Private network; RFC1918
		"100.64.0.0/10",   // shared address space; RFC6598
		"127.0.0.1/8",     // loopback; RFC6890
		"169.254.0.0/16",  // link local address; RFC3927
		"172.16.0.0/12",   // Private network; RFC1918
		"192.0.0.0/24",    // IETF protocol assignments; RFC6890
		"192.0.2.0/24",    // TEST-NET-1 documentation and examples; RFC5737
		"192.168.0.0/16",  // Private network; RFC1918
		"192.88.99.0/24",  // IPv6 to IPv4 relay; RFC7626
		"198.18.0.0/15",   // Benchmarking tests; RFC2544
		"198.51.100.0/24", // TEST-NET-2 documentation and examples; RFC5737
		"203.0.113.0/24",  // TEST-NET-3 documentation and examples; RFC5737
		"224.0.0.0/4",     // Multicast; RFC 5771
		"240.0.0.0/4",     // Reserved (includes broadcast / 255.255.255.255); RFC 3232

		//"::/0",          // Default route
		"::/128",  // Unspecified address
		"::1/128", // Loopback

		"fc00::/7",  // Unique local address IPv6; RFC4193
		"fe80::/10", // link local address
		"ff00::/8",  // Multicast

		// TODO: these cause wrong matches; I need to look in to this.
		//"2000::/3", // Unicast
		// "2001:db8::/32", // Documentations and examples; RFC3849
		//"2002::/16", // IPv6 to IPv4 relay; RFC7626
	}

	privateCIDR = make([]*net.IPNet, 0, len(blocks))
	for _, b := range blocks {
		_, cidr, _ := net.ParseCIDR(b)
		privateCIDR = append(privateCIDR, cidr)
	}
}

// PrivateIP reports if this is a private non-public IP address.
func PrivateIP(addr net.IP) bool {
	privateCIDROnce.Do(setupPrivateCIDR)

	for _, c := range privateCIDR {
		if c.Contains(addr) {
			return true
		}
	}
	return false
}

// PrivateIPString reports if this is a private non-public IP address.
//
// This will return true for anything that is not an IP address, such as
// "example.com" or "localhost".
func PrivateIPString(ip string) bool {
	addr := net.ParseIP(RemovePort(strings.TrimSpace(ip)))
	if addr == nil { // Not an IP address?
		return true
	}
	return PrivateIP(addr)
}

// RemovePort removes the "port" part of an hostname.
//
// This only works for "host:port", and not URLs. See net.SplitHostPort.
func RemovePort(host string) string {
	shost, _, err := net.SplitHostPort(host)
	if err != nil { // Probably doesn't have a port
		return host
	}
	return shost
}

// SafeDialer is only alllowed to connect to the listed networks and ports on
// non-private addresses.
//
// Any attempt to connect to e.g. "127.0.0.1" will return an error. This is
// intended for clients that should only connect to external resources from user
// input.
//
// If the allowed lists are empty then "tcp4", "tcp6", "80", and "443" are used.
//
// The Timeout and KeepAlive are set to 30 seconds.
//
// Also see zhttputil.SafeTransport() and zhttputil.SafeClient().
func SafeDialer(allowedNets []string, allowedPorts []int) *net.Dialer {
	return &net.Dialer{
		Control: socketControl(allowedNets, allowedPorts),

		// Same defaults as net/http.DefaultTransport
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
}

func socketControl(allowedNets []string, allowedPorts []int) func(string, string, syscall.RawConn) error {
	if len(allowedNets) == 0 {
		allowedNets = []string{"tcp4", "tcp6"}
	}
	ports := []string{"80", "443"}
	if len(allowedPorts) > 0 {
		ports = make([]string, 0, len(allowedPorts))
		for _, p := range allowedPorts {
			ports = append(ports, strconv.Itoa(p))
		}
	}

	return func(network, address string, _ syscall.RawConn) error {
		if !zstring.Contains(allowedNets, network) {
			return fmt.Errorf("znet.SafeDialer: network not in allowed list %v: %q", allowedNets, network)
		}

		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return fmt.Errorf("znet.SafeDialer: invalid host/port pair: %q: %w", address, err)
		}
		if !zstring.Contains(ports, port) {
			return fmt.Errorf("znet.SafeDialer: port not in allowed list %v: %q", ports, port)
		}

		ip := net.ParseIP(host)
		if ip == nil {
			return fmt.Errorf("znet.SafeDialer: invalid IP address: %q", host)
		}
		if PrivateIP(ip) {
			return fmt.Errorf("znet.SafeDialer: not a public IP: %q", ip)
		}

		return nil
	}
}
