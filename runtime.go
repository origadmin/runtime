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
func NewConfig(cfg *configv1.SourceConfig, ss ...config.RuntimeConfigSetting) (config.Config, error) {
	return build.NewConfig(cfg, ss...)
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
func SyncConfig(cfg *configv1.SourceConfig, v any, ss ...config.RuntimeConfigSetting) error {
	return build.SyncConfig(cfg, v, ss...)
}

func RegisterConfigSync(name string, syncFunc ConfigSyncFunc) {
	build.RegisterConfigSync(name, syncFunc)
}

// NewDiscovery creates a new Discovery using the registered RegistryBuilder.
func NewDiscovery(cfg *configv1.Registry, ss ...config.RuntimeConfigSetting) (registry.Discovery, error) {
	return build.NewDiscovery(cfg, ss...)
}

// NewRegistrar creates a new Registrar using the registered RegistryBuilder.
func NewRegistrar(cfg *configv1.Registry, ss ...config.RuntimeConfigSetting) (registry.Registrar, error) {
	return build.NewRegistrar(cfg, ss...)
}

// RegisterRegistry registers a RegistryBuilder with the builder.
func RegisterRegistry(name string, registryBuilder RegistryBuilder) {
	build.RegisterRegistryBuilder(name, registryBuilder)
}

// NewMiddlewareClient creates a new Middleware with the builder.
func NewMiddlewareClient(name string, cm *configv1.Customize_Config, runtimeConfig *config.RuntimeConfig) (middleware.Middleware, error) {
	return build.NewMiddlewareClient(name, cm, runtimeConfig)
}

// NewMiddlewareServer creates a new Middleware with the builder.
func NewMiddlewareServer(name string, cm *configv1.Customize_Config, runtimeConfig *config.RuntimeConfig) (middleware.Middleware, error) {
	return build.NewMiddlewareServer(name, cm, runtimeConfig)
}

// NewMiddlewaresClient creates a new Middleware with the builder.
func NewMiddlewaresClient(cc *configv1.Customize, ss ...config.RuntimeConfigSetting) []middleware.Middleware {
	return build.NewMiddlewaresClient(nil, cc, ss...)
}

// NewMiddlewaresServer creates a new Middleware with the builder.
func NewMiddlewaresServer(cc *configv1.Customize, ss ...config.RuntimeConfigSetting) []middleware.Middleware {
	return build.NewMiddlewaresServer(nil, cc, ss...)
}

// RegisterMiddleware registers a MiddlewareBuilder with the builder.
func RegisterMiddleware(name string, middlewareBuilder MiddlewareBuilder) {
	build.RegisterMiddlewareBuilder(name, middlewareBuilder)
}

// NewHTTPServiceServer creates a new HTTP server using the provided configuration
func NewHTTPServiceServer(cfg *configv1.Service, ss ...config.RuntimeConfigSetting) (*service.HTTPServer, error) {
	// Call the build.NewHTTPServer function with the provided configuration
	return build.NewHTTPServer(cfg, ss...)
}

// NewHTTPServiceClient creates a new HTTP client using the provided context and configuration
func NewHTTPServiceClient(ctx context.Context, cfg *configv1.Service, ss ...config.RuntimeConfigSetting) (*service.HTTPClient, error) {
	// Call the build.NewHTTPClient function with the provided context and configuration
	return build.NewHTTPClient(ctx, cfg, ss...)
}

// NewGRPCServiceServer creates a new GRPC server using the provided configuration
func NewGRPCServiceServer(cfg *configv1.Service, ss ...config.RuntimeConfigSetting) (*service.GRPCServer, error) {
	// Call the build.NewGRPCServer function with the provided configuration
	return build.NewGRPCServer(cfg, ss...)
}

// NewGRPCServiceClient creates a new GRPC client using the provided context and configuration
func NewGRPCServiceClient(ctx context.Context, cfg *configv1.Service, ss ...config.RuntimeConfigSetting) (*service.GRPCClient, error) {
	// Call the build.NewGRPCClient function with the provided context and configuration
	return build.NewGRPCClient(ctx, cfg, ss...)
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
