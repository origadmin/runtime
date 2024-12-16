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
	log.Debugf("Applying default settings to options: %+v", ss)
	option := settings.ApplyDefaultsOrZero(ss...)
	log.Debugf("Applied options: %+v", option)

	options := []transhttp.ServerOption{
		transhttp.Middleware(option.Middlewares...),
	}
	log.Debugf("Initial server options: %+v", options)

	if serviceHttp := cfg.GetHttp(); serviceHttp != nil {
		log.Debugf("Configured HTTP service: %+v", serviceHttp)
		if serviceHttp.Network != "" {
			log.Debugf("Setting network to: %s", serviceHttp.Network)
			options = append(options, transhttp.Network(serviceHttp.Network))
		}
		if serviceHttp.Addr != "" {
			log.Debugf("Setting address to: %s", serviceHttp.Addr)
			options = append(options, transhttp.Address(serviceHttp.Addr))
		}
		if serviceHttp.Timeout != nil {
			log.Debugf("Setting timeout to: %s", serviceHttp.Timeout.AsDuration())
			options = append(options, transhttp.Timeout(serviceHttp.Timeout.AsDuration()))
		}
		if cfg.DynamicEndpoint && serviceHttp.Endpoint == "" {
			log.Debugf("Generating endpoint using custom endpointURL function or default service discovery method")
			var err error
			endpointParse := helpers.ServiceEndpoint
			// Obtain an endpoint using the custom endpointURL function or the default service discovery method
			endpointParse = option.EndpointFunc

			var host string
			if cfg.HostName != "" {
				log.Debugf("Using hostname: %s", cfg.HostName)
				host = env.Var(cfg.HostName)
			} else {
				log.Debugf("Using default hostname: %s", hostName)
				host = env.Var(hostName)
			}
			hostIP := cfg.HostIp
			if hostIP == "" {
				log.Debugf("Resolving host IP: %s", host)
				hostIP = net.HostAddr(host)
			}

			endpointStr, err := endpointParse("http", hostIP, serviceHttp.Addr)
			if err != nil {
				log.Errorf("Failed to generate endpoint: %v", err)
			} else {
				log.Debugf("Generated endpoint: %s", endpointStr)
				serviceHttp.Endpoint = endpointStr
			}
		}
		log.Infof("HTTP endpoint: %s", serviceHttp.Endpoint)
		if serviceHttp.Endpoint != "" {
			endpoint, err := url.Parse(serviceHttp.Endpoint)
			if err != nil {
				log.Errorf("Failed to parse endpoint: %v", err)
			} else {
				log.Debugf("Parsed endpoint: %+v", endpoint)
				// If there are no errors, add an endpoint to options
				options = append(options, transhttp.Endpoint(endpoint))
			}
		}
	}

	log.Debugf("Final server options: %+v", options)
	srv := transhttp.NewServer(options...)
	log.Debugf("Created new HTTP server: %+v", srv)
	return srv, nil
}
