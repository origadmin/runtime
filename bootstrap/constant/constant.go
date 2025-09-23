package constant

// Core component keys used for identifying components in the bootstrap paths map.
const (
	// ComponentLogger is the key for the logger configuration path.
	ComponentLogger = "logger"

	// ComponentRegistries is the key for the registries/discoveries configuration path.
	ComponentRegistries = "registries"

	// ComponentComponents is the key for the generic components configuration block.
	ComponentComponents = "components"
)

// DefaultComponentPaths provides the framework's default path map for core components.
// This map is used as the lowest priority base for path resolution.
// It can be inspected or copied by users to build their own custom path maps.
var DefaultComponentPaths = map[string]string{
	ComponentLogger:     "logger",
	ComponentRegistries: "registries",
}
