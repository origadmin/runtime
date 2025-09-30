// Package bootstrap implements the functions, types, and interfaces for the module.
package bootstrap

import (
	"fmt"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"

	bootstrapv1 "github.com/origadmin/runtime/api/gen/go/bootstrap/v1"
	sourcev1 "github.com/origadmin/runtime/api/gen/go/source/v1"
	"github.com/origadmin/runtime/bootstrap/constant"
	bootstrapconfig "github.com/origadmin/runtime/bootstrap/internal/config"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/config/file"
	"github.com/origadmin/runtime/interfaces"
)

// defaultComponentPaths provides the framework's default path map for core components.
// It is now a private variable within the bootstrap package, ensuring that the default
// path logic is cohesive and contained within this package.
var defaultComponentPaths = map[string]string{
	constant.ConfigApp:            "app",
	constant.ComponentLogger:      "logger",
	constant.ComponentRegistries:  "registries",
	constant.ComponentMiddlewares: "middlewares",
}

// LoadConfig creates a new configuration decoder instance.
// It orchestrates the entire configuration decoding process, following a clear, layered approach.
func LoadConfig(bootstrapPath string, opts ...Option) (interfaces.StructuredConfig, error) {
	// 1. Apply Options to determine the configuration flow.
	providerOpts := FromConfigLoadOptions(opts...)

	var baseConfig interfaces.Config

	// Case 1: A fully custom interfaces.Config is provided.
	if providerOpts.config != nil {
		// If it's already a StructuredConfig, the user has provided a complete, loaded implementation.
		if sc, ok := providerOpts.config.(interfaces.StructuredConfig); ok {
			return sc, nil
		}
		// Otherwise, we'll use it as the base for our default structured implementation.
		baseConfig = providerOpts.config

		// Case 2: Default flow - load from bootstrapPath.
	} else {
		// Load the bootstrap file to get sources.
		bootstrapCfg, err := LoadBootstrapConfig(bootstrapPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load bootstrap config: %w", err)
		}

		var sources *sourcev1.Sources
		if bootstrapCfg != nil {
			sources = &sourcev1.Sources{Sources: bootstrapCfg.GetSources()}
		}

		baseConfig, err = runtimeconfig.NewConfig(sources, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create base config: %w", err)
		}
	}

	// Step 2: Merge paths. This logic now applies to all flows that provide a baseConfig.
	// We create a new map from our private default paths and then merge user-provided paths on top.
	paths := make(map[string]string, len(defaultComponentPaths))
	for k, v := range defaultComponentPaths {
		paths[k] = v
	}
	if providerOpts.defaultPaths != nil {
		for component, path := range providerOpts.defaultPaths {
			paths[component] = path
		}
	}

	// Step 3: Load the configuration. This is the responsibility of the bootstrap module.
	if err := baseConfig.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Step 4: (Optional) Apply a high-level transformer if provided.
	if providerOpts.configTransformer != nil {
		return providerOpts.configTransformer.Transform(baseConfig)
	}

	// Final Step: If no transformer is used, apply the default structured implementation.
	// We take the loaded base interfaces.Config and enhance it with structured, path-based decoding.
	return bootstrapconfig.NewStructured(baseConfig, paths), nil
}

// LoadBootstrapConfig loads the bootstrapv1.Bootstrap definition from a local bootstrap configuration file.
// This function is the first step in the configuration process.
// The temporary config used to load the sources is closed internally.
func LoadBootstrapConfig(bootstrapPath string) (*bootstrapv1.Bootstrap, error) {
	// Create a temporary Kratos config instance to load the bootstrap.yaml file.
	bootConfig := kratosconfig.New(
		kratosconfig.WithSource(file.NewSource(bootstrapPath)),
	)

	// Defer closing the config and handle its error
	defer func() {
		if err := bootConfig.Close(); err != nil {
			log.Errorf("failed to close temporary bootstrap config: %v", err)
		}
	}()

	// Load the config to read the bootstrap file.
	if err := bootConfig.Load(); err != nil {
		return nil, fmt.Errorf("failed to load temporary bootstrap config: %w", err)
	}

	var bc bootstrapv1.Bootstrap
	if err := bootConfig.Scan(&bc); err != nil {
		return nil, fmt.Errorf("failed to scan bootstrap config from %s: %w", bootstrapPath, err)
	}

	return &bc, nil
}
