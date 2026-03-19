/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package component

import (
	"context"
	"errors"

	"github.com/origadmin/runtime/contracts/iterator"
)

var (
	// ErrRequirementNotFound is returned when a requested requirement is not found.
	ErrRequirementNotFound = errors.New("engine: requirement not found")
)

type (
	Category string
	Scope    string
	Priority int
)

const (
	// ReservedPrefix defines identifiers owned by the system.
	ReservedPrefix = "_"
	// DefaultName defines the name of the default component instance.
	DefaultName = "_default"
)

type Locator interface {
	Get(ctx context.Context, name ...string) (any, error)
	Iter(ctx context.Context) iterator.Iterator
	In(cat Category, opts ...InOption) Registry
	WithInScope(s Scope) Locator
	WithInTags(tags ...string) Locator
	Skip(names ...string) Locator
	Category() Category
	Scope() Scope
	Scopes() []Scope
	Tags() []string // Returns the "Package" of identities carried by this locator
}

type Handle interface {
	Config() any
	Name() string
	Category() Category
	Scope() Scope
	Locator() Locator
	Tag() string // Returns the SINGLE IDENTITY currently being processed by the provider
	Require(purpose string) (any, error)
}

type Provider func(ctx context.Context, h Handle) (any, error)

const (
	// PriorityFramework is the lowest priority, used for framework-level defaults.
	PriorityFramework Priority = 0
	// PriorityInfrastructure is used for common infrastructure (DB, Cache, etc).
	PriorityInfrastructure Priority = 100
	// PriorityDefault is the standard priority for business components.
	PriorityDefault Priority = 200
	// PriorityImportant is used for components that should override defaults.
	PriorityImportant Priority = 300
	// PriorityCritical is the highest priority, for emergency overrides.
	PriorityCritical Priority = 400
)

type Registry interface {
	Locator
	Register(p Provider, opts ...RegisterOption)
	Inject(name string, inst any, opts ...RegisterOption)
	IsRegistered(opts ...RegisterOption) bool
	Requirement(purpose string, resolver RequirementResolver)
}

type Container interface {
	Register(cat Category, p Provider, opts ...RegisterOption)
	Inject(cat Category, name string, inst any, opts ...RegisterOption)
	IsRegistered(cat Category, opts ...RegisterOption) bool
	Requirement(cat Category, purpose string, res RequirementResolver)
	Load(ctx context.Context, source any, opts ...LoadOption) error
	In(cat Category, opts ...InOption) Registry
}

// ConfigResolver resolves raw configuration source into ModuleConfig.
type ConfigResolver func(ctx context.Context, source any, opts *LoadOptions) (*ModuleConfig, error)

type Registration struct {
	Category Category
	Provider Provider
	Options  []RegisterOption
}

// ModuleConfig represents the configuration for an entire category.
// RequirementResolver at this level acts as a default for all entries in the category.
type ModuleConfig struct {
	Entries             []ConfigEntry
	Active              string
	RequirementResolver RequirementResolver
}

// ConfigEntry represents the configuration for a single component instance.
// RequirementResolver at this level takes precedence over the ModuleConfig level.
type ConfigEntry struct {
	Name                string
	Value               any
	RequirementResolver RequirementResolver
}
type RequirementResolver func(ctx context.Context, h Handle, purpose string) (any, error)

type RegistrationOptions struct {
	Scopes              []Scope
	Tag                 string
	Priority            Priority
	DefaultEntries      []string
	ConfigResolver      ConfigResolver
	RequirementResolver RequirementResolver
}

type RegisterOption func(*RegistrationOptions)

// InOption is a functional option that modifies a Registry.
type InOption func(Registry) Registry

type LoadOptions struct {
	Category Category
	Scope    Scope
	Name     string
	Tags     []string
	Resolver ConfigResolver
}

type LoadOption func(*LoadOptions)
