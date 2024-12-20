/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type (

	// Builder is an interface that defines a method for registering a RegistryBuilder.
	Builder interface {
		Factory
		RegisterRegistryBuilder(string, Factory)
	}

	// Factory is an interface that defines methods for creating a discovery and a Registrar.
	Factory interface {
		NewRegistrar(*configv1.Registry, ...OptionSetting) (Registrar, error)
		NewDiscovery(*configv1.Registry, ...OptionSetting) (Discovery, error)
	}
	Registry interface {
		Registrar
		Discovery
	}
)

// RegistrarBuildFunc is a function type that takes a *config.RegistryConfig and returns a Registrar and an error.
type RegistrarBuildFunc func(*configv1.Registry, ...OptionSetting) (Registrar, error)

// NewRegistrar is a method that calls the RegistrarBuildFunc with the given config.
func (fn RegistrarBuildFunc) NewRegistrar(cfg *configv1.Registry, ss ...OptionSetting) (Registrar, error) {
	return fn(cfg, ss...)
}

// DiscoveryBuildFunc is a function type that takes a *config.RegistryConfig and returns a discovery and an error.
type DiscoveryBuildFunc func(*configv1.Registry, ...OptionSetting) (Discovery, error)

// NewDiscovery is a method that calls the DiscoveryBuildFunc with the given config.
func (fn DiscoveryBuildFunc) NewDiscovery(cfg *configv1.Registry, ss ...OptionSetting) (Discovery, error) {
	return fn(cfg, ss...)
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