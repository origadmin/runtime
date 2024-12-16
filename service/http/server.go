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
	hostName = "ORIGADMIN_SERVICE_HOST"
)

// NewServer Create an HTTP server instance.
func NewServer(cfg *configv1.Service, ss ...OptionSetting) (*transhttp.Server, error) {
	log.Debugf("Creating new HTTP server with config: %+v", cfg)
	if cfg == nil {
		log.Errorf("Service config is nil")
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
			log.Debugf("Generating endpoint using custom endpointURL function or default service discovery method")
			var err error
			endpointParse := helpers.ServiceEndpoint
			// Obtain an endpoint using the custom endpointURL function or the default service discovery method
			if option.EndpointFunc != nil {
				endpointParse = option.EndpointFunc
			}
			var host string
			if cfg.HostName != "" {
				host = env.Var(cfg.HostName)
			} else {
				host = env.Var(hostName)
			}
			hostIP := cfg.HostIp
			if hostIP == "" {
				hostIP = net.HostAddr(host)
			}
			log.Debugf("Resolving host IP: %s", host)
			endpointStr, err := endpointParse("http", hostIP, serviceHttp.Addr)
			if err == nil {
				log.Debugf("Generated endpoint: %s", endpointStr)
				serviceHttp.Endpoint = endpointStr
			} else {
				log.Errorf("Failed to generate endpoint: %v", err)
			}
		}
		log.Infof("HTTP endpoint: %s", serviceHttp.Endpoint)
		if serviceHttp.Endpoint != "" {
			endpoint, err := url.Parse(serviceHttp.Endpoint)
			if err == nil {
				log.Debugf("Parsed endpoint: %+v", endpoint)
				options = append(options, transhttp.Endpoint(endpoint))
			} else {
				log.Errorf("Failed to parse endpoint: %v", err)
			}
		}
	}

	srv := transhttp.NewServer(options...)
	return srv, nil
}
