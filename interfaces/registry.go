/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package interfaces

import (
	registry "github.com/go-kratos/kratos/v2/registry"
)

type (
	// KDiscovery is an alias for github.com/go-kratos/kratos/v2/registry.Discovery
	KDiscovery = registry.Discovery
	// KRegistrar is an alias for github.com/go-kratos/kratos/v2/registry.Registrar
	KRegistrar = registry.Registrar

	// Registry defines the interface for service registration and discovery.
	Registry interface {
		KRegistrar
		KDiscovery
	}
)
