package bootstrap

import (
	"fmt"
	"sync"

	"github.com/origadmin/runtime/interfaces"
)

// ComponentFactoryFunc defines the signature for a function that can create a generic component.
// It receives the global configuration and the specific configuration map for the component instance.
type ComponentFactoryFunc func(config interfaces.Config, componentConfig map[string]interface{}) (interface{}, error)

var (
	factories      = make(map[string]ComponentFactoryFunc)
	factoriesMutex sync.RWMutex
)

// RegisterComponentFactory registers a new component factory with the bootstrap system.
// This function is typically called from the init() function of a component's package.
// It is safe for concurrent use.
func RegisterComponentFactory(componentType string, factory ComponentFactoryFunc) {
	factoriesMutex.Lock()
	defer factoriesMutex.Unlock()

	if _, exists := factories[componentType]; exists {
		// In a real-world scenario, you might want to panic or log a fatal error here
		// to prevent accidental overwrites during development.
		panic(fmt.Sprintf("component factory for type '%s' is already registered", componentType))
	}

	factories[componentType] = factory
}

// getFactory retrieves a component factory by its type.
// It is safe for concurrent use.
func getFactory(componentType string) (ComponentFactoryFunc, bool) {
	factoriesMutex.RLock()
	defer factoriesMutex.RUnlock()

	factory, ok := factories[componentType]
	return factory, ok
}
