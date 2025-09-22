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

	// GetLogger returns the initialized logger instance.
	GetLogger() kratoslog.Logger // Changed to Kratos log.Logger

	// GetDiscoveries returns a map of initialized Kratos Discovery clients.
	GetDiscoveries() map[string]registry.Discovery

	// GetRegistrars returns a map of initialized Kratos Registrar clients.
	GetRegistrars() map[string]registry.Registrar

	// GetDefaultRegistrar returns the default Kratos Registrar for self-registration.
	GetDefaultRegistrar() registry.Registrar

	// --- Generic Service Locator for Extensibility ---

	// GetComponent retrieves a component by its registered name.
	// This allows for future components to be added without changing the interface.
	GetComponent(name string) (component interface{}, ok bool)
}
