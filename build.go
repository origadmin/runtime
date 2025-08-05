/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/interfaces"
)

// NewConfig creates a new SelectorServer using the registered ConfigBuilder.
func NewConfig(cfg *configv1.SourceConfig, ss ...interfaces.Option) (kratosconfig.Config, error) {
	return defaultManager.ConfigBuilder.NewConfig(cfg, ss...)
}

// RegisterConfig registers a ConfigBuilder with the builder.
func RegisterConfig(name string, factory interfaces.ConfigFactory) {
	defaultManager.ConfigBuilder.Register(name, factory)
}

// RegisterConfigFunc registers a ConfigBuilder with the builder.
func RegisterConfigFunc(name string, buildFunc config.BuildFunc) {
	defaultManager.ConfigBuilder.Register(name, buildFunc)
}

// SyncConfig synchronizes the given configuration with the given value.
func SyncConfig(cfg *configv1.SourceConfig, v any, ss ...interfaces.Option) error {
	return defaultManager.ConfigBuilder.SyncConfig(cfg, v, ss...)
}

//// NewDiscovery creates a new discovery using the registered RegistryBuilder.
//func NewDiscovery(cfg *configv1.Discovery, ss ...interface{}) (registry.KDiscovery, error) {
//	return defaultManager.RegistryBuilder.NewDiscovery(cfg, ss...)
//}
//
//// NewRegistrar creates a new KRegistrar using the registered RegistryBuilder.
//func NewRegistrar(cfg *configv1.Discovery, ss ...interface{}) (registry.KRegistrar, error) {
//	return defaultManager.RegistryBuilder.NewRegistrar(cfg, ss...)
//}
//
//// RegisterRegistry registers a RegistryBuilder with the builder.
//func RegisterRegistry(name string, factory registry.Factory) {
//	defaultManager.RegistryBuilder.Register(name, factory)
//}
//
//// NewMiddlewareClient creates a new KMiddleware with the builder.
//func NewMiddlewareClient(name string, cm *middlewarev1.Middleware, ss ...middleware.Option) (middleware.KMiddleware, error) {
//	return defaultManager.MiddlewareProvider.BuildClient(cm, ss...)
//}
//
//// NewMiddlewareServer creates a new KMiddleware with the builder.
//func NewMiddlewareServer(name string, cm *middlewarev1.Middleware, ss ...middleware.Option) (middleware.KMiddleware, error) {
//	return defaultManager.MiddlewareProvider.BuildServer(cm, ss...)
//}
//
//// NewMiddlewaresClient creates a new KMiddleware with the builder.
//func NewMiddlewaresClient(cc *middlewarev1.Middleware, ss ...middleware.Option) []middleware.KMiddleware {
//	return defaultManager.MiddlewareProvider.BuildClient(cc, ss...)
//}
//
//// NewMiddlewaresServer creates a new KMiddleware with the builder.
//func NewMiddlewaresServer(cc *middlewarev1.Middleware, ss ...middleware.Option) []middleware.KMiddleware {
//	return defaultManager.MiddlewareProvider.BuildServer(cc, ss...)
//}
//
//// RegisterMiddleware registers a MiddlewareBuilder with the builder.
//func RegisterMiddleware(name string, builder middleware.Factory) {
//	defaultManager.MiddlewareProvider.Register(name, builder)
//}
//
//// NewHTTPServiceServer creates a new HTTP server using the provided configuration
//func NewHTTPServiceServer(cfg *configv1.Service, ss ...service.HTTPOption) (*service.HTTPServer, error) {
//	// Call the runtimeBuilder.NewHTTPServer function with the provided configuration
//	return defaultManager.ServerBuilder.NewHTTPServer(cfg, ss...)
//}
//
//// NewHTTPServiceClient creates a new HTTP client using the provided context and configuration
//func NewHTTPServiceClient(ctx context.Context, cfg *configv1.Service, ss ...service.HTTPOption) (*service.HTTPClient, error) {
//	// Call the runtimeBuilder.NewHTTPClient function with the provided context and configuration
//	return defaultManager.ServerBuilder.NewHTTPClient(ctx, cfg, ss...)
//}
//
//// NewGRPCServiceServer creates a new GRPC server using the provided configuration
//func NewGRPCServiceServer(cfg *configv1.Service, ss ...service.GRPCOption) (*service.GRPCServer, error) {
//	// Call the runtimeBuilder.NewGRPCServer function with the provided configuration
//	return defaultManager.ServerBuilder.NewGRPCServer(cfg, ss...)
//}
//
//// NewGRPCServiceClient creates a new GRPC client using the provided context and configuration
//func NewGRPCServiceClient(ctx context.Context, cfg *configv1.Service, ss ...service.GRPCOption) (*service.GRPCClient, error) {
//	// Call the runtimeBuilder.NewGRPCClient function with the provided context and configuration
//	return defaultManager.ServerBuilder.NewGRPCClient(ctx, cfg, ss...)
//}
//
//// RegisterService registers a service builder with the provided name
//func RegisterService(name string, factory service.ServerFactory) {
//	// Call the runtimeBuilder.RegisterServiceBuilder function with the provided name and service builder
//	defaultManager.ServerBuilder.Register(name, factory)
//}
