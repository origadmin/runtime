/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package factory

import (
	"context"

	"github.com/origadmin/runtime/contracts/options"
)

// Handle is a minimalist interface for context-bound instance retrieval.
// Implementation is provided by the internal IoC container.
type Handle interface {
	// Get retrieves an instance by category and name.
	Get(category string, name string) any
}

// ComponentFactory is the unified interface for all component creators.
// It decouples the implementation from the framework's core lifecycle.
type ComponentFactory interface {
	// New creates a new instance using the provided config and container handle.
	// cfg: The raw configuration object (usually a proto segment).
	// h:   The handle to retrieve other dependencies from the container.
	New(ctx context.Context, cfg any, h Handle, opts ...options.Option) (any, error)
}

// ConfigGetter is a function that extracts a specific configuration segment from the root config.
type ConfigGetter func(root any) any

// Binding represents the explicit link between a component identity and its creator.
type Binding struct {
	Category string
	Name     string
	Factory  ComponentFactory
	Getter   ConfigGetter
}
