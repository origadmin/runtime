package constant

// Component keys used for path mapping and dependency injection.
const (
	// ConfigApp is the key for the application configuration.
	ConfigApp = "app"

	// ComponentLogger is the key for the logger component.
	ComponentLogger = "logger"

	// ComponentRegistries is the key for the service registries component.
	ComponentRegistries = "registries"

	// ComponentMiddlewares is the key for the middlewares component.
	ComponentMiddlewares = "middlewares"
)

// defaultComponentPaths provides the framework's default path map for core components.
// This is a private variable to prevent external mutation.
var defaultComponentPaths = map[string]string{
	ConfigApp:            "app",
	ComponentLogger:      "logger",
	ComponentRegistries:  "registries",
	ComponentMiddlewares: "middlewares",
}

// DefaultComponentPaths returns a copy of the default component path map.
// Returning a copy ensures that the original map cannot be modified by external packages,
// which is a critical safety measure.
func DefaultComponentPaths() map[string]string {
	// Create a new map with the same capacity.
	paths := make(map[string]string, len(defaultComponentPaths))
	// Copy key-value pairs.
	for k, v := range defaultComponentPaths {
		paths[k] = v
	}
	return paths
}
