/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"sync"

	"github.com/go-kratos/kratos/v2/transport"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/registry"
	"github.com/origadmin/runtime/service"
	"github.com/origadmin/runtime/service/grpc"
	"github.com/origadmin/runtime/service/http"
)

type Builder interface {
	Config() config.Builder
	Registry() registry.Builder
	Service() service.ServerBuilder
	Middleware() middleware.Builder
	NewMiddlewareClient(name string, config *middlewarev1.Middleware, ss ...middleware.Option) (middleware.KMiddleware, error)
	NewMiddlewareServer(name string, config *middlewarev1.Middleware, ss ...middleware.Option) (middleware.KMiddleware, error)
	NewMiddlewaresClient(config *middlewarev1.Middleware, ss ...middleware.Option) []middleware.KMiddleware
	NewMiddlewaresServer(config *middlewarev1.Middleware, ss ...middleware.Option) []middleware.KMiddleware
	RegisterMiddlewareBuilder(name string, builder middleware.Factory)
	NewRegistrar(cfg *configv1.Discovery, ss ...registry.Option) (registry.KRegistrar, error)
	NewDiscovery(cfg *configv1.Discovery, ss ...registry.Option) (registry.KDiscovery, error)
	RegisterRegistryBuilder(name string, factory registry.Factory)
	RegisterRegistryFunc(name string, registryBuilder registry.RegistrarBuildFunc, discoveryBuilder registry.DiscoveryBuildFunc)
	NewConfig(sourceConfig *configv1.SourceConfig, options ...config.Option) (config.KConfig, error)
	RegisterConfigBuilder(s string, factory config.Factory)
	NewServer(name string, c *configv1.Service, options ...service.ServerOption) (transport.Server, error)
	NewGRPCServer(c *configv1.Service, options ...service.GRPCOption) (*service.GRPCServer, error)
	NewHTTPServer(c *configv1.Service, options ...service.HTTPOption) (*service.HTTPServer, error)
	NewGRPCClient(c context.Context, c2 *configv1.Service, options ...service.GRPCOption) (*service.GRPCClient, error)
	NewHTTPClient(c context.Context, c2 *configv1.Service, options ...service.HTTPOption) (*service.HTTPClient, error)
	RegisterServiceBuilder(name string, factory service.ServerFactory)
	SyncConfig(cfg *configv1.SourceConfig, v any, ss ...config.Option) error
	RegisterConfigSyncer(name string, configSyncer config.Syncer)
	RegisterConfigSync(name string, configSyncer config.Syncer)
}

// builder is a struct that holds a map of ConfigBuilders and a map of RegistryBuilders.
type builder struct {
	syncMux           sync.RWMutex
	syncs             map[string]config.Syncer
	ConfigBuilder     config.Builder
	RegistryBuilder   registry.Builder
	ServiceBuilder    service.ServerBuilder
	MiddlewareBuilder middleware.Builder
}

func (b *builder) Config() config.Builder {
	return b.ConfigBuilder
}

func (b *builder) Registry() registry.Builder {
	return b.RegistryBuilder
}

func (b *builder) Service() service.ServerBuilder {
	return b.ServiceBuilder
}

func (b *builder) Middleware() middleware.Builder {
	return b.MiddlewareBuilder
}

func (b *builder) NewConfig(sourceConfig *configv1.SourceConfig, options ...config.Option) (config.KConfig, error) {
	return b.ConfigBuilder.NewConfig(sourceConfig, options...)
}

func (b *builder) RegisterConfigBuilder(s string, factory config.Factory) {
	b.ConfigBuilder.Register(s, factory)
}

func (b *builder) NewServer(name string, c *configv1.Service, options ...service.ServerOption) (transport.Server, error) {
	return b.ServiceBuilder.Build(name, c, options...)
}

func (b *builder) NewGRPCServer(cfg *configv1.Service, options ...service.GRPCOption) (*service.GRPCServer, error) {
	return grpc.NewServer(cfg, options...)
}

func (b *builder) NewHTTPServer(cfg *configv1.Service, options ...service.HTTPOption) (*service.HTTPServer, error) {
	return http.NewServer(cfg, options...)
}

func (b *builder) NewGRPCClient(c context.Context, cfg *configv1.Service, options ...service.GRPCOption) (*service.GRPCClient, error) {
	return grpc.NewClient(c, cfg, options...)
}

func (b *builder) NewHTTPClient(c context.Context, cfg *configv1.Service, options ...service.HTTPOption) (*service.HTTPClient, error) {
	return http.NewClient(c, cfg, options...)
}

func (b *builder) RegisterServiceBuilder(name string, factory service.ServerFactory) {
	b.ServiceBuilder.Register(name, factory)
}

// init initializes the builder struct.
func (b *builder) init() Builder {
	b.ConfigBuilder = config.DefaultBuilder
	b.RegistryBuilder = registry.DefaultBuilder
	b.ServiceBuilder = service.DefaultBuilder
	b.MiddlewareBuilder = middleware.DefaultBuilder
	return b
}

// NewConfig creates a new SelectorServer using the registered ConfigBuilder.
func NewConfig(cfg *configv1.SourceConfig, ss ...config.Option) (config.KConfig, error) {
	return runtimeBuilder.Config().NewConfig(cfg, ss...)
}

// RegisterConfig registers a ConfigBuilder with the builder.
func RegisterConfig(name string, factory config.Factory) {
	runtimeBuilder.Config().Register(name, factory)
}

// RegisterConfigFunc registers a ConfigBuilder with the builder.
func RegisterConfigFunc(name string, buildFunc config.BuildFunc) {
	runtimeBuilder.Config().Register(name, buildFunc)
}

// SyncConfig synchronizes the given configuration with the given value.
func SyncConfig(cfg *configv1.SourceConfig, v any, ss ...config.Option) error {
	return runtimeBuilder.SyncConfig(cfg, v, ss...)
}

func RegisterConfigSync(name string, syncFunc config.Syncer) {
	runtimeBuilder.RegisterConfigSync(name, syncFunc)
}

// NewDiscovery creates a new discovery using the registered RegistryBuilder.
func NewDiscovery(cfg *configv1.Discovery, ss ...registry.Option) (registry.KDiscovery, error) {
	return runtimeBuilder.NewDiscovery(cfg, ss...)
}

// NewRegistrar creates a new KRegistrar using the registered RegistryBuilder.
func NewRegistrar(cfg *configv1.Discovery, ss ...registry.Option) (registry.KRegistrar, error) {
	return runtimeBuilder.NewRegistrar(cfg, ss...)
}

// RegisterRegistry registers a RegistryBuilder with the builder.
func RegisterRegistry(name string, factory registry.Factory) {
	runtimeBuilder.RegisterRegistryBuilder(name, factory)
}

// NewMiddlewareClient creates a new KMiddleware with the builder.
func NewMiddlewareClient(name string, cm *middlewarev1.Middleware, ss ...middleware.Option) (middleware.KMiddleware, error) {
	return runtimeBuilder.NewMiddlewareClient(name, cm, ss...)
}

// NewMiddlewareServer creates a new KMiddleware with the builder.
func NewMiddlewareServer(name string, cm *middlewarev1.Middleware, ss ...middleware.Option) (middleware.KMiddleware, error) {
	return runtimeBuilder.NewMiddlewareServer(name, cm, ss...)
}

// NewMiddlewaresClient creates a new KMiddleware with the builder.
func NewMiddlewaresClient(cc *middlewarev1.Middleware, ss ...middleware.Option) []middleware.KMiddleware {
	return runtimeBuilder.NewMiddlewaresClient(cc, ss...)
}

// NewMiddlewaresServer creates a new KMiddleware with the builder.
func NewMiddlewaresServer(cc *middlewarev1.Middleware, ss ...middleware.Option) []middleware.KMiddleware {
	return runtimeBuilder.NewMiddlewaresServer(cc, ss...)
}

// RegisterMiddleware registers a MiddlewareBuilder with the builder.
func RegisterMiddleware(name string, builder middleware.Factory) {
	runtimeBuilder.RegisterMiddlewareBuilder(name, builder)
}

// NewHTTPServiceServer creates a new HTTP server using the provided configuration
func NewHTTPServiceServer(cfg *configv1.Service, ss ...service.HTTPOption) (*service.HTTPServer, error) {
	// Call the runtimeBuilder.NewHTTPServer function with the provided configuration
	return runtimeBuilder.NewHTTPServer(cfg, ss...)
}

// NewHTTPServiceClient creates a new HTTP client using the provided context and configuration
func NewHTTPServiceClient(ctx context.Context, cfg *configv1.Service, ss ...service.HTTPOption) (*service.HTTPClient, error) {
	// Call the runtimeBuilder.NewHTTPClient function with the provided context and configuration
	return runtimeBuilder.NewHTTPClient(ctx, cfg, ss...)
}

// NewGRPCServiceServer creates a new GRPC server using the provided configuration
func NewGRPCServiceServer(cfg *configv1.Service, ss ...service.GRPCOption) (*service.GRPCServer, error) {
	// Call the runtimeBuilder.NewGRPCServer function with the provided configuration
	return runtimeBuilder.NewGRPCServer(cfg, ss...)
}

// NewGRPCServiceClient creates a new GRPC client using the provided context and configuration
func NewGRPCServiceClient(ctx context.Context, cfg *configv1.Service, ss ...service.GRPCOption) (*service.GRPCClient, error) {
	// Call the runtimeBuilder.NewGRPCClient function with the provided context and configuration
	return runtimeBuilder.NewGRPCClient(ctx, cfg, ss...)
}

// RegisterService registers a service builder with the provided name
func RegisterService(name string, factory service.ServerFactory) {
	// Call the runtimeBuilder.RegisterServiceBuilder function with the provided name and service builder
	runtimeBuilder.Service().Register(name, factory)
}

// NewBuilder creates a new Builder.
func NewBuilder() Builder {
	b := &builder{
		syncs: make(map[string]config.Syncer),
	}
	return b.init()
}
