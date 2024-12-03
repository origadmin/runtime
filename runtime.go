/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime provides functions for loading configurations and registering services.
package runtime

import (
	"sync"

	"github.com/origadmin/toolkits/errors"

	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/registry"
	"github.com/origadmin/runtime/service"
)

type Builder interface {
	ConfigBuilder
	RegistryBuilder
	ServiceBuilder
	MiddlewareBuilders

	configBuildRegistry
	registryBuildRegistry
	serviceBuildRegistry
	middlewareBuildRegistry
}

// build is a global variable that holds an instance of the builder struct.
var (
	once  = &sync.Once{}
	build = &builder{}
)

// ErrNotFound is an error that is returned when a ConfigBuilder or RegistryBuilder is not found.
var ErrNotFound = errors.String("not found")

// init initializes the builder struct.
func init() {
	once.Do(func() {
		build.init()
	})
}

// Global returns the global instance of the builder.
func Global() Builder {
	return build
}

// NewConfig creates a new Selector using the registered ConfigBuilder.
func NewConfig(cfg *configv1.SourceConfig, rc *config.RuntimeConfig) (config.Config, error) {
	if rc == nil {
		rc = config.DefaultRuntimeConfig
	}
	return build.NewConfig(cfg, rc)
}

// RegisterConfig registers a ConfigBuilder with the builder.
func RegisterConfig(name string, configBuilder ConfigBuilder) {
	build.RegisterConfigBuilder(name, configBuilder)
}

// RegisterConfigFunc registers a ConfigBuilder with the builder.
func RegisterConfigFunc(name string, buildFunc ConfigBuildFunc) {
	build.RegisterConfigBuilder(name, buildFunc)
}

// SyncConfig synchronizes the given configuration with the given value.
func SyncConfig(cfg *configv1.SourceConfig, v any, rc *config.RuntimeConfig) error {
	if rc == nil {
		rc = config.DefaultRuntimeConfig
	}
	return build.SyncConfig(cfg, v, rc)
}

func RegisterConfigSync(name string, syncFunc ConfigSyncFunc) {
	build.RegisterConfigSync(name, syncFunc)
}

// NewDiscovery creates a new Discovery using the registered RegistryBuilder.
func NewDiscovery(cfg *configv1.Registry, rc *config.RuntimeConfig) (registry.Discovery, error) {
	if rc == nil {
		rc = config.DefaultRuntimeConfig
	}
	return build.NewDiscovery(cfg, rc)
}

// NewRegistrar creates a new Registrar using the registered RegistryBuilder.
func NewRegistrar(cfg *configv1.Registry, rc *config.RuntimeConfig) (registry.Registrar, error) {
	if rc == nil {
		rc = config.DefaultRuntimeConfig
	}
	return build.NewRegistrar(cfg, rc)
}

// RegisterRegistry registers a RegistryBuilder with the builder.
func RegisterRegistry(name string, registryBuilder RegistryBuilder) {
	build.RegisterRegistryBuilder(name, registryBuilder)
}

// NewMiddlewareClient creates a new Middleware with the builder.
func NewMiddlewareClient(name string, cm *configv1.Customize_Config, rc *config.RuntimeConfig) (middleware.Middleware, error) {
	if rc == nil {
		rc = config.DefaultRuntimeConfig
	}
	return build.NewMiddlewareClient(name, cm, rc)
}

// NewMiddlewareServer creates a new Middleware with the builder.
func NewMiddlewareServer(name string, cm *configv1.Customize_Config, rc *config.RuntimeConfig) (middleware.Middleware, error) {
	if rc == nil {
		rc = config.DefaultRuntimeConfig
	}
	return build.NewMiddlewareServer(name, cm, rc)
}

// NewMiddlewaresClient creates a new Middleware with the builder.
func NewMiddlewaresClient(cc *configv1.Customize, rc *config.RuntimeConfig) []middleware.Middleware {
	if rc == nil {
		rc = config.DefaultRuntimeConfig
	}
	return build.NewMiddlewaresClient(nil, cc, rc)
}

// NewMiddlewaresServer creates a new Middleware with the builder.
func NewMiddlewaresServer(cc *configv1.Customize, rc *config.RuntimeConfig) []middleware.Middleware {
	if rc == nil {
		rc = config.DefaultRuntimeConfig
	}
	return build.NewMiddlewaresServer(nil, cc, rc)
}

// RegisterMiddleware registers a MiddlewareBuilder with the builder.
func RegisterMiddleware(name string, middlewareBuilder MiddlewareBuilder) {
	build.RegisterMiddlewareBuilder(name, middlewareBuilder)
}

// NewHTTPServiceServer creates a new HTTP server using the provided configuration
func NewHTTPServiceServer(cfg *configv1.Service, rc *config.RuntimeConfig) (*service.HTTPServer, error) {
	if rc == nil {
		rc = config.DefaultRuntimeConfig
	}
	// Call the build.NewHTTPServer function with the provided configuration
	return build.NewHTTPServer(cfg, rc)
}

// NewHTTPServiceClient creates a new HTTP client using the provided context and configuration
func NewHTTPServiceClient(ctx context.Context, cfg *configv1.Service, rc *config.RuntimeConfig) (*service.HTTPClient, error) {
	if rc == nil {
		rc = config.DefaultRuntimeConfig
	}
	// Call the build.NewHTTPClient function with the provided context and configuration
	return build.NewHTTPClient(ctx, cfg, rc)
}

// NewGRPCServiceServer creates a new GRPC server using the provided configuration
func NewGRPCServiceServer(cfg *configv1.Service, rc *config.RuntimeConfig) (*service.GRPCServer, error) {
	if rc == nil {
		rc = config.DefaultRuntimeConfig
	}
	// Call the build.NewGRPCServer function with the provided configuration
	return build.NewGRPCServer(cfg, rc)
}

// NewGRPCServiceClient creates a new GRPC client using the provided context and configuration
func NewGRPCServiceClient(ctx context.Context, cfg *configv1.Service, rc *config.RuntimeConfig) (*service.GRPCClient, error) {
	if rc == nil {
		rc = config.DefaultRuntimeConfig
	}
	// Call the build.NewGRPCClient function with the provided context and configuration
	return build.NewGRPCClient(ctx, cfg, rc)
}

// RegisterService registers a service builder with the provided name
func RegisterService(name string, serviceBuilder ServiceBuilder) {
	// Call the build.RegisterServiceBuilder function with the provided name and service builder
	build.RegisterServiceBuilder(name, serviceBuilder)
}

// New creates a new Builder.
func New() Builder {
	return newBuilder()
}
