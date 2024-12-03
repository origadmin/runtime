/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package http implements the functions, types, and interfaces for the module.
package http

import (
	"net/url"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goexts/generic/settings"
	"github.com/origadmin/toolkits/helpers"

	"github.com/origadmin/runtime/config"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/middleware"
)

// NewServer Create an HTTP server instance.
func NewServer(cfg *configv1.Service, ss ...config.RuntimeConfigSetting) *transhttp.Server {
	var options []transhttp.ServerOption

	option := settings.Apply(&config.RuntimeConfig{}, ss).Service()
	var ms []middleware.Middleware
	ms = middleware.NewServer(cfg.GetMiddleware())
	if option.Middlewares != nil {
		ms = append(ms, option.Middlewares...)
	}
	options = append(options, transhttp.Middleware(ms...))

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
		if cfg.DynamicEndpoint {
			var endpoint *url.URL
			var err error

			// Obtain an endpoint using the custom EndpointURL function or the default service discovery method
			if option.EndpointURL != nil {
				endpoint, err = option.EndpointURL(serviceHttp.Endpoint, "http", cfg.Host, serviceHttp.Addr)
			} else {
				endpointStr := helpers.ServiceDiscoveryEndpoint(serviceHttp.Endpoint, "http", cfg.Host, serviceHttp.Addr)
				endpoint, err = url.Parse(endpointStr)
			}

			// If there are no errors, add an endpoint to options
			if err == nil {
				options = append(options, transhttp.Endpoint(endpoint))
			} else {
				// Record errors for easy debugging
				// log.Printf("Failed to get or parse endpoint: %v", err)
			}
		} else {
			endpoint, err := url.Parse(serviceHttp.Endpoint)
			if err == nil {
				options = append(options, transhttp.Endpoint(endpoint))
			}
		}
	}

	srv := transhttp.NewServer(options...)
	return srv
}
