/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

import (
	"errors"

	"github.com/go-kratos/kratos/v2/registry"
)

// This is only alias type for wrapped
type (
	KWatcher         = registry.Watcher
	KServiceInstance = registry.ServiceInstance
	KDiscovery       = registry.Discovery
	KRegistrar       = registry.Registrar
)

var (
	ErrRegistryNotFound = errors.New("registry not found")
)
