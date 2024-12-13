/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package service implements the functions, types, and interfaces for the module.
package service

import (
	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/service/grpc"
	"github.com/origadmin/runtime/service/http"
	"github.com/origadmin/runtime/service/selector"
)

// DefaultServiceBuilder is the default instance of the service builder.
var DefaultServiceBuilder = &serviceBuilder{}

// ServiceBuilder is a struct that implements the service builder interface.
// It provides methods for creating new gRPC and HTTP servers and clients.
type serviceBuilder struct{}

// NewGRPCServer creates a new gRPC server based on the provided configuration.
// It returns a pointer to the new server and an error if any.
func (s serviceBuilder) NewGRPCServer(cfg *configv1.Service, ss ...OptionSetting) (*GRPCServer, error) {
	option := settings.ApplyOrZero(ss...)
	if cfg.GetSelector() != nil {
		filter, err := selector.NewFilter(cfg.GetSelector())
		if err != nil {
			return nil, err
		}
		option.grpc = append([]grpc.OptionSetting{
			grpc.WithNodeFilter(filter),
		}, option.grpc...)
	}
	// Create a new gRPC server using the provided configuration and options.
	return grpc.NewServer(cfg, option.grpc...)
}

// NewHTTPServer creates a new HTTP server based on the provided configuration.
// It returns a pointer to the new server and an error if any.
func (s serviceBuilder) NewHTTPServer(cfg *configv1.Service, ss ...OptionSetting) (*HTTPServer, error) {
	option := settings.ApplyOrZero(ss...)
	if cfg.GetSelector() != nil {
		filter, err := selector.NewFilter(cfg.GetSelector())
		if err != nil {
			return nil, err
		}
		option.http = append([]http.OptionSetting{
			http.WithNodeFilter(filter),
		}, option.http...)
	}

	// Create a new HTTP server using the provided configuration and options.
	return http.NewServer(cfg, option.http...)
}

// NewGRPCClient creates a new gRPC client based on the provided context and configuration.
// It returns a pointer to the new client and an error if any.
func (s serviceBuilder) NewGRPCClient(ctx context.Context, cfg *configv1.Service, ss ...OptionSetting) (*GRPCClient, error) {
	option := settings.ApplyOrZero(ss...)
	// Create a new gRPC client using the provided context, configuration, and options.
	return grpc.NewClient(ctx, cfg, option.grpc...)
}

// NewHTTPClient creates a new HTTP client based on the provided context and configuration.
// It returns a pointer to the new client and an error if any.
func (s serviceBuilder) NewHTTPClient(ctx context.Context, cfg *configv1.Service, ss ...OptionSetting) (*HTTPClient, error) {
	option := settings.ApplyOrZero(ss...)
	// Create a new HTTP client using the provided context, configuration, and options.
	return http.NewClient(ctx, cfg, option.http...)
}
