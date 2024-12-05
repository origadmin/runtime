/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package grpc implements the functions, types, and interfaces for the module.
package grpc

import (
	"net/url"

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/origadmin/toolkits/env"
	"github.com/origadmin/toolkits/helpers"
	"github.com/origadmin/toolkits/net"

	"github.com/origadmin/runtime/config"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/middleware"
)

const (
	Scheme = "grpc"
)

// NewServer Create a GRPC server instance
func NewServer(cfg *configv1.Service, rc *config.RuntimeConfig) *transgrpc.Server {
	var options []transgrpc.ServerOption
	if rc == nil {
		rc = config.DefaultRuntimeConfig
	}

	var ms []middleware.Middleware
	ms = middleware.NewServer(cfg.GetMiddleware())
	service := rc.Service()
	if service.Middlewares != nil {
		ms = append(ms, service.Middlewares...)
	}
	options = append(options, transgrpc.Middleware(ms...))

	if serviceGrpc := cfg.GetGrpc(); serviceGrpc != nil {
		if serviceGrpc.Network != "" {
			options = append(options, transgrpc.Network(serviceGrpc.Network))
		}
		if serviceGrpc.Addr != "" {
			options = append(options, transgrpc.Address(serviceGrpc.Addr))
		}
		if serviceGrpc.Timeout != nil {
			options = append(options, transgrpc.Timeout(serviceGrpc.Timeout.AsDuration()))
		}
		if cfg.DynamicEndpoint && serviceGrpc.Endpoint == "" {
			var err error
			endpointParse := helpers.ServiceEndpoint
			// Obtain an endpoint using the custom EndpointURL function or the default service discovery method
			if service.EndpointURL != nil {
				endpointParse = service.EndpointURL
			}

			host := env.Var(rc.Bootstrap().EnvPrefix, "host")
			if cfg.Host != "" {
				host = env.Var(rc.Bootstrap().EnvPrefix, cfg.Host)
			}
			endpointStr, err := endpointParse("grpc", net.HostAddr(host), serviceGrpc.Addr)
			if err == nil {
				serviceGrpc.Endpoint = endpointStr
			}
		}
		if serviceGrpc.Endpoint != "" {
			log.Infof("GRPC endpoint: %s", serviceGrpc.Endpoint)
			endpoint, err := url.Parse(serviceGrpc.Endpoint)
			// If there are no errors, add an endpoint to options
			if err == nil {
				options = append(options, transgrpc.Endpoint(endpoint))
			} else {
				// Record errors for easy debugging
				// log.Printf("Failed to get or parse endpoint: %v", err)
			}
		}
	}

	srv := transgrpc.NewServer(options...)
	return srv
}
