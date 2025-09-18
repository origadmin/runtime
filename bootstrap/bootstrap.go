package bootstrap

import (
	"fmt"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/goexts/generic/configure"

	sourcev1 "github.com/origadmin/runtime/api/gen/go/source/v1"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/config/file"
)

// DefaultBootstrapPath the default bootstrap configuration file path
const DefaultBootstrapPath = "configs/bootstrap.toml"

// DefaultSources create a default configuration source definition
func DefaultSources() *sourcev1.Sources {
	return &sourcev1.Sources{
		Sources: []*sourcev1.SourceConfig{
			{
				Name:     "default",
				Type:     "file",
				Priority: 10,
				Config: &sourcev1.SourceConfig_File{
					File: &sourcev1.FileSource{
						Path:   "configs/config.toml",
						Format: "toml",
						Reload: true,
					},
				},
			},
		},
	}
}

// LoadDefault loads the bootstrap configuration for the default path
func LoadDefault() (kratosconfig.Config, error) {
	return Load(DefaultBootstrapPath)
}

// LoadDefaultWithOptions loads the default bootstrap configuration with custom bootstrap options
func LoadDefaultWithOptions(options ...Option) (kratosconfig.Config, error) {
	return LoadWithOptions(DefaultBootstrapPath, options...)
}

// Load initializes the configuration by loading a bootstrap file, which defines
// the configuration sources to be used. It then delegates the creation and
// merging of these sources to the `runtime/config` package, which now handles
// priority sorting.
//
// This function serves as the primary entry point for loading the application's
// configuration.
func Load(bootstrapPath string, options ...runtimeconfig.Option) (kratosconfig.Config, error) {
	return LoadWithOptions(bootstrapPath, WithConfigOptions(options...))
}

// LoadWithOptions allows loading bootstrap configuration with both bootstrap config options
// and runtime config options. This function provides more flexibility when you need to
// customize both the bootstrap loading process and the final configuration creation.
func LoadWithOptions(bootstrapPath string, options ...Option) (kratosconfig.Config, error) {
	// 1. Load the bootstrap file using Kratos file source with provided options.
	bootstrapOptions := []kratosconfig.Option{
		kratosconfig.WithSource(file.NewSource(bootstrapPath)),
	}
	opts := configure.Apply(&Options{
		BootstrapOptions: bootstrapOptions,
	}, options)

	bootstrapConf := kratosconfig.New(opts.BootstrapOptions...)
	if err := bootstrapConf.Load(); err != nil {
		return nil, fmt.Errorf("failed to load bootstrap config file %s: %w", bootstrapPath, err)
	}

	// 2. Unmarshal the bootstrap config into a `sourcev1.Sources` struct.
	var sourcesDef sourcev1.Sources
	if err := bootstrapConf.Scan(&sourcesDef); err != nil {
		return nil, fmt.Errorf("failed to scan bootstrap config: %w", err)
	}
	// Close the bootstrap config as it's no longer needed.
	defer bootstrapConf.Close()

	// 2. Validate the sources definition.
	if err := validateSources(&sourcesDef); err != nil {
		return nil, fmt.Errorf("invalid sources definition: %w", err)
	}

	// 3. Delegate the creation and merging of sources to `runtimeconfig.NewConfig`.
	// The priority sorting logic is now handled within `runtimeconfig.NewConfig`.
	finalConfig, err := runtimeconfig.NewConfig(&sourcesDef, opts.ConfigOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create final config from sources: %w", err)
	}

	return finalConfig, nil
}

func LoadFromSource(source *sourcev1.SourceConfig, options ...Option) (kratosconfig.Config, error) {
	if source == nil {
		return nil, fmt.Errorf("source cannot be nil")
	}

	if source.GetType() != "file" {
		return nil, fmt.Errorf("source type must be file, but got %s", source.GetType())
	}

	return LoadWithOptions(source.GetFile().GetPath(), options...)
}

// LoadAndScan loads configuration and scans it into the specified struct
func LoadAndScan(bootstrapPath string, target interface{}) error {
	config, err := Load(bootstrapPath)
	if err != nil {
		return err
	}
	defer config.Close()

	if err := config.Scan(target); err != nil {
		return fmt.Errorf("failed to scan config into target: %w", err)
	}

	return nil
}

// validateSources validates the effectiveness of configuration source definitions
func validateSources(sources *sourcev1.Sources) error {
	if sources == nil {
		return fmt.Errorf("sources definition cannot be nil")
	}

	// Check if configuration sources are defined
	if len(sources.Sources) == 0 {
		return fmt.Errorf("no configuration sources defined")
	}

	// Check required fields for each configuration source
	for i, source := range sources.Sources {
		if source == nil {
			return fmt.Errorf("source #%d cannot be nil", i)
		}
		if source.Name == "" {
			return fmt.Errorf("source #%d must have a name", i)
		}
		if source.Type == "" {
			return fmt.Errorf("source #%d must have a type", i)
		}
	}

	return nil
}
