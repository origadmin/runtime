/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"net/url"

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/registry"
)

type (
	EndpointURLFunc = func(endpoint string, scheme string, host string, addr string) (*url.URL, error)
)

// SelectorOption represents a set of configuration options for a service selector.
type SelectorOption interface {
	// GRPC define the gRPC client options
	GRPC(cfg *configv1.Service_Selector) (transgrpc.ClientOption, error)
	// HTTP define the HTTP client options
	HTTP(cfg *configv1.Service_Selector) (transhttp.ClientOption, error)
}

// ServiceOption represents a set of configuration options for a service.
type ServiceOption struct {
	// Selector is an option for selecting a service instance.
	Selector SelectorOption
	// Discovery is an interface for discovering service instances.
	Discovery registry.Discovery
	// Middlewares is a list of middleware functions to be applied to the service.
	Middlewares []middleware.Middleware
	// EndpointURL is a function that generates a URL for a service endpoint.
	EndpointURL func(endpoint string, scheme string, host string, addr string) (*url.URL, error)
}

// ServiceSetting is a function that modifies a ServiceOption.
type ServiceSetting = func(config *ServiceOption)

// WithDiscovery returns a ServiceSetting that sets the Discovery field of a ServiceOption.
func WithDiscovery(discovery registry.Discovery) ServiceSetting {
	return func(config *ServiceOption) {
		// Set the Discovery field of the ServiceOption.
		config.Discovery = discovery
	}
}

// WithMiddlewares returns a ServiceSetting that sets the Middlewares field of a ServiceOption.
func WithMiddlewares(middlewares ...middleware.Middleware) ServiceSetting {
	return func(config *ServiceOption) {
		// Set the Middlewares field of the ServiceOption.
		config.Middlewares = middlewares
	}
}

// WithEndpointURL returns a ServiceSetting that sets the EndpointURL field of a ServiceOption.
func WithEndpointURL(endpoint EndpointURLFunc) ServiceSetting {
	return func(config *ServiceOption) {
		// Set the EndpointURL field of the ServiceOption.
		config.EndpointURL = endpoint
	}
}

// WithSelector returns a ServiceSetting that sets the Selector field of a ServiceOption.
func WithSelector(selector SelectorOption) ServiceSetting {
	return func(config *ServiceOption) {
		// Set the Selector field of the ServiceOption.
		config.Selector = selector
	}
}
