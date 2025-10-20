// Package bootstrap implements the functions, types, and interfaces for the module.
package bootstrap

import (
	"fmt"
	"path/filepath"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	bootstrapv1 "github.com/origadmin/runtime/api/gen/go/runtime/bootstrap/v1"
	sourcev1 "github.com/origadmin/runtime/api/gen/go/runtime/source/v1"
	"github.com/origadmin/runtime/bootstrap/constant"
	bootstrapconfig "github.com/origadmin/runtime/bootstrap/internal/config"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/config/envsource"
	"github.com/origadmin/runtime/config/file"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/log"
)

// defaultComponentPaths provides the framework's default path map for core components.
// It is now a private variable within the bootstrap package, ensuring that the default
// path logic is cohesive and contained within this package.
var defaultComponentPaths = map[string]string{
	constant.ConfigApp:            "app",
	constant.ComponentLogger:      "logger",
	constant.ComponentRegistries:  "discoveries", // Corrected to match full_config.yaml
	constant.ComponentMiddlewares: "middlewares",
}

// LoadConfig creates a new configuration decoder instance.
// It orchestrates the entire configuration decoding process, following a clear, layered approach.
func LoadConfig(bootstrapPath string, opts ...Option) (interfaces.StructuredConfig, error) {
	logger := log.NewHelper(log.FromOptions(opts))
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
		// Determine the full path for the initial bootstrap file.
		fullBootstrapPath := bootstrapPath
		if providerOpts.directory != "" && !filepath.IsAbs(bootstrapPath) {
			fullBootstrapPath = filepath.Join(providerOpts.directory, bootstrapPath)
		}

		var sources []*sourcev1.SourceConfig
		// If not in 'directly' mode, try to load sources from the bootstrap file.
		if !providerOpts.directly {
			sources = loadSourcesFromBootstrapFile(fullBootstrapPath, providerOpts, logger)
		}

		// Fallback or 'directly' mode: use the bootstrapPath as the single source.
		if len(sources) == 0 {
			logger.Infof("No sources found in bootstrap file, using it directly: %s", fullBootstrapPath)
			sources = append(sources, SourceWithFile(fullBootstrapPath))
		}

		// Create the base config from the collected and resolved sources.
		var err error
		baseConfig, err = runtimeconfig.NewConfig(&sourcev1.Sources{Sources: sources}, opts...)
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

// loadSourcesFromBootstrapFile attempts to load a bootstrap configuration and resolve the paths of its sources.
// It returns a slice of source configurations or nil if loading fails or no sources are found.
func loadSourcesFromBootstrapFile(fullBootstrapPath string, providerOpts *ConfigLoadOptions, logger *log.Helper) []*sourcev1.SourceConfig {
	bootstrapCfg, err := LoadBootstrapConfig(fullBootstrapPath)
	if err != nil || bootstrapCfg == nil || len(bootstrapCfg.GetSources()) == 0 {
		return nil
	}

	// Successfully loaded sources, now resolve their paths.
	var sources []*sourcev1.SourceConfig
	bootstrapDir := filepath.Dir(fullBootstrapPath) // Base for relative paths
	for _, source := range bootstrapCfg.GetSources() {
		if fileSource, ok := source.Config.(*sourcev1.SourceConfig_File); ok {
			path := fileSource.File.Path
			resolvedPath := path

			// Use custom path resolver if provided, otherwise use default logic.
			if providerOpts.pathResolver != nil {
				resolvedPath = providerOpts.pathResolver(bootstrapDir, path)
			} else if !filepath.IsAbs(path) {
				// Default logic: It's a relative path, join it with the bootstrap file's directory.
				resolvedPath = filepath.Join(bootstrapDir, path)
			}

			fileSource.File.Path = resolvedPath
			logger.Infof("Load bootstrap file: %s", resolvedPath)
		}
		sources = append(sources, source)
	}
	return sources
}

// LoadBootstrapConfig loads the bootstrapv1.Bootstrap definition from a local bootstrap configuration file.
// This function is the first step in the configuration process.
// The temporary config used to load the sources is closed internally.
func LoadBootstrapConfig(bootstrapPath string, opts ...Option) (*bootstrapv1.Bootstrap, error) {
	providerOpts := FromOptions(opts...)
	// Create a temporary Kratos config instance to load the bootstrap.yaml file.
	bootConfig := kratosconfig.New(
		kratosconfig.WithSource(file.NewSource(bootstrapPath), envsource.NewSource(providerOpts.bootstrapPrefix)),
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

func SourceWithFile(path string) *sourcev1.SourceConfig {
	return &sourcev1.SourceConfig{
		Type: "file",
		Config: &sourcev1.SourceConfig_File{
			File: &sourcev1.FileSource{
				Path: path,
			},
		},
	}
}
