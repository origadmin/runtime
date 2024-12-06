/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package http implements the functions, types, and interfaces for the module.
package http

import (
	"net/url"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/origadmin/toolkits/env"
	"github.com/origadmin/toolkits/helpers"
	"github.com/origadmin/toolkits/net"

	"github.com/origadmin/runtime/config"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
)

const (
	Scheme = "http"
)

// NewServer Create an HTTP server instance.
func NewServer(cfg *configv1.Service, rc *config.RuntimeConfig) *transhttp.Server {
	if rc == nil {
		rc = config.DefaultRuntimeConfig
	}

	service := rc.Service()
	options := []transhttp.ServerOption{
		transhttp.Middleware(service.Middlewares...),
	}
	if serviceHttp := cfg.GetHttp(); serviceHttp != nil {
		if serviceHttp.Network != "" {
			options = append(options, transhttp.Network(serviceHttp.Network))
		}
		if serviceHttp.Addr != "" {
			options = append(options, transhttp.Address(serviceHttp.Addr))
		}
		if serviceHttp.Timeout != nil {
			options = append(options, transhttp.Timeout(serviceHttp.Timeout.AsDuration()))
		}
		if cfg.DynamicEndpoint && serviceHttp.Endpoint == "" {
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
			endpointStr, err := endpointParse("http", net.HostAddr(host), serviceHttp.Addr)
			if err == nil {
				serviceHttp.Endpoint = endpointStr
			}
		}
		if serviceHttp.Endpoint != "" {
			log.Infof("HTTP endpoint: %s", serviceHttp.Endpoint)
			endpoint, err := url.Parse(serviceHttp.Endpoint)
			// If there are no errors, add an endpoint to options
			if err == nil {
				options = append(options, transhttp.Endpoint(endpoint))
			} else {
				// Record errors for easy debugging
				// log.Printf("Failed to get or parse endpoint: %v", err)
			}
		}
	}

	srv := transhttp.NewServer(options...)
	return srv
}
