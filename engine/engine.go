package engine

import (
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/engine/container"
	"github.com/origadmin/runtime/engine/metadata"
)

type (
	Category = metadata.Category
	Scope    = metadata.Scope
	Handle   = component.Handle
	Provider = component.Provider
	Registry = component.Registry

	Extractor           = component.Extractor
	ModuleConfig        = component.ModuleConfig
	ConfigEntry         = component.ConfigEntry
	RegistrationOptions = component.RegistrationOptions
	RegisterOption      = component.RegisterOption
)

const (
	GlobalScope = metadata.GlobalScope
	ServerScope = metadata.ServerScope
	ClientScope = metadata.ClientScope

	CategoryInfrastructure = metadata.CategoryInfrastructure
	CategoryLogger         = metadata.CategoryLogger
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

// NewContainer creates a new engine container.
func NewContainer() Registry {
	return container.NewContainer()
}

// In is a helper to get a scoped handle from a registry.
func In(h Handle, cat Category, opts ...RegisterOption) Handle {
	return h.In(cat, opts...)
}

// WithScope is a functional option to specify the scope during registration.
func WithScope(s Scope) RegisterOption {
	return func(o *RegistrationOptions) {
		o.Scope = s
	}
}

// WithPriority is a functional option to specify the initialization priority.
func WithPriority(p int) RegisterOption {
	return func(o *RegistrationOptions) {
		o.Priority = p
	}
}
