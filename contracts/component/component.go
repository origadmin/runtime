package component

import (
	"context"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/engine/metadata"
)

// --- Configuration Sniffing Contracts ---

type (
	// AppConfigGetter defines the contract for retrieving application metadata.
	AppConfigGetter interface{ GetApp() *appv1.App }
	// LoggerConfigGetter defines the contract for retrieving logger configuration.
	LoggerConfigGetter interface{ GetLogger() *loggerv1.Logger }
	// MiddlewareConfigGetter defines the contract for retrieving middleware stacks.
	MiddlewareConfigGetter interface {
		GetMiddlewares() *middlewarev1.Middlewares
	}
	// DataConfigGetter defines the contract for retrieving data/storage configuration.
	DataConfigGetter interface{ GetData() *datav1.Data }
	// RegistryConfigGetter defines the contract for retrieving service discovery configuration.
	RegistryConfigGetter interface {
		GetDiscoveries() *discoveryv1.Discoveries
	}
)

// --- Engine Core Contracts ---

// Handle defines the interface for component retrieval and configuration resolution.
type Handle interface {
	// Get retrieves a component in the current Category context. "" means Default.
	Get(ctx context.Context, name string) (any, error)
	// In switches navigation to a different category/scope.
	In(category metadata.Category, opts ...RegisterOption) Handle
	// BindConfig performs high-performance type-safe assignment of current config.
	BindConfig(target any) error
	// Config returns the raw business configuration object.
	Config() any
	// Scope returns current isolation level.
	Scope() metadata.Scope
	// Category returns current functional category.
	Category() metadata.Category
}

// Provider is a function that creates a component instance.
type Provider func(ctx context.Context, h Handle, opts ...options.Option) (any, error)

// ConfigEntry represents a named configuration item.
type ConfigEntry struct {
	Name  string
	Value any
}

// ModuleConfig is the standardized output of an Extractor.
type ModuleConfig struct {
	Entries []ConfigEntry
	Active  string
}

// Extractor is a function that extracts a normalized configuration block from the root config.
type Extractor func(root any) (*ModuleConfig, error)

// Registry defines the interface for component management.
type Registry interface {
	Handle
	// Register registers a component factory.
	Register(c metadata.Category, e Extractor, p Provider, opts ...RegisterOption)
	// DefaultRegister registers a component factory only if none is already registered.
	DefaultRegister(c metadata.Category, e Extractor, p Provider, opts ...RegisterOption)
	// BindRoot injects the global business configuration object.
	BindRoot(root any)
	// Init performs sorted, sequential initialization of all registered components.
	Init(ctx context.Context) error
}

// RegisterOption is a functional option for component registration.
type RegisterOption func(any)
