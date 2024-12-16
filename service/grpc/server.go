/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package grpc implements the functions, types, and interfaces for the module.
package grpc

import (
	"net/url"

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
	hostName = "ORIGADMIN_SERVICE_HOST"
)

// NewServer Create a GRPC server instance
func NewServer(cfg *configv1.Service, ss ...OptionSetting) (*transgrpc.Server, error) {
	log.Debugf("Creating new GRPC server instance with config: %+v", cfg)
	if cfg == nil {
		log.Errorf("Service config is nil")
		return nil, errors.New("service config is nil")
	}
	log.Debugf("Applying default settings to options: %+v", ss)
	option := settings.ApplyDefaultsOrZero(ss...)
	log.Debugf("Applied options: %+v", option)
	options := []transgrpc.ServerOption{
		transgrpc.Middleware(option.Middlewares...),
	}
	log.Debugf("Initial options: %+v", options)
	if serviceGrpc := cfg.GetGrpc(); serviceGrpc != nil {
		log.Debugf("Found GRPC config: %+v", serviceGrpc)
		if serviceGrpc.Network != "" {
			log.Debugf("Adding network option: %s", serviceGrpc.Network)
			options = append(options, transgrpc.Network(serviceGrpc.Network))
		}
		if serviceGrpc.Addr != "" {
			log.Debugf("Adding address option: %s", serviceGrpc.Addr)
			options = append(options, transgrpc.Address(serviceGrpc.Addr))
		}
		if serviceGrpc.Timeout != nil {
			log.Debugf("Adding timeout option: %s", serviceGrpc.Timeout.AsDuration())
			options = append(options, transgrpc.Timeout(serviceGrpc.Timeout.AsDuration()))
		}
		if cfg.DynamicEndpoint && serviceGrpc.Endpoint == "" {
			log.Debugf("Dynamic endpoint is enabled and endpoint is empty, generating endpoint")
			var err error
			endpointParse := helpers.ServiceEndpoint
			// Obtain an endpoint using the custom endpointURL function or the default service discovery method
			if option.EndpointFunc != nil {
				endpointParse = option.EndpointFunc
				log.Debugf("Using custom endpoint function: %+v", option.EndpointFunc)
			}

			var host string
			if cfg.HostName != "" {
				host = env.Var(cfg.HostName)
				log.Debugf("Using host name from config: %s", host)
			} else {
				host = env.Var(hostName)
				log.Debugf("Using default host name: %s", host)
			}
			hostIP := cfg.HostIp
			if hostIP == "" {
				hostIP = net.HostAddr(host)
				log.Debugf("Resolved host IP: %s", hostIP)
			}
			endpointStr, err := endpointParse("grpc", hostIP, serviceGrpc.Addr)
			if err != nil {
				log.Errorf("Failed to generate endpoint: %v", err)
			} else {
				serviceGrpc.Endpoint = endpointStr
				log.Debugf("Generated endpoint: %s", serviceGrpc.Endpoint)
			}
		}
		log.Infof("GRPC endpoint: %s", serviceGrpc.Endpoint)
		if serviceGrpc.Endpoint != "" {
			endpoint, err := url.Parse(serviceGrpc.Endpoint)
			if err != nil {
				log.Errorf("Failed to parse endpoint: %v", err)
			} else {
				log.Debugf("Parsed endpoint: %+v", endpoint)
				options = append(options, transgrpc.Endpoint(endpoint))
			}
		}
	}

	log.Debugf("Final options: %+v", options)
	srv := transgrpc.NewServer(options...)
	log.Debugf("Created new GRPC server instance: %+v", srv)
	return srv, nil
}
