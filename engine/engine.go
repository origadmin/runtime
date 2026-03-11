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
	CategoryResolvers map[Category]Resolver
	Registrations     []Registration
}

type RegistryOption func(*RegistryOptions)

// WithCategoryResolvers sets the default resolvers for specific categories.
func WithCategoryResolvers(res map[Category]Resolver) RegistryOption {
	return func(o *RegistryOptions) {
		if o.CategoryResolvers == nil {
			o.CategoryResolvers = make(map[Category]Resolver)
		}
		for k, v := range res {
			o.CategoryResolvers[k] = v
		}
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
	o := &RegistryOptions{
		CategoryResolvers: make(map[Category]Resolver),
	}
	for _, opt := range opts {
		opt(o)
	}

	// Transform public options to internal container options
	var internalOpts []container.Option
	if len(o.CategoryResolvers) > 0 {
		internalOpts = append(internalOpts, container.WithCategoryResolvers(o.CategoryResolvers))
	}

	reg := container.NewContainer(internalOpts...)
	for _, r := range o.Registrations {
		reg.Register(r.Category, r.Provider, r.Options...)
	}
	return reg
}

// --- Registration Options (STRICT SINGLE TAG) ---

// WithScope registers the component for a specific scope.
func WithScope(s Scope) RegisterOption {
	return func(o *RegistrationOptions) {
		o.Scopes = append(o.Scopes, s)
	}
}

// WithScopes registers the component for multiple scopes.
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

// WithTag registers the component with a SINGLE identity tag.
// Architecture Rule: A registered component has one unique identity.
func WithTag(tag string) RegisterOption {
	return func(o *RegistrationOptions) {
		o.Tag = tag
	}
}

// WithDefaultEntry ensures that a specific component name is always present in the container
// for this provider, even if it's missing from the external configuration source.
func WithDefaultEntry(name string) RegisterOption {
	return func(o *RegistrationOptions) {
		o.DefaultEntries = append(o.DefaultEntries, name)
	}
}

// WithResolverOption specifies a local configuration resolver for a component.
func WithResolverOption(res Resolver) RegisterOption {
	return func(o *RegistrationOptions) {
		o.Resolver = res
	}
}

// --- Perspective Options (USING INTERFACE METHODS) ---

// WithInScope specifies the perspective scope.
func WithInScope(s Scope) InOption {
	return func(l component.Locator) component.Locator {
		return l.WithInScope(s)
	}
}

// WithInTags specifies the perspective tags.
func WithInTags(tags ...string) InOption {
	return func(l component.Locator) component.Locator {
		return l.WithInTags(tags...)
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

func ForScope(s Scope) LoadOption {
	return func(o *LoadOptions) {
		o.Scope = s
	}
}

func WithLoadResolver(res Resolver) LoadOption {
	return func(o *LoadOptions) {
		o.Resolver = res
	}
}
