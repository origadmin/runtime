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

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service/endpoint"
	"github.com/origadmin/runtime/service/tls"
	"github.com/origadmin/toolkits/errors"
)

const (
	Scheme   = "grpc"
	hostName = "HOST"
)

// NewServer Create a GRPC server instance
func NewServer(cfg *configv1.Service, options ...Option) (*transgrpc.Server, error) {
	if cfg == nil {
		return nil, errors.New("service config is nil")
	}
	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/grpc"))
	ll.Debugf("Creating new GRPC server instance with config: %+v", cfg)
	option := settings.ApplyZero(options)
	timeout := defaultTimeout
	serverOptions := option.ServerOptions
	if serviceGrpc := cfg.GetGrpc(); serviceGrpc != nil {
		if serviceGrpc.UseTls {
			tlsConfig, err := tls.NewServerTLSConfig(serviceGrpc.GetTlsConfig())
			if err != nil {
				return nil, err
			}
			if tlsConfig != nil {
				serverOptions = append(serverOptions, transgrpc.TLSConfig(tlsConfig))
			}
		}
		if serviceGrpc.Network != "" {
			serverOptions = append(serverOptions, transgrpc.Network(serviceGrpc.Network))
		}
		if serviceGrpc.Addr != "" {
			serverOptions = append(serverOptions, transgrpc.Address(serviceGrpc.Addr))
		}
		if serviceGrpc.Timeout != 0 {
			timeout = time.Duration(serviceGrpc.Timeout * 1e6)
		}
		if cfg.DynamicEndpoint && serviceGrpc.Endpoint == "" {
			ep := parseEndpointOption(option)
			dynamic, err := endpoint.GenerateDynamic(ep, Scheme, serviceGrpc.Addr)
			if err != nil {
				return nil, err
			}
			serviceGrpc.Endpoint = dynamic
		}
		ll.Debugw("msg", "GRPC", "endpoint", serviceGrpc.Endpoint)
		if serviceGrpc.Endpoint != "" {
			parsedEndpoint, err := url.Parse(serviceGrpc.Endpoint)
			if err == nil {
				serverOptions = append(serverOptions, transgrpc.Endpoint(parsedEndpoint))
			} else {
				ll.Errorf("Failed to parse endpoint: %v", err)
			}
		}
	}
	serverOptions = append(serverOptions, transgrpc.Timeout(timeout))
	if len(option.Middlewares) > 0 {
		serverOptions = append(serverOptions, transgrpc.Middleware(option.Middlewares...))
	}
	srv := transgrpc.NewServer(serverOptions...)
	return srv, nil
}
