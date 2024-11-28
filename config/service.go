/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"net/url"

	"github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/registry"
)

type (
	EndpointURLFunc = func(endpoint string, scheme string, host string, addr string) (*url.URL, error)
)

type ServiceOption struct {
	Discovery   registry.Discovery
	Middlewares []middleware.Middleware
	EndpointURL func(endpoint string, scheme string, host string, addr string) (*url.URL, error)
}

type ServiceSetting = func(config *ServiceOption)

func WithDiscovery(discovery registry.Discovery) ServiceSetting {
	return func(config *ServiceOption) {
		config.Discovery = discovery
	}
}

func WithMiddlewares(middlewares ...middleware.Middleware) ServiceSetting {
	return func(config *ServiceOption) {
		config.Middlewares = middlewares
	}
}

func WithEndpointURL(endpoint EndpointURLFunc) ServiceSetting {
	return func(config *ServiceOption) {
		config.EndpointURL = endpoint
	}
}
