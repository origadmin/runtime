package container

import (
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
)

// DefaultComponentPriority defines the default priority for components.
// Components with lower priority values are initialized first.
// It is recommended to use values in increments of 100 to leave space for future adjustments.
const DefaultComponentPriority = 1000

// ComponentFunc is an adapter to allow the use of ordinary functions as ComponentFactory.
type ComponentFunc func(cfg interfaces.StructuredConfig, container Container, opts ...options.Option) (interfaces.Component, error)

// NewComponent calls the wrapped function.
func (c ComponentFunc) NewComponent(cfg interfaces.StructuredConfig, container Container, opts ...options.Option) (interfaces.Component, error) {
	return c(cfg, container, opts...)
}

// Priority returns the default priority for a ComponentFunc.
// This ensures that function-based components have a predictable, lower priority.
func (c ComponentFunc) Priority() int {
	return DefaultComponentPriority
}

// ComponentFactory defines the interface for creating generic components.
// It includes a priority system to manage initialization order.
type ComponentFactory interface {
	// Priority determines the initialization order of the component.
	// Components with lower priority values are created and initialized first.
	Priority() int
	// NewComponent creates a new component instance.
	// It receives a component-specific configuration and the container instance,
	// allowing it to register other components or access other services.
	NewComponent(cfg interfaces.StructuredConfig, container Container, opts ...options.Option) (interfaces.Component, error)
}
