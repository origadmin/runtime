/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

import (
	"github.com/go-kratos/kratos/v2/registry"
)

// This is only alias type for registry
type (
	Watcher         = registry.Watcher
	ServiceInstance = registry.ServiceInstance
	Discovery       = registry.Discovery
	Registrar       = registry.Registrar
)
