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

// LoadConfig creates a new configuration decoder instance.
// It orchestrates the entire configuration decoding process, following a clear, layered approach.
func LoadConfig(bootstrapPath string, opts ...ConfigLoadOption) (interfaces.StructuredConfig, error) {
	// 1. Apply Options
	decoderOpts := &decoderOptions{}
	for _, o := range opts {
		o(decoderOpts)
	}

	var (
		baseConfig interfaces.Config
		paths      map[string]string
	)

	// Case 1: A fully custom interfaces.Config is provided.
	if decoderOpts.customConfig != nil {
		// If it's already a StructuredConfig, the user has provided a complete implementation.
		if sc, ok := decoderOpts.customConfig.(interfaces.StructuredConfig); ok {
			return sc, nil
		}
		// Otherwise, we'll use it as the base for our default structured implementation.
		baseConfig = decoderOpts.customConfig
		paths = nil // Paths are not applicable for a fully custom config.

		// Case 2: A direct Kratos config is provided.
	} else if decoderOpts.customConfig != nil {
		// A transformer takes precedence, allowing custom logic on the Kratos config.
		if decoderOpts.configTransformer != nil {
			sc, err := decoderOpts.configTransformer.Transform(decoderOpts.customConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to transform provided kratos config: %w", err)
			}
			return sc, nil
		}
		baseConfig = decoderOpts.customConfig
		paths = nil // Paths are not applicable here either.
		// Case 3: Default flow - load from bootstrapPath.
	} else {
		// Load the bootstrap file to get sources.
		bootstrapCfg, err := LoadBootstrapConfig(bootstrapPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load bootstrap config: %w", err)
		}

		// Create the final Kratos config from all sources.
		var sources *sourcev1.Sources
		if bootstrapCfg != nil {
			sources = &sourcev1.Sources{Sources: bootstrapCfg.GetSources()}
		}
		if sources == nil {
			sources = &sourcev1.Sources{}
		}

		// Create a ready-to-use, loaded, and adapted config object.
		// The runtimeconfig package now handles creation, loading, and adapting.
		finalConfig, err := runtimeconfig.NewConfig(sources, decoderOpts.configOptions...)
		if err != nil {
			return nil, fmt.Errorf("failed to create final config: %w", err)
		}
		if err := finalConfig.Load(); err != nil {
			return nil, fmt.Errorf("failed to load final kratos config: %w", err)
		}

		// A transformer can intercept the final Kratos config.
		if decoderOpts.configTransformer != nil {
			baseConfig, err = decoderOpts.configTransformer.Transform(finalConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to transform final kratos config: %w", err)
			}
		} else {
			// Default behavior: adapt the final Kratos config.
			baseConfig = finalConfig
		}

		// Merge paths for the default flow.
		finalPaths := make(map[string]string)
		for component, path := range constant.DefaultComponentPaths {
			finalPaths[component] = path
		}
		if decoderOpts.defaultPaths != nil {
			for component, path := range decoderOpts.defaultPaths {
				finalPaths[component] = path
			}
		}
		paths = finalPaths
	}

	// Final Step: All flows converge here.
	// We take the base interfaces.Config and enhance it with structured, path-based decoding.
	return bootstrapconfig.NewStructured(baseConfig, paths), nil
}

// LoadBootstrapConfig loads the bootstrapv1.Bootstrap definition from a local bootstrap configuration file.
// This function is the first step in the configuration process.
// The temporary config used to load the sources is closed internally.
func LoadBootstrapConfig(bootstrapPath string) (*bootstrapv1.Bootstrap, error) {
	// Create a temporary Kratos config instance to load the bootstrap.yaml file.
	// We assume runtimeconfig.NewFileSource exists and creates a file source.
	bootConfig := kratosconfig.New(
		kratosconfig.WithSource(file.NewSource(bootstrapPath)),
	)
	// Defer closing the config and handle its error
	defer func() {
		if err := bootConfig.Close(); err != nil {
			// Log the error, as we can't return it from a deferred function
			log.Errorf("failed to close temporary bootstrap config: %v", err)
		}
	}()

	// CRITICAL FIX: Load the config after creating it
	if err := bootConfig.Load(); err != nil {
		return nil, fmt.Errorf("failed to load temporary bootstrap config: %w", err)
	}

	var bc bootstrapv1.Bootstrap
	if err := bootConfig.Scan(&bc); err != nil { // Reverted to direct Scan
		return nil, fmt.Errorf("failed to scan bootstrap config from %s: %w", bootstrapPath, err)
	}

	return &bc, nil
}
