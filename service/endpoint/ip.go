/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package endpoint implements the functions, types, and interfaces for the module.
package endpoint

import (
	"errors"
	"fmt"
	"net"

	toolnet "github.com/origadmin/toolkits/net"
)

func ExtractIP(scheme, host, hostPort string) (string, error) {
	_, port, err := net.SplitHostPort(hostPort)
	if err != nil && host == "" {
		return "", fmt.Errorf("invalid host: %w", err)
	}
	if len(host) > 0 && (host != "0.0.0.0" && host != "[::]" && host != "::") {
		return schemeHost(scheme, net.JoinHostPort(host, port)), nil
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get local ip: %w", err)
	}
	minIndex := int(^uint(0) >> 1)
	ips := make([]net.IP, 0)
	for _, iface := range ifaces {
		if (iface.Flags & net.FlagUp) == 0 {
			continue
		}
		if iface.Index >= minIndex && len(ips) != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for i, rawAddr := range addrs {
			var ip net.IP
			switch addr := rawAddr.(type) {
			case *net.IPAddr:
				ip = addr.IP
			case *net.IPNet:
				ip = addr.IP
			default:
				continue
			}
			if isValidIP(ip.String()) {
				minIndex = iface.Index
				if i == 0 {
					ips = make([]net.IP, 0, 1)
				}
				ips = append(ips, ip)
				if ip.To4() != nil {
					break
				}
			}
		}
	}
	if len(ips) != 0 {
		return schemeHost(scheme, net.JoinHostPort(ips[len(ips)-1].String(), port)), nil
	}
	return "", errors.New("no local ip found")
}

func schemeHost(scheme, host string) string {
	return scheme + "://" + host
}

func isValidIP(addr string) bool {
	ip := net.ParseIP(addr)
	return ip.IsGlobalUnicast() && !ip.IsInterfaceLocalMulticast()
}

type Options struct {
	EnvVar       string
	HostIP       string
	EndpointFunc func(scheme string, host string, addr string) (string, error)
}

func GenerateDynamic(option *Options, scheme, addr string) (string, error) {
	endpointParse := ExtractIP
	if option.EndpointFunc != nil {
		endpointParse = option.EndpointFunc
	}

	hostIP := option.HostIP
	if hostIP == "" {
		hostIP = toolnet.HostAddr(toolnet.WithEnvVar(option.EnvVar))
	}
	return endpointParse(scheme, hostIP, addr)
}
