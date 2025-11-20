package interfaces

import (
	"github.com/origadmin/runtime/interfaces/options"
)

// Component is a generic runtime component.
type Component interface{} // Minimal definition

// ComponentFactory defines the interface for creating generic components.
type ComponentFactory interface {
	NewComponent(cfg *StructuredConfig, container Container, opts ...options.Option) (Component, error)
}
