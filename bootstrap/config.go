// Package bootstrap implements the functions, types, and interfaces for the module.
package bootstrap

import (
	"fmt"
	"path/filepath"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	bootstrapv1 "github.com/origadmin/runtime/api/gen/go/runtime/bootstrap/v1"
	sourcev1 "github.com/origadmin/runtime/api/gen/go/runtime/source/v1"
	"github.com/origadmin/runtime/bootstrap/constant"
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
	constant.ConfigApp:                "app",
	constant.ComponentLogger:          "logger",
	constant.ComponentRegistries:      "discoveries",
	constant.ComponentDefaultRegistry: "default_registry_name",
	constant.ComponentMiddlewares:     "middlewares",
	constant.ComponentServers:         "servers",
	constant.ComponentClients:         "clients",
}

// --- Options for LoadConfig ---

// PathResolverFunc defines the signature for a function that resolves a configuration path.
// It takes the base directory (of the bootstrap file) and the path from the config source
// and returns the final, resolved path.
type PathResolverFunc func(baseDir, path string) string

// ConfigTransformer defines an interface for custom transformation of kratosconfig.Config to interfaces.Config.
type ConfigTransformer interface {
	Transform(interfaces.Config, interfaces.StructuredConfig) (interfaces.StructuredConfig, error)
}

// ConfigTransformFunc is a function type that implements the ConfigTransformer interface.
type ConfigTransformFunc func(interfaces.Config, interfaces.StructuredConfig) (interfaces.StructuredConfig, error)

// Transform implements the ConfigTransformer interface for ConfigTransformFunc.
func (f ConfigTransformFunc) Transform(config interfaces.Config, sc interfaces.StructuredConfig) (
	interfaces.StructuredConfig, error) {
	return f(config, sc)
}

// LoadConfig creates a new configuration decoder instance.
// It orchestrates the entire configuration decoding process, following a clear, layered approach.
func LoadConfig(bootstrapPath string, providerOpts *ProviderOptions) (interfaces.Config, error) {
	logger := log.NewHelper(log.FromOptions(providerOpts.rawOptions))
	// 1. Apply Options to determine the configuration flow.

	var baseConfig interfaces.Config

	// Case 1: A fully custom interfaces.Config is provided.
	if providerOpts.config != nil { // The user has provided a pre-configured config instance.
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
			sources = append(sources, WithFileSource(fullBootstrapPath))
		}

		// Create the base config from the collected and resolved sources.
		var err error
		baseConfig, err = runtimeconfig.NewConfig(&sourcev1.Sources{Sources: sources}, providerOpts.rawOptions...)
		if err != nil {
			return nil, fmt.Errorf("failed to create base config: %w", err)
		}
	}

	// Step 2: Load the configuration. This is the responsibility of the bootstrap module.
	if err := baseConfig.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return baseConfig, nil
}

// loadSourcesFromBootstrapFile attempts to load a bootstrap configuration and resolve the paths of its sources.
// It returns a slice of source configurations or nil if loading fails or no sources are found.
func loadSourcesFromBootstrapFile(fullBootstrapPath string, providerOpts *ProviderOptions, logger *log.Helper) []*sourcev1.SourceConfig {
	bootstrapCfg, err := LoadBootstrapConfig(fullBootstrapPath, providerOpts.rawOptions...)
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

func WithFileSource(path string) *sourcev1.SourceConfig {
	return &sourcev1.SourceConfig{
		Type: "file",
		Config: &sourcev1.SourceConfig_File{
			File: &sourcev1.FileSource{
				Path: path,
			},
		},
	}
}
