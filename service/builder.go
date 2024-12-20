/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package builder implements the functions, types, and interfaces for the module.
package service

import (
	"sync"

	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/service/grpc"
	"github.com/origadmin/runtime/service/http"
	"github.com/origadmin/runtime/service/selector"
)

// DefaultServiceBuilder is the default instance of the builder builder.
var DefaultServiceBuilder = &factory{}

// ServiceBuilder is a struct that implements the builder builder interface.
// It provides methods for creating new gRPC and HTTP servers and clients.
type factory struct{}

// NewGRPCServer creates a new gRPC server based on the provided configuration.
// It returns a pointer to the new server and an error if any.
func (f factory) NewGRPCServer(cfg *configv1.Service, ss ...OptionSetting) (*GRPCServer, error) {
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
func (f factory) NewHTTPServer(cfg *configv1.Service, ss ...OptionSetting) (*HTTPServer, error) {
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
func (f factory) NewGRPCClient(ctx context.Context, cfg *configv1.Service, ss ...OptionSetting) (*GRPCClient, error) {
	option := settings.ApplyOrZero(ss...)
	// Create a new gRPC client using the provided context, configuration, and options.
	return grpc.NewClient(ctx, cfg, option.grpc...)
}

// NewHTTPClient creates a new HTTP client based on the provided context and configuration.
// It returns a pointer to the new client and an error if any.
func (f factory) NewHTTPClient(ctx context.Context, cfg *configv1.Service, ss ...OptionSetting) (*HTTPClient, error) {
	option := settings.ApplyOrZero(ss...)
	// Create a new HTTP client using the provided context, configuration, and options.
	return http.NewClient(ctx, cfg, option.http...)
}

type builder struct {
	factoryMux sync.RWMutex
	factories  map[string]Factory
}

func (s *builder) RegisterServiceBuilder(name string, factory Factory) {
	s.factoryMux.Lock()
	defer s.factoryMux.Unlock()
	s.factories[name] = factory
}

// NewGRPCServer creates a new gRPC server based on the given ServiceConfig.
func (s *builder) NewGRPCServer(cfg *configv1.Service, ss ...OptionSetting) (*GRPCServer, error) {
	s.factoryMux.RLock()
	defer s.factoryMux.RUnlock()
	if serviceBuilder, ok := s.factories[cfg.Name]; ok {
		return serviceBuilder.NewGRPCServer(cfg, ss...)
	}
	return nil, ErrServiceNotFound
}

// NewHTTPServer creates a new HTTP server based on the given ServiceConfig.
func (s *builder) NewHTTPServer(cfg *configv1.Service, ss ...OptionSetting) (*HTTPServer, error) {
	s.factoryMux.RLock()
	defer s.factoryMux.RUnlock()
	if serviceBuilder, ok := s.factories[cfg.Name]; ok {
		return serviceBuilder.NewHTTPServer(cfg, ss...)
	}
	return nil, ErrServiceNotFound
}

// NewGRPCClient creates a new gRPC client based on the given ServiceConfig.
func (s *builder) NewGRPCClient(ctx context.Context, cfg *configv1.Service, ss ...OptionSetting) (*GRPCClient, error) {
	s.factoryMux.RLock()
	defer s.factoryMux.RUnlock()
	if serviceBuilder, ok := s.factories[cfg.Name]; ok {
		return serviceBuilder.NewGRPCClient(ctx, cfg, ss...)
	}
	return nil, ErrServiceNotFound
}

// NewHTTPClient creates a new HTTP client based on the given ServiceConfig.
func (s *builder) NewHTTPClient(ctx context.Context, cfg *configv1.Service, ss ...OptionSetting) (*HTTPClient, error) {
	s.factoryMux.RLock()
	defer s.factoryMux.RUnlock()
	if serviceBuilder, ok := s.factories[cfg.Name]; ok {
		return serviceBuilder.NewHTTPClient(ctx, cfg, ss...)
	}
	return nil, ErrServiceNotFound
}

func NewBuilder() Builder {
	return &builder{
		factories: make(map[string]Factory),
	}
}
