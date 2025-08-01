/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

import (
	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces/factory"
)

type (

	// Builder is an interface that defines a method for registering a RegistryBuilder.
	Builder interface {
		factory.Registry[Factory]
		Factory
	}
	// Factory is an interface that defines methods for creating a discovery and a KRegistrar.
	Factory interface {
		NewRegistrar(*configv1.Discovery, ...Option) (KRegistrar, error)
		NewDiscovery(*configv1.Discovery, ...Option) (KDiscovery, error)
	}

	Registry interface {
		KRegistrar
		KDiscovery
	}
)

// RegistrarBuildFunc is a function type that takes a *config.RegistryConfig and returns a KRegistrar and an error.
type RegistrarBuildFunc func(*configv1.Discovery, ...Option) (KRegistrar, error)

// NewRegistrar is a method that calls the RegistrarBuildFunc with the given config.
func (fn RegistrarBuildFunc) NewRegistrar(cfg *configv1.Discovery, ss ...Option) (KRegistrar, error) {
	return fn(cfg, ss...)
}

// DiscoveryBuildFunc is a function type that takes a *config.RegistryConfig and returns a discovery and an error.
type DiscoveryBuildFunc func(*configv1.Discovery, ...Option) (KDiscovery, error)

// NewDiscovery is a method that calls the DiscoveryBuildFunc with the given config.
func (fn DiscoveryBuildFunc) NewDiscovery(cfg *configv1.Discovery, ss ...Option) (KDiscovery, error) {
	return fn(cfg, ss...)
}

type FuncFactory struct {
	RegistrarFunc func(*configv1.Discovery, ...Option) (KRegistrar, error)
	DiscoveryFunc func(*configv1.Discovery, ...Option) (KDiscovery, error)
}

func (f FuncFactory) NewRegistrar(cfg *configv1.Discovery, opts ...Option) (KRegistrar, error) {
	return f.RegistrarFunc(cfg, opts...)
}

func (f FuncFactory) NewDiscovery(cfg *configv1.Discovery, opts ...Option) (KDiscovery, error) {
	return f.DiscoveryFunc(cfg, opts...)
}

// wrapped is a struct that embeds RegistrarBuildFunc and DiscoveryBuildFunc.
type wrapped struct {
	RegistrarBuildFunc
	DiscoveryBuildFunc
}

func WrapFactory(registrar RegistrarBuildFunc, discovery DiscoveryBuildFunc) Factory {
	return &wrapped{
		RegistrarBuildFunc: registrar,
		DiscoveryBuildFunc: discovery,
	}
}

// _ is a blank identifier that is used to satisfy the interface requirement for RegistryBuilder.
var _ Factory = &wrapped{}
