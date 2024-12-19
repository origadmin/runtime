/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package service implements the functions, types, and interfaces for the module.
package service

import (
	"sync"

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"

	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type (
	// ServiceBuildRegistry is an interface that defines a method for registering a service builder.
	serviceBuildRegistry interface {
		RegisterServiceBuilder(name string, builder Builder)
	}
	// Builder is an interface that defines a method for creating a new service.
	Builder interface {
		NewGRPCServer(*configv1.Service, ...OptionSetting) (*GRPCServer, error)
		NewHTTPServer(*configv1.Service, ...OptionSetting) (*HTTPServer, error)
		NewGRPCClient(context.Context, *configv1.Service, ...OptionSetting) (*GRPCClient, error)
		NewHTTPClient(context.Context, *configv1.Service, ...OptionSetting) (*HTTPClient, error)
	}
)
type Service struct {
	serviceMux sync.RWMutex
	services   map[string]Builder
}

func (s *Service) RegisterServiceBuilder(name string, builder Builder) {
	s.serviceMux.Lock()
	defer s.serviceMux.Unlock()
	s.services[name] = builder
}

// NewGRPCServer creates a new gRPC server based on the given ServiceConfig.
func (s *Service) NewGRPCServer(cfg *configv1.Service, ss ...OptionSetting) (*transgrpc.Server, error) {
	s.serviceMux.RLock()
	defer s.serviceMux.RUnlock()
	if serviceBuilder, ok := s.services[cfg.Name]; ok {
		return serviceBuilder.NewGRPCServer(cfg, ss...)
	}
	return nil, ErrServiceNotFound
}

// NewHTTPServer creates a new HTTP server based on the given ServiceConfig.
func (s *Service) NewHTTPServer(cfg *configv1.Service, ss ...OptionSetting) (*transhttp.Server, error) {
	s.serviceMux.RLock()
	defer s.serviceMux.RUnlock()
	if serviceBuilder, ok := s.services[cfg.Name]; ok {
		return serviceBuilder.NewHTTPServer(cfg, ss...)
	}
	return nil, ErrServiceNotFound
}

// NewGRPCClient creates a new gRPC client based on the given ServiceConfig.
func (s *Service) NewGRPCClient(ctx context.Context, cfg *configv1.Service, ss ...OptionSetting) (*grpc.ClientConn, error) {
	s.serviceMux.RLock()
	defer s.serviceMux.RUnlock()
	if serviceBuilder, ok := s.services[cfg.Name]; ok {
		return serviceBuilder.NewGRPCClient(ctx, cfg, ss...)
	}
	return nil, ErrServiceNotFound
}

// NewHTTPClient creates a new HTTP client based on the given ServiceConfig.
func (s *Service) NewHTTPClient(ctx context.Context, cfg *configv1.Service, ss ...OptionSetting) (*transhttp.Client, error) {
	s.serviceMux.RLock()
	defer s.serviceMux.RUnlock()
	if serviceBuilder, ok := s.services[cfg.Name]; ok {
		return serviceBuilder.NewHTTPClient(ctx, cfg, ss...)
	}
	return nil, ErrServiceNotFound
}

func New() *Service {
	return &Service{
		services: make(map[string]Builder),
	}
}
