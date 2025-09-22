package bootstrap

import (
	kratoslog "github.com/go-kratos/kratos/v2/log" // Import Kratos log package
	"github.com/go-kratos/kratos/v2/registry"      // Import Kratos registry types
)

// componentProviderImpl is the concrete implementation of the interfaces.ComponentProvider.
// It holds all the initialized components created during the bootstrap process.
// This struct is private to the bootstrap package.
type componentProviderImpl struct {
	logger           kratoslog.Logger              // Changed to kratoslog.Logger
	discoveries      map[string]registry.Discovery // Changed to Kratos registry.Discovery
	registrars       map[string]registry.Registrar // Added for Kratos registrars
	defaultRegistrar registry.Registrar            // Added for default Kratos registrar

	// components holds components that are not strongly-typed in the interface,
	// allowing for dynamic extension without modifying the ComponentProvider interface.
	components map[string]interface{}
}

// Logger returns the initialized logger instance.
func (o *componentProviderImpl) Logger() kratoslog.Logger { // Changed to kratoslog.Logger
	return o.logger
}

// Discoveries returns a map of initialized Kratos Discovery clients.
func (o *componentProviderImpl) Discoveries() map[string]registry.Discovery {
	return o.discoveries
}

// Registrars returns a map of initialized Kratos Registrar clients.
func (o *componentProviderImpl) Registrars() map[string]registry.Registrar {
	return o.registrars
}

// DefaultRegistrar returns the default Kratos Registrar for self-registration.
func (o *componentProviderImpl) DefaultRegistrar() registry.Registrar {
	return o.defaultRegistrar
}

// Component retrieves a component by its registered name from the components map.
func (o *componentProviderImpl) Component(name string) (component interface{}, ok bool) {
	component, ok = o.components[name]
	return
}
