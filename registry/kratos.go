/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

import (
	"errors"

	"github.com/go-kratos/kratos/v2/registry"
)

//go:generate adptool .
//go:adapter:package github.com/go-kratos/kratos/v2/middleware
//go:adapter:package:type *
//go:adapter:package:type:prefix Kratos
//go:adapter:package:func *
//go:adapter:package:func:prefix Kratos

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
