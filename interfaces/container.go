package interfaces

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/origadmin/runtime/middleware"
)

// ComponentFactory defines the signature for a function that can create a generic component.
// It receives the global configuration and the specific configuration map for the component instance.
type ComponentFactory func(cfg StructuredConfig, container Container) (interface{}, error)

// Container defines the interface for retrieving fully-initialized application components.
// It is the return type of bootstrap.NewProvider and the input for runtime.New.
type Container interface {
	// Logger returns the configured Kratos logger.
	Logger() log.Logger

	// Discoveries returns a map of all configured service discovery components.
	Discoveries() map[string]registry.Discovery

	// Registrars returns a map of all configured service registrar components.
	Registrars() map[string]registry.Registrar

	// DefaultRegistrar returns the default service registrar, used for service self-registration.
	// It may be nil if no default registry is configured.
	DefaultRegistrar() registry.Registrar

	ServerMiddleware(name middleware.Name) (middleware.KMiddleware, bool)
	ClientMiddleware(name middleware.Name) (middleware.KMiddleware, bool)

	// --- Generic Service Locator for Extensibility ---

	// Component retrieves a generic component by its registered name.
	// This allows for future components to be added without changing the interface.
	Component(name string) (component interface{}, ok bool)
}
