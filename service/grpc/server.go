/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package grpc implements the functions, types, and interfaces for the module.
package grpc

import (
	"net/url"
	"time"

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/goexts/generic/settings"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/toolkits/env"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/helpers"
	"github.com/origadmin/toolkits/net"
)

const (
	Scheme   = "grpc"
	hostName = "HOST"
)

// NewServer Create a GRPC server instance
func NewServer(cfg *configv1.Service, ss ...OptionSetting) (*transgrpc.Server, error) {
	log.Debugf("Creating new GRPC server instance with config: %+v", cfg)
	if cfg == nil {
		return nil, errors.New("service config is nil")
	}
	option := settings.ApplyDefaultsOrZero(ss...)
	serverOptions := []transgrpc.ServerOption{
		transgrpc.Middleware(option.Middlewares...),
	}
	if len(option.ServerOptions) > 0 {
		serverOptions = append(serverOptions, option.ServerOptions...)
	}
	if serviceGrpc := cfg.GetGrpc(); serviceGrpc != nil {
		if serviceGrpc.Network != "" {
			serverOptions = append(serverOptions, transgrpc.Network(serviceGrpc.Network))
		}
		if serviceGrpc.Addr != "" {
			serverOptions = append(serverOptions, transgrpc.Address(serviceGrpc.Addr))
		}
		if serviceGrpc.Timeout != 0 {
			serverOptions = append(serverOptions, transgrpc.Timeout(time.Duration(serviceGrpc.Timeout)))
		}
		if cfg.DynamicEndpoint && serviceGrpc.Endpoint == "" {
			log.Debugf("Dynamic endpoint is enabled and endpoint is empty, generating endpoint")
			var err error
			endpointParse := helpers.ServiceEndpoint
			// Obtain an endpoint using the custom endpointURL function or the default service discovery method
			if option.EndpointFunc != nil {
				endpointParse = option.EndpointFunc
			}

			var host string
			if option.Prefix != "" {
				host = env.Var(option.Prefix, hostName)
			}
			hostIP := option.HostIp
			if hostIP == "" {
				hostIP = net.HostAddr(host)
			}
			log.Debugf("Resolved host IP: %s", hostIP)
			endpointStr, err := endpointParse("grpc", hostIP, serviceGrpc.Addr)
			if err == nil {
				serviceGrpc.Endpoint = endpointStr
			} else {
				log.Errorf("Failed to generate endpoint: %v", err)
			}
		}
		log.Debugf("GRPC endpoint: %s", serviceGrpc.Endpoint)
		if serviceGrpc.Endpoint != "" {
			endpoint, err := url.Parse(serviceGrpc.Endpoint)
			if err == nil {
				serverOptions = append(serverOptions, transgrpc.Endpoint(endpoint))
			} else {
				log.Errorf("Failed to parse endpoint: %v", err)
			}
		}
	}

	srv := transgrpc.NewServer(serverOptions...)
	return srv, nil
}
