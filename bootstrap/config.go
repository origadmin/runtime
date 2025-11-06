// Package bootstrap implements the functions, types, and interfaces for the module.
package bootstrap

import (
	"fmt"
	"path/filepath"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	bootstrapv1 "github.com/origadmin/runtime/api/gen/go/config/bootstrap/v1"
	sourcev1 "github.com/origadmin/runtime/api/gen/go/config/source/v1"
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
	constant.ComponentData:            "data",
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
	var err error
	// Case 1: A fully custom interfaces.Config is provided.
	if providerOpts.config != nil { // The user has provided a pre-configured config instance.
		// Otherwise, we'll use it as the base for our default structured implementation.
		baseConfig = providerOpts.config
		// Case 2: Default flow - load from bootstrapPath.
	} else {

		var sources []*sourcev1.SourceConfig
		// If not in 'directly' mode, try to load sources from the bootstrap file.
		if !providerOpts.directly {
			sources = loadSourcesFromBootstrapFile(bootstrapPath, providerOpts, logger)
		}

		// Fallback or 'directly' mode: use the bootstrapPath as the single source.
		if len(sources) == 0 {
			logger.Infof("No sources found in bootstrap file, using it directly: %s", bootstrapPath)
			sources = append(sources, SourceWithFile(bootstrapPath))
		}

		// Create the base config from the collected and resolved sources.
		baseConfig, err = runtimeconfig.NewConfig(&sourcev1.Sources{Configs: sources}, providerOpts.rawOptions...)
		if err != nil {
			return nil, fmt.Errorf("failed to create base config: %w", err)
		}
	}

	// Step 2: Load the configuration. This is the responsibility of the bootstrap module.
	if err := baseConfig.Load(); err != nil {
		return nil, err
	}

	return baseConfig, nil
}

func isFileSource(source *sourcev1.SourceConfig) bool {
	return source != nil && source.Type == "file" && source.File != nil
}

// loadSourcesFromBootstrapFile attempts to load a bootstrap configuration and resolve the paths of its sources.
// It returns a slice of source configurations or nil if loading fails or no sources are found.
func loadSourcesFromBootstrapFile(bootstrapPath string, providerOpts *ProviderOptions, logger *log.Helper) []*sourcev1.SourceConfig {
	bootstrapCfg, err := LoadBootstrapConfig(bootstrapPath, providerOpts.rawOptions...)
	if err != nil || len(bootstrapCfg.GetSources()) == 0 {
		if err != nil {
			logger.Warnf("Failed to load or parse sources from bootstrap file: %v", err)
		}
		return nil
	}

	// Successfully loaded sources, now resolve their paths.
	bootstrapDir := filepath.Dir(bootstrapPath) // Base for relative paths
	for i, source := range bootstrapCfg.Sources {
		if !isFileSource(source) {
			continue
		}
		path := source.File.Path
		resolvedPath := path

		// Use custom path resolver if provided, otherwise use default logic.
		if providerOpts.pathResolver != nil {
			resolvedPath = providerOpts.pathResolver(bootstrapDir, path)
		} else if !filepath.IsAbs(path) {
			// Default logic: It's a relative path, join it with the bootstrap file's directory.
			resolvedPath = filepath.Join(bootstrapDir, path)
		}

		source.File.Path = resolvedPath
		bootstrapCfg.Sources[i] = source
		logger.Infof("Load bootstrap file: %s", resolvedPath)
	}
	return bootstrapCfg.Sources
}

func bootstrapSources(path string, prefixes ...string) kratosconfig.Option {
	return kratosconfig.WithSource(
		file.NewSource(path),
		envsource.NewSource(prefixes...),
	)
}

// LoadBootstrapConfig loads the bootstrapv1.Bootstrap definition from a local bootstrap configuration file.
// This function is the first step in the configuration process.
// The temporary config used to load the sources is closed internally.
func LoadBootstrapConfig(bootstrapPath string, opts ...Option) (*bootstrapv1.Bootstrap, error) {
	providerOpts := FromOptions(opts...)
	// Create a temporary Kratos config instance to load the bootstrap.yaml file.
	bootConfig := kratosconfig.New(bootstrapSources(bootstrapPath, providerOpts.prefixes...))

	// Defer closing the config and handle its error
	defer func() {
		if err := bootConfig.Close(); err != nil {
			log.Errorf("failed to close temporary bootstrap config: %v", err)
		}
	}()

	// Load the config to read the bootstrap file.
	if err := bootConfig.Load(); err != nil {
		return nil, err
	}

	var bc bootstrapv1.Bootstrap
	if err := bootConfig.Scan(&bc); err != nil {
		return nil, err
	}

	return &bc, nil
}

func SourceWithFile(path string) *sourcev1.SourceConfig {
	return &sourcev1.SourceConfig{
		Type: "file",
		File: &sourcev1.FileSource{
			Path: path,
		},
	}
}
