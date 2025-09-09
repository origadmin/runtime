/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package endpoint

import "net"

// IPSelectionStrategy defines how to select an IP from multiple available IPs
type IPSelectionStrategy func([]net.IP) (net.IP, error)

// EndpointOptions contains configuration for endpoint generation
type EndpointOptions struct {
	// EnvVar specifies the environment variable name that may contain the host IP
	EnvVar string

	// HostIP directly specifies the host IP to use (overrides EnvVar if both are set)
	HostIP string

	// EndpointFunc allows custom endpoint generation logic
	EndpointFunc func(scheme, host, port string) (string, error)

	// IPStrategy defines how to select an IP when multiple are available
	IPStrategy IPSelectionStrategy

	// PreferPublicIP indicates whether to prefer public IPs over private ones
	PreferPublicIP bool

	// AllowedInterfaces specifies which network interfaces to consider
	AllowedInterfaces []string

	// CloudMetadataEnabled enables/disables cloud provider metadata lookup
	CloudMetadataEnabled bool
}

// DefaultOptions returns the default endpoint options
func DefaultOptions() *EndpointOptions {
	return &EndpointOptions{
		EnvVar:            "SERVICE_IP",
		IPStrategy:        defaultIPStrategy,
		PreferPublicIP:    true,
		CloudMetadataEnabled: true,
	}
}

// WithEnvVar sets the environment variable name for host IP
func (o *EndpointOptions) WithEnvVar(name string) *EndpointOptions {
	o.EnvVar = name
	return o
}

// WithHostIP directly sets the host IP
func (o *EndpointOptions) WithHostIP(ip string) *EndpointOptions {
	o.HostIP = ip
	return o
}

// WithIPStrategy sets the IP selection strategy
func (o *EndpointOptions) WithIPStrategy(strategy IPSelectionStrategy) *EndpointOptions {
	o.IPStrategy = strategy
	return o
}

// WithCloudMetadata enables/disables cloud metadata lookup
func (o *EndpointOptions) WithCloudMetadata(enabled bool) *EndpointOptions {
	o.CloudMetadataEnabled = enabled
	return o
}
