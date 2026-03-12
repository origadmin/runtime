/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package component

import (
	"context"
	"iter"
)

type (
	Category string
	Scope    string
	Priority int
)

const (
	CategoryLogger      Category = "logger"
	CategoryRegistrar   Category = "registrar"
	CategoryDiscovery   Category = "discovery"
	CategoryClient      Category = "client"
	CategoryServer      Category = "server"
	CategoryMiddleware  Category = "middleware"
	CategoryDatabase    Category = "database"
	CategoryCache       Category = "cache"
	CategoryObjectStore Category = "objectstore"
	CategoryQueue       Category = "queue"
	CategoryTask        Category = "task"
	CategoryMail        Category = "mail"
	CategoryStorage     Category = "storage"
	CategorySecurity    Category = "security"
	CategorySkipper     Category = "skipper"
)

const (
	// GlobalScope is the default system fallback scope.
	GlobalScope Scope = "_global"
	// ServerScope is the standard scope for server-side components.
	ServerScope Scope = "server"
	// ClientScope is the standard scope for client-side components.
	ClientScope Scope = "client"

	// DefaultName is the system key for the active/default instance.
	DefaultName = "_default"
	// ReservedPrefix defines identifiers owned by the system.
	ReservedPrefix = "_"
)

type Locator interface {
	Get(ctx context.Context, name string) (any, error)
	Iter(ctx context.Context) iter.Seq2[string, any]
	In(cat Category, opts ...InOption) Locator
	WithInScope(s Scope) Locator
	WithInTags(tags ...string) Locator
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
	Register(cat Category, p Provider, opts ...RegisterOption)
	Inject(cat Category, name string, inst any, opts ...RegisterOption)
	Has(cat Category, opts ...RegisterOption) bool
	Load(ctx context.Context, source any, opts ...LoadOption) error
	In(cat Category, opts ...InOption) Locator
}

type Resolver func(source any, cat Category) (*ModuleConfig, error)

type Registration struct {
	Category Category
	Provider Provider
	Options  []RegisterOption
}

type ModuleConfig struct {
	Entries []ConfigEntry
	Active  string
}

type ConfigEntry struct {
	Name  string
	Value any
}

type RegistrationOptions struct {
	Resolver       Resolver
	Scopes         []Scope
	Priority       Priority
	Tag            string
	DefaultEntries []string
}

type RegisterOption func(*RegistrationOptions)

// InOption is a functional option that modifies a Locator.
// It directly uses the Locator interface to support fluent perspective switching.
type InOption func(Locator) Locator

type LoadOptions struct {
	Category Category
	Scope    Scope
	Name     string
	Resolver Resolver
	Tags     []string
}

type LoadOption func(*LoadOptions)

// Helper interfaces for identifying configuration entries
type (
	// Named represents an object that has a unique name.
	Named interface {
		GetName() string
	}

	// Typed represents an object that has a specific type or category.
	Typed interface {
		GetType() string
	}

	// Dialectal represents an object that specifies a database dialect.
	Dialectal interface {
		GetDialect() string
	}

	// Driver represents an object that specifies a underlying driver.
	Driver interface {
		GetDriver() string
	}
)
