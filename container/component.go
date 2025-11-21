package container

import (
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
)

type ComponentFunc func(cfg *interfaces.StructuredConfig, container Container, opts ...options.Option) (interfaces.Component,
	error)

func (c ComponentFunc) NewComponent(cfg *interfaces.StructuredConfig, container Container, opts ...options.Option) (interfaces.Component, error) {
	return c(cfg, container, opts...)
}

// ComponentFactory defines the interface for creating generic components.
type ComponentFactory interface {
	NewComponent(cfg *interfaces.StructuredConfig, container Container, opts ...options.Option) (interfaces.Component, error)
}
