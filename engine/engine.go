package engine

import (
	"github.com/origadmin/runtime/engine/container"
	"github.com/origadmin/runtime/engine/metadata"
	"github.com/origadmin/runtime/engine/protocol"
)

type (
	Category = metadata.Category
	Scope    = metadata.Scope
	Handle   = container.Handle
	Provider = container.Provider
	Registry = container.Registry

	Extractor      = protocol.Extractor
	ModuleConfig   = protocol.ModuleConfig
	ConfigEntry    = protocol.ConfigEntry
	RegisterOption = container.RegisterOption
)

const (
	GlobalScope = metadata.GlobalScope
	ServerScope = metadata.ServerScope
	ClientScope = metadata.ClientScope

	CategoryInfrastructure = metadata.CategoryInfrastructure
	CategoryRegistry       = metadata.CategoryRegistry
	CategoryClient         = metadata.CategoryClient
	CategoryServer         = metadata.CategoryServer
	CategoryMiddleware     = metadata.CategoryMiddleware
	CategoryDatabase       = metadata.CategoryDatabase
	CategoryCache          = metadata.CategoryCache
	CategoryObjectStore    = metadata.CategoryObjectStore
	CategoryQueue          = metadata.CategoryQueue
	CategoryTask           = metadata.CategoryTask
	CategoryMail           = metadata.CategoryMail
	CategoryStorage        = metadata.CategoryStorage
)

// Standard Priorities
const (
	PriorityInfrastructure = metadata.PriorityInfrastructure
	PriorityRegistry       = metadata.PriorityRegistry
	PriorityStorage        = metadata.PriorityStorage
	PriorityClientStack    = metadata.PriorityClientStack
	PriorityServerStack    = metadata.PriorityServerStack
)

// NewContainer creates a new engine container with the root business config.
func NewContainer(root any) Registry {
	return container.NewContainer(root)
}

// In is a helper to get a scoped handle from a registry.
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
