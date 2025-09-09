/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements a pluggable mechanism for service registration and discovery.
package registry

import (
	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces/factory"
)

// --- Error Definitions ---
// Following the project's error handling specification.
// const (
// 	ReasonRegistryNotFound = "REGISTRY_NOT_FOUND"
// 	ReasonInvalidConfig    = "INVALID_CONFIG"
// 	ReasonCreationFailure  = "CREATION_FAILURE"
// )

// Registry defines the interface for service registration and discovery.
type Registry interface {
	KRegistrar
	KDiscovery
}

type Builder interface {
	factory.Registry[Factory]
	Factory
}

// --- Top-Level API ---

// Register registers a new registry factory with the DefaultBuilder.
// It is a convenience wrapper around the builder's Register method.
func Register(name string, factory Factory) {
	DefaultBuilder().Register(name, factory)
}

// NewRegistrar creates a new KRegistrar instance using the DefaultBuilder.
func NewRegistrar(cfg *configv1.Discovery, opts ...Option) (KRegistrar, error) {
	return DefaultBuilder().NewRegistrar(cfg, opts...)
}

// NewDiscovery creates a new KDiscovery instance using the DefaultBuilder.
func NewDiscovery(cfg *configv1.Discovery, opts ...Option) (KDiscovery, error) {
	return DefaultBuilder().NewDiscovery(cfg, opts...)
}
