/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package service implements the functions, types, and interfaces for the module.
package service

import (
	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces/builder"
	"github.com/origadmin/runtime/service/grpc"
	"github.com/origadmin/runtime/service/http"
	"github.com/origadmin/runtime/service/selector"
)

// DefaultServiceBuilder is the default instance of the buildImpl.
var DefaultServiceBuilder = &factory{}

// ServiceBuilder is a struct that implements the buildImpl interface.
// It provides methods for creating new gRPC and HTTP servers and clients.
type factory struct{}

// NewGRPCServer creates a new gRPC server based on the provided configuration.
// It returns a pointer to the new server and an error if any.
func (f factory) NewGRPCServer(cfg *configv1.Service, ss ...GRPCOption) (*GRPCServer, error) {
	if cfg.GetSelector() != nil {
		filter, err := selector.NewFilter(cfg.GetSelector())
		if err != nil {
			return nil, err
		}
		ss = append([]GRPCOption{
			grpc.WithNodeFilter(filter),
		}, ss...)
	}
	// Create a new gRPC server using the provided configuration and options.
	return grpc.NewServer(cfg, ss...)
}

// NewHTTPServer creates a new HTTP server based on the provided configuration.
// It returns a pointer to the new server and an error if any.
func (f factory) NewHTTPServer(cfg *configv1.Service, ss ...HTTPOption) (*HTTPServer, error) {
	if cfg.GetSelector() != nil {
		filter, err := selector.NewFilter(cfg.GetSelector())
		if err != nil {
			return nil, err
		}
		ss = append([]http.Option{
			http.WithNodeFilter(filter),
		}, ss...)
	}

	// Create a new HTTP server using the provided configuration and options.
	return http.NewServer(cfg, ss...)
}

// NewGRPCClient creates a new gRPC client based on the provided context and configuration.
// It returns a pointer to the new client and an error if any.
func (f factory) NewGRPCClient(ctx context.Context, cfg *configv1.Service, ss ...GRPCOption) (*GRPCClient,
	error) {
	// Create a new gRPC client using the provided context, configuration, and options.
	return grpc.NewClient(ctx, cfg, ss...)
}

// NewHTTPClient creates a new HTTP client based on the provided context and configuration.
// It returns a pointer to the new client and an error if any.
func (f factory) NewHTTPClient(ctx context.Context, cfg *configv1.Service, ss ...HTTPOption) (*HTTPClient, error) {
	// Create a new HTTP client using the provided context, configuration, and options.
	return http.NewClient(ctx, cfg, ss...)
}

type buildImpl struct {
	builder.Builder[Factory]
}

// NewGRPCServer creates a new gRPC server based on the given ServiceConfig.
func (b *buildImpl) NewGRPCServer(cfg *configv1.Service, ss ...GRPCOption) (*GRPCServer, error) {
	serviceBuilder, ok := b.Get(cfg.Name)
	if ok {
		return serviceBuilder.NewGRPCServer(cfg, ss...)
	}
	return nil, ErrServiceNotFound
}

// NewHTTPServer creates a new HTTP server based on the given ServiceConfig.
func (b *buildImpl) NewHTTPServer(cfg *configv1.Service, ss ...HTTPOption) (*HTTPServer, error) {
	serviceBuilder, ok := b.Get(cfg.Name)
	if ok {
		return serviceBuilder.NewHTTPServer(cfg, ss...)
	}
	return nil, ErrServiceNotFound
}

// NewGRPCClient creates a new gRPC client based on the given ServiceConfig.
func (b *buildImpl) NewGRPCClient(ctx context.Context, cfg *configv1.Service, ss ...GRPCOption) (*GRPCClient, error) {
	serviceBuilder, ok := b.Get(cfg.Name)
	if ok {
		return serviceBuilder.NewGRPCClient(ctx, cfg, ss...)
	}
	return nil, ErrServiceNotFound
}

// NewHTTPClient creates a new HTTP client based on the given ServiceConfig.
func (b *buildImpl) NewHTTPClient(ctx context.Context, cfg *configv1.Service, ss ...HTTPOption) (*HTTPClient, error) {
	serviceBuilder, ok := b.Get(cfg.Name)
	if ok {
		return serviceBuilder.NewHTTPClient(ctx, cfg, ss...)
	}
	return nil, ErrServiceNotFound
}

func NewBuilder() Builder {
	return &buildImpl{
		Builder: builder.New[Factory](),
	}
}
