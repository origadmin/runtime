/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"sync"

	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/registry"
	"github.com/origadmin/runtime/service"
)

// builder is a struct that holds a map of ConfigBuilders and a map of RegistryBuilders.
type builder struct {
	syncMux         sync.RWMutex
	syncs           map[string]config.Syncer
	ConfigBuilder   config.Builder
	RegistryBuilder registry.Builder
	ServiceBuilder  service.Builder
	middlewareMux   sync.RWMutex
	middlewares     map[string]MiddlewareBuilder
}

func (b *builder) NewConfig(sourceConfig *configv1.SourceConfig, setting ...config.OptionSetting) (config.KConfig, error) {
	return b.ConfigBuilder.NewConfig(sourceConfig, setting...)
}

func (b *builder) RegisterConfigBuilder(s string, factory config.Factory) {
	b.ConfigBuilder.RegisterConfigBuilder(s, factory)
}

func (b *builder) NewGRPCServer(c *configv1.Service, setting ...service.OptionSetting) (*service.GRPCServer, error) {
	return b.ServiceBuilder.NewGRPCServer(c, setting...)
}

func (b *builder) NewHTTPServer(c *configv1.Service, setting ...service.OptionSetting) (*service.HTTPServer, error) {
	return b.ServiceBuilder.NewHTTPServer(c, setting...)
}

func (b *builder) NewGRPCClient(c context.Context, c2 *configv1.Service, setting ...service.OptionSetting) (*service.GRPCClient, error) {
	return b.ServiceBuilder.NewGRPCClient(c, c2, setting...)
}

func (b *builder) NewHTTPClient(c context.Context, c2 *configv1.Service, setting ...service.OptionSetting) (*service.HTTPClient, error) {
	return b.ServiceBuilder.NewHTTPClient(c, c2, setting...)
}

func (b *builder) RegisterServiceBuilder(name string, factory service.Factory) {
	b.ServiceBuilder.RegisterServiceBuilder(name, factory)
}

// init initializes the builder struct.

func (b *builder) init() {

	b.syncs = make(map[string]config.Syncer)
	b.ConfigBuilder = config.NewBuilder()
	b.RegistryBuilder = registry.NewBuilder()
	b.ServiceBuilder = service.NewBuilder()
	b.middlewares = make(map[string]MiddlewareBuilder)
}

func newBuilder() *builder {
	b := &builder{}
	return b
}

// NewConfig creates a new Selector using the registered ConfigBuilder.
func NewConfig(cfg *configv1.SourceConfig, ss ...config.OptionSetting) (config.KConfig, error) {
	return runtime.builder.ConfigBuilder.NewConfig(cfg, ss...)
}

// RegisterConfig registers a ConfigBuilder with the builder.
func RegisterConfig(name string, factory config.Factory) {
	runtime.builder.ConfigBuilder.RegisterConfigBuilder(name, factory)
}

// RegisterConfigFunc registers a ConfigBuilder with the builder.
func RegisterConfigFunc(name string, buildFunc config.BuildFunc) {
	runtime.builder.ConfigBuilder.RegisterConfigBuilder(name, buildFunc)
}

// SyncConfig synchronizes the given configuration with the given value.
func SyncConfig(cfg *configv1.SourceConfig, v any, ss ...config.OptionSetting) error {
	return runtime.builder.SyncConfig(cfg, v, ss...)
}

func RegisterConfigSync(name string, syncFunc ConfigSyncFunc) {
	runtime.builder.RegisterConfigSync(name, syncFunc)
}

// NewDiscovery creates a new discovery using the registered RegistryBuilder.
func NewDiscovery(cfg *configv1.Registry, ss ...registry.OptionSetting) (registry.KDiscovery, error) {
	return runtime.builder.NewDiscovery(cfg, ss...)
}

// NewRegistrar creates a new KRegistrar using the registered RegistryBuilder.
func NewRegistrar(cfg *configv1.Registry, ss ...registry.OptionSetting) (registry.KRegistrar, error) {
	return runtime.builder.NewRegistrar(cfg, ss...)
}

// RegisterRegistry registers a RegistryBuilder with the builder.
func RegisterRegistry(name string, factory registry.Factory) {
	runtime.builder.RegisterRegistryBuilder(name, factory)
}

// NewMiddlewareClient creates a new KMiddleware with the builder.
func NewMiddlewareClient(name string, cm *configv1.Customize_Config, ss ...middleware.OptionSetting) (middleware.KMiddleware, error) {
	return runtime.builder.NewMiddlewareClient(name, cm, ss...)
}

// NewMiddlewareServer creates a new KMiddleware with the builder.
func NewMiddlewareServer(name string, cm *configv1.Customize_Config, ss ...middleware.OptionSetting) (middleware.KMiddleware, error) {
	return runtime.builder.NewMiddlewareServer(name, cm, ss...)
}

// NewMiddlewaresClient creates a new KMiddleware with the builder.
func NewMiddlewaresClient(cc *configv1.Customize, ss ...middleware.OptionSetting) []middleware.KMiddleware {
	return runtime.builder.NewMiddlewaresClient(nil, cc, ss...)
}

// NewMiddlewaresServer creates a new KMiddleware with the builder.
func NewMiddlewaresServer(cc *configv1.Customize, ss ...middleware.OptionSetting) []middleware.KMiddleware {
	return runtime.builder.NewMiddlewaresServer(nil, cc, ss...)
}

// RegisterMiddleware registers a MiddlewareBuilder with the builder.
func RegisterMiddleware(name string, middlewareBuilder MiddlewareBuilder) {
	runtime.builder.RegisterMiddlewareBuilder(name, middlewareBuilder)
}

// NewHTTPServiceServer creates a new HTTP server using the provided configuration
func NewHTTPServiceServer(cfg *configv1.Service, ss ...service.OptionSetting) (*service.HTTPServer, error) {
	// Call the runtime.builder.NewHTTPServer function with the provided configuration
	return runtime.builder.NewHTTPServer(cfg, ss...)
}

// NewHTTPServiceClient creates a new HTTP client using the provided context and configuration
func NewHTTPServiceClient(ctx context.Context, cfg *configv1.Service, ss ...service.OptionSetting) (*service.HTTPClient, error) {
	// Call the runtime.builder.NewHTTPClient function with the provided context and configuration
	return runtime.builder.NewHTTPClient(ctx, cfg, ss...)
}

// NewGRPCServiceServer creates a new GRPC server using the provided configuration
func NewGRPCServiceServer(cfg *configv1.Service, ss ...service.OptionSetting) (*service.GRPCServer, error) {
	// Call the runtime.builder.NewGRPCServer function with the provided configuration
	return runtime.builder.NewGRPCServer(cfg, ss...)
}

// NewGRPCServiceClient creates a new GRPC client using the provided context and configuration
func NewGRPCServiceClient(ctx context.Context, cfg *configv1.Service, ss ...service.OptionSetting) (*service.GRPCClient, error) {
	// Call the runtime.builder.NewGRPCClient function with the provided context and configuration
	return runtime.builder.NewGRPCClient(ctx, cfg, ss...)
}

// RegisterService registers a service builder with the provided name
func RegisterService(name string, factory service.Factory) {
	// Call the runtime.builder.RegisterServiceBuilder function with the provided name and service builder
	runtime.builder.ServiceBuilder.RegisterServiceBuilder(name, factory)
}

// NewBuilder creates a new Builder.
func NewBuilder() Builder {
	b := newBuilder()
	b.init()
	return b
}
