/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package http implements the functions, types, and interfaces for the module.
package http

import (
	"net/url"
	"time"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goexts/generic/settings"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service/endpoint"
	"github.com/origadmin/runtime/service/tls"
	"github.com/origadmin/toolkits/env"
	"github.com/origadmin/toolkits/errors"
)

const (
	Scheme   = "http"
	hostName = "HOST"
)

// NewServer Create an HTTP server instance.
func NewServer(cfg *configv1.Service, options ...Option) (*transhttp.Server, error) {
	if cfg == nil {
		log.Errorf("Service config is nil")
		return nil, errors.New("service config is nil")
	}
	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/http"))
	ll.Debugf("Creating new HTTP server instance with config: %+v", cfg)
	option := settings.ApplyZero(options)
	timeout := defaultTimeout
	serverOptions := option.ServerOptions
	if serviceHttp := cfg.GetHttp(); serviceHttp != nil {
		if serviceHttp.UseTls {
			tlsConfig, err := tls.NewServerTLSConfig(serviceHttp.GetTlsConfig())
			if err != nil {
				return nil, err
			}
			if tlsConfig != nil {
				serverOptions = append(serverOptions, transhttp.TLSConfig(tlsConfig))
			}
		}
		if serviceHttp.Network != "" {
			serverOptions = append(serverOptions, transhttp.Network(serviceHttp.Network))
		}
		if serviceHttp.Addr != "" {
			serverOptions = append(serverOptions, transhttp.Address(serviceHttp.Addr))
		}
		if serviceHttp.Timeout != 0 {
			timeout = time.Duration(serviceHttp.Timeout * 1e6)
		}
		if cfg.DynamicEndpoint && serviceHttp.Endpoint == "" {
			hostEnv := hostName
			if option.Prefix != "" {
				hostEnv = env.Var(option.Prefix, hostName)
			}
			opts := &endpoint.Options{
				EnvVar:       hostEnv,
				HostIP:       option.HostIp,
				EndpointFunc: nil,
			}
			dynamic, err := endpoint.GenerateDynamic(opts, Scheme, serviceHttp.Addr)
			if err != nil {
				return nil, err
			}
			serviceHttp.Endpoint = dynamic
		}
		ll.Debugw("msg", "HTTP", "endpoint", serviceHttp.Endpoint)
		if serviceHttp.Endpoint != "" {
			parsedEndpoint, err := url.Parse(serviceHttp.Endpoint)
			if err == nil {
				serverOptions = append(serverOptions, transhttp.Endpoint(parsedEndpoint))
			} else {
				ll.Errorf("Failed to parse endpoint: %v", err)
			}
		}
	}
	serverOptions = append(serverOptions, transhttp.Timeout(timeout))
	if len(option.Middlewares) > 0 {
		serverOptions = append(serverOptions, transhttp.Middleware(option.Middlewares...))
	}

	srv := transhttp.NewServer(serverOptions...)
	return srv, nil
}
