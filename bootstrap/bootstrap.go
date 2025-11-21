package bootstrap

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log" // Keep for error logging during config close

	bootstrapconfig "github.com/origadmin/runtime/bootstrap/internal/config"
)

// New creates a new component provider, which is the main entry point for application startup.
// It orchestrates the entire process of configuration loading and component initialization.
// It now returns the Result interface.
func New(bootstrapPath string, opts ...Option) (res Result, err error) {
	// 1. Apply options to get access to WithAppInfo.
	// Assuming options have been flattened as per our discussion.
	providerOpts := FromOptions(opts...)

	// 2. Load configuration.
	cfg, err := LoadConfig(bootstrapPath, providerOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err) // Early exit, no cleanup needed yet.
	}

	// --- Create StructuredConfig ---
	// Step 1: Merge default and user-provided paths.
	paths := make(map[string]string, len(defaultComponentPaths))
	for k, v := range defaultComponentPaths {
		paths[k] = v
	}
	if providerOpts.defaultPaths != nil {
		for component, path := range providerOpts.defaultPaths {
			paths[component] = path
		}
	}

	// Step 2: Create the base structured config implementation.
	// Step 3: (Optional) Apply a high-level transformer if provided.
	sc := bootstrapconfig.NewStructured(cfg, paths)
	if providerOpts.configTransformer != nil {
		sc, err = providerOpts.configTransformer.Transform(cfg, sc)
		if err != nil {
			// Ensure config is closed on error
			if closeErr := cfg.Close(); closeErr != nil {
				log.Errorf("failed to close config after transform error: %v", closeErr)
			}
			return nil, fmt.Errorf("failed to transform config: %w", err)
		}
	}

	configAppInfo, err := sc.DecodeApp()
	if err != nil {
		log.Debugf("failed to decode app info from config, will rely on WithAppInfo option or defaults: %v", err)
	}
	// 4. Assemble and return the final result.
	res = &resultImpl{
		config:           cfg,
		structuredConfig: sc,
		appInfo:          configAppInfo,
	}
	return res, nil
}
