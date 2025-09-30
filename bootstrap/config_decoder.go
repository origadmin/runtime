// Package bootstrap implements the functions, types, and interfaces for the module.
package bootstrap

import (
	"fmt"

	sourcev1 "github.com/origadmin/runtime/api/gen/go/source/v1"
	"github.com/origadmin/runtime/bootstrap/constant"
	bootstrapConfig "github.com/origadmin/runtime/bootstrap/internal/config"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/interfaces"
)

// NewDecoder creates a new configuration decoder instance.
// It orchestrates the entire configuration decoding process, including path resolution and source merging.
// The returned interfaces.Config is ready to be consumed by New or other tools.
// It no longer returns *bootstrapv1.Bootstrap directly, as that should be scanned from the returned interfaces.Config.
func NewDecoder(bootstrapPath string, opts ...DecoderOption) (interfaces.Config, error) {
	// 1. Apply Options
	decoderOpts := &decoderOptions{}
	for _, o := range opts {
		o(decoderOpts)
	}

	// If a custom config is provided, use it directly.
	if decoderOpts.customConfig != nil {
		return decoderOpts.customConfig, nil
	}

	// If a direct Kratos config is provided, use it to create the interfaces.Config.
	if decoderOpts.kratosConfig != nil {
		// If a transformer is provided, use it to transform the Kratos config.
		if decoderOpts.configTransformer != nil {
			cfg, err := decoderOpts.configTransformer.Transform(decoderOpts.kratosConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to create interfaces.Config from kratosConfig using transformer: %w", err)
			}
			return cfg, nil
		}
		// Otherwise, use the default implementation.
		return bootstrapConfig.NewConfigImpl(decoderOpts.kratosConfig, nil), nil
	}

	// 2. Load BootstrapConfig from file
	bootstrapCfg, err := LoadBootstrapConfig(bootstrapPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load bootstrap config: %w", err)
	}

	// 3. Merge paths with a clear priority:
	//    1. DefaultComponentPaths (lowest)
	//    2. WithDefaultPaths option
	//    3. Paths from bootstrap.yaml (highest)

	// Start with a copy of the public default map
	finalPaths := make(map[string]string)
	for component, path := range constant.DefaultComponentPaths {
		finalPaths[component] = path
	}

	// Apply paths from WithDefaultPaths option
	if decoderOpts.defaultPaths != nil {
		for component, path := range decoderOpts.defaultPaths {
			finalPaths[component] = path
		}
	}

	// 4. Create the final Kratos config.Config from all sources
	var sources *sourcev1.Sources
	if bootstrapCfg != nil {
		sources = &sourcev1.Sources{Sources: bootstrapCfg.GetSources()}
	}
	// If no sources are found, it means bootstrap.yaml was empty or invalid, but we can still proceed
	// if a config factory is provided to handle the empty Kratos config.
	if sources == nil {
		sources = &sourcev1.Sources{}
	}

	finalKratosConfig, err := runtimeconfig.NewConfig(sources, decoderOpts.configOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create final kratos config: %w", err)
	}

	// CRITICAL FIX: Load the configuration after creating it
	if err := finalKratosConfig.Load(); err != nil {
		return nil, fmt.Errorf("failed to load final kratos config: %w", err)
	}

	// 5. Create the interfaces.Config implementation with the final merged paths
	// If a config transformer is provided, use it to transform the finalKratosConfig.
	if decoderOpts.configTransformer != nil {
		cfg, err := decoderOpts.configTransformer.Transform(finalKratosConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create interfaces.Config from finalKratosConfig using transformer: %w", err)
		}
		return cfg, nil
	}

	// Otherwise, use the default implementation.
	// We assume bootstrapConfig.NewConfigImpl exists in the internal package.
	decoder := bootstrapConfig.NewConfigImpl(finalKratosConfig, finalPaths)

	return decoder, nil
}
