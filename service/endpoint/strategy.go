/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package endpoint

import (
	"net"
	"sort"
)

// defaultIPStrategy is the default strategy that prefers IPv4 addresses
func defaultIPStrategy(ips []net.IP) (net.IP, error) {
	// Make a copy to avoid modifying the original slice
	ipList := make([]net.IP, len(ips))
	copy(ipList, ips)

	// Sort IPs: IPv4 first, then IPv6, and by string representation for determinism
	sort.Slice(ipList, func(i, j int) bool {
		iIsIPv4 := ipList[i].To4() != nil
		jIsIPv4 := ipList[j].To4() != nil

		// Prefer IPv4 over IPv6
		if iIsIPv4 != jIsIPv4 {
			return iIsIPv4
		}

		// Then sort by string representation
		return ipList[i].String() < ipList[j].String()
	})

	if len(ipList) > 0 {
		return ipList[0], nil
	}
	return nil, ErrNoIPFound
}

// PreferIPv4Strategy prefers IPv4 addresses over IPv6
func PreferIPv4Strategy(ips []net.IP) (net.IP, error) {
	for _, ip := range ips {
		if ip.To4() != nil {
			return ip, nil
		}
	}
	if len(ips) > 0 {
		return ips[0], nil
	}
	return nil, ErrNoIPFound
}

// PreferPublicIPStrategy prefers public IPs over private ones
func PreferPublicIPStrategy(ips []net.IP) (net.IP, error) {
	var publicIPs, privateIPs []net.IP

	for _, ip := range ips {
		if isPrivateIP(ip) {
			privateIPs = append(privateIPs, ip)
		} else {
			publicIPs = append(publicIPs, ip)
		}
	}

	// Return first public IP if available, otherwise first private IP
	if len(publicIPs) > 0 {
		return publicIPs[0], nil
	}
	if len(privateIPs) > 0 {
		return privateIPs[0], nil
	}
	return nil, ErrNoIPFound
}

// isPrivateIP checks if an IP is in a private network
func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	// Check for private address ranges
	ip4 := ip.To4()
	if ip4 == nil {
		// Not an IPv4 address
		return false
	}

	// 10.0.0.0/8
	if ip4[0] == 10 {
		return true
	}
	// 172.16.0.0/12
	if ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31 {
		return true
	}
	// 192.168.0.0/16
	if ip4[0] == 192 && ip4[1] == 168 {
		return true
	}

	return false
}
