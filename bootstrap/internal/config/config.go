package config

import (
	"errors"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/runtime/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/runtime/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	"github.com/origadmin/runtime/bootstrap/constant"
	"github.com/origadmin/runtime/interfaces"
)

// structuredConfigImpl implements the interfaces.StructuredConfig interface.
// It wraps a generic interfaces.Config and provides type-safe, path-based decoding methods.
type structuredConfigImpl struct {
	interfaces.Config // Embed the generic config interface
	paths             map[string]string
}

// Statically assert that structuredConfigImpl implements the full StructuredConfig interface.
var _ interfaces.StructuredConfig = (*structuredConfigImpl)(nil)

// NewStructured creates a new structured config implementation.
// It takes a generic interfaces.Config and a path map to provide high-level decoding methods.
func NewStructured(cfg interfaces.Config, paths map[string]string) interfaces.StructuredConfig {
	if paths == nil {
		paths = make(map[string]string)
	}
	return &structuredConfigImpl{
		Config: cfg,
		paths:  paths,
	}
}

// decodeComponent implements a simple and robust decoding logic.
// It no longer contains any fallback logic. It trusts the `paths` map provided by the bootstrap package.
func (c *structuredConfigImpl) decodeComponent(componentKey string, value any) error {
	path, ok := c.paths[componentKey]

	// If the key is not in the paths map, or the path is explicitly empty, it's considered disabled or not configured.
	if !ok || path == "" {
		return nil // This is not an error.
	}

	// Attempt to decode using the provided path.
	err := c.Decode(path, value)
	if err != nil {
		// If the error is specifically a "not found" error, it's not a fatal issue.
		if errors.Is(err, kratosconfig.ErrNotFound) {
			return nil
		}
		// Any other error (e.g., parsing) is a real problem.
		return err
	}
	return nil
}

// DecodeApp implements the AppConfigDecoder interface.
func (c *structuredConfigImpl) DecodeApp() (*appv1.App, error) {
	appConfig := new(appv1.App)
	if err := c.decodeComponent(constant.ConfigApp, appConfig); err != nil {
		return nil, err
	}
	// If the struct is still zero-valued, it means the key was not found or disabled.
	if appConfig.Id == "" && appConfig.Name == "" {
		return nil, nil
	}
	return appConfig, nil
}

// DecodeLogger implements the LoggerConfigDecoder interface.
func (c *structuredConfigImpl) DecodeLogger() (*loggerv1.Logger, error) {
	loggerConfig := new(loggerv1.Logger)
	if err := c.decodeComponent(constant.ComponentLogger, loggerConfig); err != nil {
		return nil, err
	}
	if loggerConfig.Name == "" && len(loggerConfig.Level) == 0 {
		return nil, nil
	}
	return loggerConfig, nil
}

// DecodeDiscoveries implements the DiscoveriesConfigDecoder interface.
func (c *structuredConfigImpl) DecodeDiscoveries() (map[string]*discoveryv1.Discovery, error) {
	// Attempt 1: Try to decode directly into a map (if config is structured as a map).
	var discoveriesMap map[string]*discoveryv1.Discovery
	err := c.decodeComponent(constant.ComponentRegistries, &discoveriesMap)
	if err == nil {
		// If decoding into a map was successful and it's not empty, return it.
		if len(discoveriesMap) > 0 {
			return discoveriesMap, nil
		}
		// If it was successful but empty, proceed to try list or return nil.
	}

	// Attempt 2: If decoding into a map failed or resulted in an empty map,
	// try to decode into a list and convert it (if config is structured as a list).
	var discoveryList []*discoveryv1.Discovery
	err = c.decodeComponent(constant.ComponentRegistries, &discoveryList)
	if err != nil {
		// If both attempts failed, return the error from the second attempt.
		return nil, err
	}

	// If the list is empty, return nil.
	if len(discoveryList) == 0 {
		return nil, nil
	}

	// Convert the list to a map.
	discoveriesMap = make(map[string]*discoveryv1.Discovery)
	for _, d := range discoveryList {
		if d.Name != "" {
			discoveriesMap[d.Name] = d
		}
	}
	return discoveriesMap, nil
}

// DecodeMiddlewares implements the MiddlewareConfigDecoder interface.
// This implementation correctly preserves the user's fix.
func (c *structuredConfigImpl) DecodeMiddlewares() (*middlewarev1.Middlewares, error) {
	var middlewares *middlewarev1.Middlewares
	if err := c.decodeComponent(constant.ComponentMiddlewares, &middlewares); err != nil {
		return nil, err
	}
	// This check correctly handles both a nil pointer and an empty inner slice.
	if middlewares == nil || (len(middlewares.Middlewares) == 0) {
		return nil, nil
	}
	return middlewares, nil
}
