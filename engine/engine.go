package engine

import (
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/engine/container"
)

type (
	Category = component.Category
	Scope    = component.Scope
	Priority = component.Priority
	Handle   = component.Handle
	Provider = component.Provider
	Registry = component.Registry

	Extractor           = component.Extractor
	ModuleConfig        = component.ModuleConfig
	ConfigEntry         = component.ConfigEntry
	RegistrationOptions = component.RegistrationOptions
	RegisterOption      = component.RegisterOption
	InOptions           = component.InOptions
	InOption            = component.InOption
	LoadOptions         = component.LoadOptions
	LoadOption          = component.LoadOption
)

const (
	GlobalScope = component.GlobalScope
)

// NewContainer creates a new engine container.
func NewContainer() Registry {
	return container.NewContainer()
}

// --- Registration Options ---

// WithScope specifies visibility for a single scope during registration.
func WithScope(s Scope) RegisterOption {
	return func(o *RegistrationOptions) {
		o.Scopes = append(o.Scopes, s)
	}
}

// WithScopes specifies visibility for multiple scopes during registration.
func WithScopes(ss ...Scope) RegisterOption {
	return func(o *RegistrationOptions) {
		o.Scopes = append(o.Scopes, ss...)
	}
}

// WithPriority specifies initialization priority.
func WithPriority(p Priority) RegisterOption {
	return func(o *RegistrationOptions) {
		o.Priority = p
	}
}

// WithExtractor specifies a local config extractor.
func WithExtractor(e Extractor) RegisterOption {
	return func(o *RegistrationOptions) {
		o.Extractor = e
	}
}

// --- Perspective Options (In) ---

// WithInScope specifies the exact perspective for a handle.
func WithInScope(s Scope) InOption {
	return func(o *InOptions) {
		o.Scope = s
	}
}

// --- Load Options ---

func ForCategory(cat Category) LoadOption {
	return func(o *LoadOptions) {
		o.Category = cat
	}
}

func ForName(name string) LoadOption {
	return func(o *LoadOptions) {
		o.Name = name
	}
}
