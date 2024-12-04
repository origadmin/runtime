/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/registry"
)

type (
	EndpointURLFunc = func(scheme string, host string, addr string) (string, error)
)

// ServiceOption represents a set of configuration options for a service.
type ServiceOption struct {
	// Discovery is an interface for discovering service instances.
	Discovery registry.Discovery
	// Middlewares is a list of middleware functions to be applied to the service.
	Middlewares []middleware.Middleware
	// EndpointURL is a function that generates a URL for a service endpoint.
	EndpointURL func(scheme string, host string, addr string) (string, error)
}

// ServiceOptionSetting is a function that modifies a ServiceOption.
type ServiceOptionSetting = func(option *ServiceOption)

// WithServiceDiscovery returns a ServiceSetting that sets the Discovery field of a ServiceOption.
func WithServiceDiscovery(discovery registry.Discovery) ServiceOptionSetting {
	return func(config *ServiceOption) {
		// Set the Discovery field of the ServiceOption.
		config.Discovery = discovery
	}
}

// WithServiceMiddlewares returns a ServiceSetting that sets the Middlewares field of a ServiceOption.
func WithServiceMiddlewares(middlewares ...middleware.Middleware) ServiceOptionSetting {
	return func(config *ServiceOption) {
		// Set the Middlewares field of the ServiceOption.
		config.Middlewares = middlewares
	}
}

// WithServiceEndpointURL returns a ServiceSetting that sets the EndpointURL field of a ServiceOption.
func WithServiceEndpointURL(endpoint EndpointURLFunc) ServiceOptionSetting {
	return func(config *ServiceOption) {
		// Set the EndpointURL field of the ServiceOption.
		config.EndpointURL = endpoint
	}
}
