/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package http implements the functions, types, and interfaces for the module.
package http

import (
	"net/url"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goexts/generic/settings"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/toolkits/env"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/helpers"
	"github.com/origadmin/toolkits/net"
)

const (
	Scheme   = "http"
	HostName = "ORIGADMIN_RUNTIME_SERVICE_HTTP_HOST"
)

// NewServer Create an HTTP server instance.
func NewServer(cfg *configv1.Service, ss ...OptionSetting) (*transhttp.Server, error) {
	if cfg == nil {
		//bootstrap = config.DefaultRuntimeConfig
		return nil, errors.New("service config is nil")
	}
	option := settings.ApplyDefaultsOrZero(ss...)
	options := []transhttp.ServerOption{
		transhttp.Middleware(option.Middlewares...),
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
			// Obtain an endpoint using the custom endpointURL function or the default service discovery method
			endpointParse = option.EndpointFunc

			host := env.Var(HostName)
			if cfg.HostName != "" {
				host = env.Var(cfg.HostName)
			}
			hostIP := cfg.HostIp
			if hostIP == "" {
				hostIP = net.HostAddr(host)
			}

			endpointStr, err := endpointParse("http", hostIP, serviceHttp.Addr)
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
	return srv, nil
}
