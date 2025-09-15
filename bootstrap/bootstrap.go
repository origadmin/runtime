package bootstrap

import (
	"fmt"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"

	sourcev1 "github.com/origadmin/runtime/api/gen/go/source/v1"
	runtimeconfig "github.com/origadmin/runtime/config"
)

// Load initializes the configuration by loading a bootstrap file, which defines
// the configuration sources to be used. It then delegates the creation and
// merging of these sources to the `runtime/config` package, which now handles
// priority sorting.
//
// This function serves as the primary entry point for loading the application's
// configuration.
func Load(bootstrapPath string) (kratosconfig.Config, error) {
	// 1. Load the bootstrap file using Kratos file source.
	bootstrapConf := kratosconfig.New(
		kratosconfig.WithSource(file.NewSource(bootstrapPath)),
	)
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

	// 3. Delegate the creation and merging of sources to `runtimeconfig.NewConfig`.
	// The priority sorting logic is now handled within `runtimeconfig.NewConfig`.
	finalConfig, err := runtimeconfig.NewConfig(&sourcesDef)
	if err != nil {
		return nil, fmt.Errorf("failed to create final config from sources: %w", err)
	}

	return finalConfig, nil
}
