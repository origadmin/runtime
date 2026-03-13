package engine

import (
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/engine/container"
)

type (
	Category = component.Category
	Scope    = component.Scope
	Priority = component.Priority

	Handle    = component.Handle
	Provider  = component.Provider
	Registry  = component.Registry
	Container = component.Container

	RequirementResolver = component.RequirementResolver
	ConfigResolver      = component.ConfigResolver
	Registration        = component.Registration
	ModuleConfig        = component.ModuleConfig
	ConfigEntry         = component.ConfigEntry
	RegistrationOptions = component.RegistrationOptions

	Option = container.Option

	InOption    = component.InOption
	Locator     = component.Locator
	LoadOption  = component.LoadOption
	LoadOptions = component.LoadOptions

	RegisterOption = component.RegisterOption
)

const (
	ReservedPrefix = component.ReservedPrefix

	PriorityFramework      = component.PriorityFramework
	PriorityInfrastructure = component.PriorityInfrastructure
	PriorityDefault        = component.PriorityDefault
	PriorityImportant      = component.PriorityImportant
	PriorityCritical       = component.PriorityCritical
)

var globalRegistrations []Registration

// Register adds a component registration to the global pool.
func Register(cat Category, p Provider, opts ...RegisterOption) {
	globalRegistrations = append(globalRegistrations, Registration{
		Category: cat,
		Provider: p,
		Options:  opts,
	})
}

// GlobalRegistrations returns a copy of all registrations in the global pool.
func GlobalRegistrations() []Registration {
	res := make([]Registration, len(globalRegistrations))
	copy(res, globalRegistrations)
	return res
}

// --- Registration Options ---

// WithConfigResolverOption specifies a local configuration resolver for a component.
func WithConfigResolverOption(res ConfigResolver) RegisterOption {
	return func(o *RegistrationOptions) {
		o.ConfigResolver = res
	}
}

// WithRequirementResolverOption specifies a local requirement resolver for a component.
func WithRequirementResolverOption(res RequirementResolver) RegisterOption {
	return func(o *RegistrationOptions) {
		o.RequirementResolver = res
	}
}

// Deprecated: Use WithConfigResolverOption instead.
var WithResolverOption = WithConfigResolverOption

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

func WithTag(tag string) RegisterOption {
	return func(o *RegistrationOptions) {
		o.Tag = tag
	}
}

func WithDefaultEntries(names ...string) RegisterOption {
	return func(o *RegistrationOptions) {
		o.DefaultEntries = append(o.DefaultEntries, names...)
	}
}

// Deprecated: Use WithDefaultEntries instead.
var WithDefaultEntry = WithDefaultEntries

// --- Perspective Options (USING INTERFACE METHODS) ---

// WithInScope specifies the perspective scope.
func WithInScope(s Scope) InOption {
	return func(l Registry) Registry {
		return l.WithInScope(s).(Registry)
	}
}

// WithInTags specifies the perspective tags.
func WithInTags(tags ...string) InOption {
	return func(l Registry) Registry {
		return l.WithInTags(tags...).(Registry)
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

func WithLoadResolver(res ConfigResolver) LoadOption {
	return func(o *LoadOptions) {
		o.Resolver = res
	}
}

// --- Container Bootstrapping ---

type RegistryOptions struct {
	CategoryResolvers map[Category]ConfigResolver
	Registrations     []Registration
}

type RegistryOption func(*RegistryOptions)

// WithCategoryResolvers sets the default resolvers for specific categories.
func WithCategoryResolvers(res map[Category]ConfigResolver) RegistryOption {
	return func(o *RegistryOptions) {
		if o.CategoryResolvers == nil {
			o.CategoryResolvers = make(map[Category]ConfigResolver)
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
func NewContainer(opts ...RegistryOption) Container {
	o := &RegistryOptions{
		CategoryResolvers: make(map[Category]ConfigResolver),
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
