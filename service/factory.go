/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package service implements the functions, types, and interfaces for the module.
package service

import (
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/service/grpc"
	"github.com/origadmin/runtime/service/http"
	"github.com/origadmin/runtime/service/selector"
)

// DefaultServiceFactory is the default instance of the buildImpl.
var DefaultServiceFactory = &factoryImpl{}

// ServiceBuilder is a struct that implements the buildImpl interface.
// It provides methods for creating new gRPC and HTTP servers and clients.
type factoryImpl struct{}

// NewGRPCServer creates a new gRPC server based on the provided configuration.
// It returns a pointer to the new server and an error if any.
func (f factoryImpl) NewGRPCServer(cfg *configv1.Service, ss ...GRPCOption) (*GRPCServer, error) {
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
func (f factoryImpl) NewHTTPServer(cfg *configv1.Service, ss ...HTTPOption) (*HTTPServer, error) {
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
func (f factoryImpl) NewGRPCClient(ctx context.Context, cfg *configv1.Service, ss ...GRPCOption) (*GRPCClient,
	error) {
	// Create a new gRPC client using the provided context, configuration, and options.
	return grpc.NewClient(ctx, cfg, ss...)
}

// NewHTTPClient creates a new HTTP client based on the provided context and configuration.
// It returns a pointer to the new client and an error if any.
func (f factoryImpl) NewHTTPClient(ctx context.Context, cfg *configv1.Service, ss ...HTTPOption) (*HTTPClient, error) {
	// Create a new HTTP client using the provided context, configuration, and options.
	return http.NewClient(ctx, cfg, ss...)
}

type buildImpl struct {
	factory.Registry[ServerFactory]
}

func (b *buildImpl) Build(name string, service *configv1.Service, options ...ServerOption) (transport.Server, error) {
	f, ok := b.Get(name)
	if !ok {
		return nil, ErrServiceNotFound
	}
	return f.New(service, options...)
}

func NewBuilder() ServerBuilder {
	return &buildImpl{
		Registry: factory.New[ServerFactory](),
	}
}
