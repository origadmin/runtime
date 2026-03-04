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

	Resolver            = component.Resolver
	Registration        = component.Registration
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

// Global Pool for init-phase registrations
var globalPool []Registration

// Register stores a component registration in the global pool (typically used in init()).
func Register(cat Category, p Provider, opts ...RegisterOption) {
	globalPool = append(globalPool, Registration{
		Category: cat,
		Provider: p,
		Options:  opts,
	})
}

// GlobalRegistrations returns a snapshot of the current global pool.
func GlobalRegistrations() []Registration {
	res := make([]Registration, len(globalPool))
	copy(res, globalPool)
	return res
}

// --- Container Bootstrapping ---

type RegistryOptions struct {
	Resolver      Resolver
	Registrations []Registration
}

type RegistryOption func(*RegistryOptions)

// WithResolver sets the default global resolver for the container.
func WithResolver(res Resolver) RegistryOption {
	return func(o *RegistryOptions) {
		o.Resolver = res
	}
}

// WithGlobalRegistrations instructs the container to load all registrations from the global pool.
func WithGlobalRegistrations() RegistryOption {
	return func(o *RegistryOptions) {
		o.Registrations = append(o.Registrations, GlobalRegistrations()...)
	}
}

// NewContainer creates a new engine container based on provided options.
func NewContainer(opts ...RegistryOption) Registry {
	o := &RegistryOptions{}
	for _, opt := range opts {
		opt(o)
	}
	reg := container.NewContainer()
	if o.Resolver != nil {
		reg.SetResolver(o.Resolver)
	}
	for _, r := range o.Registrations {
		reg.Register(r.Category, r.Provider, r.Options...)
	}
	return reg
}

// --- Registration Options ---

func WithScope(s Scope) RegisterOption {
	return func(o *RegistrationOptions) {
		o.Scopes = append(o.Scopes, s)
	}
}

func WithScopes(ss ...Scope) RegisterOption {
	return func(o *RegistrationOptions) {
		o.Scopes = append(o.Scopes, ss...)
	}
}

func WithPriority(p Priority) RegisterOption {
	return func(o *RegistrationOptions) {
		o.Priority = p
	}
}

// WithResolverOption specifies a local configuration resolver for a component.
func WithResolverOption(res Resolver) RegisterOption {
	return func(o *RegistrationOptions) {
		o.Resolver = res
	}
}

// --- Perspective Options (In) ---

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

func WithLoadResolver(res Resolver) LoadOption {
	return func(o *LoadOptions) {
		o.Resolver = res
	}
}
