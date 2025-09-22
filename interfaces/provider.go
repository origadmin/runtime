package interfaces

import (
	kratoslog "github.com/go-kratos/kratos/v2/log" // Import Kratos log package
	"github.com/go-kratos/kratos/v2/registry"
)

// ComponentProvider defines the contract for accessing core application components.
// It supports both strongly-typed access for core components and a generic
// service locator pattern for future extensibility.
type ComponentProvider interface {
	// --- Strongly-Typed Accessors for Core Components ---

	// Logger returns the initialized logger instance.
	Logger() kratoslog.Logger // Changed to Kratos log.Logger

	// Discoveries returns a map of initialized Kratos Discovery clients.
	Discoveries() map[string]registry.Discovery

	// Registrars returns a map of initialized Kratos Registrar clients.
	Registrars() map[string]registry.Registrar

	// DefaultRegistrar returns the default Kratos Registrar for self-registration.
	DefaultRegistrar() registry.Registrar

	// --- Generic Service Locator for Extensibility ---

	// Component retrieves a component by its registered name.
	// This allows for future components to be added without changing the interface.
	Component(name string) (component interface{}, ok bool)
}
