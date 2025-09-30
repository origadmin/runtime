// Package bootstrap implements the functions, types, and interfaces for the module.
package bootstrap

import (
	"fmt"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goexts/generic/configure"

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
	// 1. Apply Options to determine the configuration flow.
	providerOpts := configure.New(opts)

	var (
		baseConfig interfaces.Config
		paths      map[string]string
	)

	// Case 1: A fully custom interfaces.Config is provided.
	if providerOpts.config != nil {
		// If it's already a StructuredConfig, the user has provided a complete, loaded implementation.
		if sc, ok := providerOpts.config.(interfaces.StructuredConfig); ok {
			return sc, nil
		}
		// Otherwise, we'll use it as the base for our default structured implementation.
		baseConfig = providerOpts.config
		paths = nil // Paths are not applicable for a fully custom config.

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

		// The runtimeconfig package now handles the creation of the base config object.
		// It returns an un-loaded interfaces.Config.
		baseConfig, err = runtimeconfig.NewConfig(sources, providerOpts.configOptions...)
		if err != nil {
			return nil, fmt.Errorf("failed to create base config: %w", err)
		}

		// Merge paths for the default flow.
		finalPaths := make(map[string]string)
		for component, path := range constant.DefaultComponentPaths {
			finalPaths[component] = path
		}
		if providerOpts.defaultPaths != nil {
			for component, path := range providerOpts.defaultPaths {
				finalPaths[component] = path
			}
		}
		paths = finalPaths
	}

	// Step 2: Load the configuration. This is the new responsibility of the bootstrap module.
	if err := baseConfig.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Final Step: All flows converge here.
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
