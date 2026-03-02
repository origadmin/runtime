package engine

import (
	"github.com/origadmin/runtime/engine/container"
	"github.com/origadmin/runtime/engine/context"
	"github.com/origadmin/runtime/engine/protocol"
)

type (
	Category = context.Category
	Scope    = context.Scope
	Handle   = container.Handle
	Provider = container.Provider
	Registry = container.Registry

	Extractor      = protocol.Extractor
	ModuleConfig   = protocol.ModuleConfig
	ConfigEntry    = protocol.ConfigEntry
	RegisterOption = container.RegisterOption
)

const (
	GlobalScope = context.GlobalScope
	ServerScope = context.ServerScope
	ClientScope = context.ClientScope

	CategoryInfrastructure = context.CategoryInfrastructure
	CategoryRegistry       = context.CategoryRegistry
	CategoryClient         = context.CategoryClient
	CategoryServer         = context.CategoryServer
	CategoryMiddleware     = context.CategoryMiddleware
	CategoryDatabase       = context.CategoryDatabase
	CategoryCache          = context.CategoryCache
	CategoryObjectStore    = context.CategoryObjectStore
)

// Standard Priorities
const (
	PriorityInfrastructure = context.PriorityInfrastructure
	PriorityRegistry       = context.PriorityRegistry
	PriorityStorage        = context.PriorityStorage
	PriorityClientStack    = context.PriorityClientStack
	PriorityServerStack    = context.PriorityServerStack
)

// NewContainer creates a new engine container with the root business config.
func NewContainer(root any) Registry {
	return container.NewContainer(root)
}

// In is a helper to get a scoped handle from a registry.
// Focuses on Category as the primary dimension.
func In(h Handle, cat Category, opts ...RegisterOption) Handle {
	return h.In(cat, opts...)
}

// WithScope is a functional option to specify the scope during registration.
func WithScope(s Scope) RegisterOption {
	return container.WithScope(s)
}

// WithPriority is a functional option to specify the initialization priority.
func WithPriority(p int) RegisterOption {
	return container.WithPriority(p)
}
