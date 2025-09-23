package interfaces

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	appv1 "github.com/origadmin/runtime/api/gen/go/app/v1"
)

// ComponentProvider defines the interface for retrieving fully-initialized application components.
// It is the return type of bootstrap.NewProvider and the input for runtime.New.
type ComponentProvider interface {
	// --- Strongly-Typed Accessors for Core Components ---

	// AppInfo returns the application's configured information (ID, name, version, metadata).
	AppInfo() *appv1.AppInfo

	// Logger returns the configured Kratos logger.
	Logger() log.Logger

	// Discoveries returns a map of all configured service discovery components.
	Discoveries() map[string]registry.Discovery

	// Registrars returns a map of all configured service registrar components.
	Registrars() map[string]registry.Registrar

	// DefaultRegistrar returns the default service registrar, used for service self-registration.
	// It may be nil if no default registry is configured.
	DefaultRegistrar() registry.Registrar

	// --- Generic Service Locator for Extensibility ---

	// Component retrieves a generic component by its registered name.
	// This allows for future components to be added without changing the interface.
	Component(name string) (component interface{}, ok bool)
}
