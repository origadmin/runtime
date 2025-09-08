/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements a pluggable mechanism for service registration and discovery.
package registry

import (
	"context"
	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/framework/runtime/interfaces"
)

// --- Error Definitions ---
// Following the project's error handling specification.
// const (
// 	ReasonRegistryNotFound = "REGISTRY_NOT_FOUND"
// 	ReasonInvalidConfig    = "INVALID_CONFIG"
// 	ReasonCreationFailure  = "CREATION_FAILURE"
// )

// Builder defines the public interface for the registry builder.
// The concrete implementation is in factory.go.
type Builder interface {
	Register(name string, factory Factory)
	NewRegistrar(cfg *configv1.Discovery, opts *Options) (KRegistrar, error)
	NewDiscovery(cfg *configv1.Discovery, opts *Options) (KDiscovery, error)
}

// --- Top-Level API ---

// Register registers a new registry factory with the DefaultBuilder.
// It is a convenience wrapper around the builder's Register method.
func Register(name string, factory Factory) {
	DefaultBuilder().Register(name, factory)
}

// NewRegistrar creates a new KRegistrar instance using the DefaultBuilder.
func NewRegistrar(cfg *configv1.Discovery, opts ...Option) (KRegistrar, error) {
	registryOpts := &Options{ContextOptions: interfaces.ContextOptions{Context: context.Background()}}
	ApplyOptions(registryOpts, opts...)
	return DefaultBuilder().NewRegistrar(cfg, registryOpts)
}

// NewDiscovery creates a new KDiscovery instance using the DefaultBuilder.
func NewDiscovery(cfg *configv1.Discovery, opts ...Option) (KDiscovery, error) {
	registryOpts := &Options{ContextOptions: interfaces.ContextOptions{Context: context.Background()}}
	ApplyOptions(registryOpts, opts...)
	return DefaultBuilder().NewDiscovery(cfg, registryOpts)
}
