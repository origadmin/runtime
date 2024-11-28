/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"

	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/service"
)

type (
	// ServiceBuildRegistry is an interface that defines a method for registering a service builder.
	serviceBuildRegistry interface {
		RegisterServiceBuilder(name string, builder ServiceBuilder)
	}
	// ServiceBuilder is an interface that defines a method for creating a new service.
	ServiceBuilder interface {
		NewGRPCServer(cfg *configv1.Service, opts ...config.ServiceSetting) (*service.GRPCServer, error)
		NewHTTPServer(cfg *configv1.Service, opts ...config.ServiceSetting) (*service.HTTPServer, error)
		NewGRPCClient(ctx context.Context, cfg *configv1.Service, opts ...config.ServiceSetting) (*service.GRPCClient, error)
		NewHTTPClient(ctx context.Context, cfg *configv1.Service, opts ...config.ServiceSetting) (*service.HTTPClient, error)
	}
)

// NewGRPCServer creates a new gRPC server based on the given ServiceConfig.
func (b *builder) NewGRPCServer(cfg *configv1.Service, opts ...config.ServiceSetting) (*transgrpc.Server, error) {
	b.serviceMux.RLock()
	defer b.serviceMux.RUnlock()
	if serviceBuilder, ok := b.services[cfg.Name]; ok {
		return serviceBuilder.NewGRPCServer(cfg, opts...)
	}
	return nil, ErrNotFound
}

// NewHTTPServer creates a new HTTP server based on the given ServiceConfig.
func (b *builder) NewHTTPServer(cfg *configv1.Service, opts ...config.ServiceSetting) (*transhttp.Server, error) {
	b.serviceMux.RLock()
	defer b.serviceMux.RUnlock()
	if serviceBuilder, ok := b.services[cfg.Name]; ok {
		return serviceBuilder.NewHTTPServer(cfg, opts...)
	}
	return nil, ErrNotFound
}

// NewGRPCClient creates a new gRPC client based on the given ServiceConfig.
func (b *builder) NewGRPCClient(ctx context.Context, cfg *configv1.Service, opts ...config.ServiceSetting) (*grpc.ClientConn, error) {
	b.serviceMux.RLock()
	defer b.serviceMux.RUnlock()
	if serviceBuilder, ok := b.services[cfg.Name]; ok {
		return serviceBuilder.NewGRPCClient(ctx, cfg, opts...)
	}
	return nil, ErrNotFound
}

// NewHTTPClient creates a new HTTP client based on the given ServiceConfig.
func (b *builder) NewHTTPClient(ctx context.Context, cfg *configv1.Service, opts ...config.ServiceSetting) (*transhttp.Client, error) {
	b.serviceMux.RLock()
	defer b.serviceMux.RUnlock()
	if serviceBuilder, ok := b.services[cfg.Name]; ok {
		return serviceBuilder.NewHTTPClient(ctx, cfg, opts...)
	}
	return nil, ErrNotFound
}

// RegisterServiceBuilder registers a new ServiceBuilder with the given service name.
func (b *builder) RegisterServiceBuilder(name string, builder ServiceBuilder) {
	b.serviceMux.Lock()
	defer b.serviceMux.Unlock()
	b.services[name] = builder
}
