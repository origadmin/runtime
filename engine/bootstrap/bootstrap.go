package bootstrap

import (
	"fmt"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	bootstrapv1 "github.com/origadmin/runtime/api/gen/go/config/bootstrap/v1"
	"github.com/origadmin/runtime/log"
)

// New creates a new component provider, which is the main entry point for application startup.
// It orchestrates the entire process of configuration loading.
// It now returns the Result interface, which contains configuration-related data and the raw App protobuf message.
func New(bootstrapPath string, opts ...Option) (res Result, err error) {
	// 1. Apply bootstrap options.
	providerOpts := FromOptions(opts...)

	// 2. Load full configuration using the sources from bootstrap config
	bootstrap, cfg, err := LoadConfig(bootstrapPath, providerOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	log.Debugf("Load bootstrap config : %+v", bootstrap)

	// Ensure app info exists
	if bootstrap == nil {
		bootstrap = &bootstrapv1.Bootstrap{}
	}
	if bootstrap.App == nil {
		bootstrap.App = &appv1.App{}
	}

	// 3. Determine business config
	var businessConfig any

	// Priority 2: Automatic decoding into target struct
	if providerOpts.configTarget != nil {
		if err := cfg.Scan(providerOpts.configTarget); err != nil {
			if closeErr := cfg.Close(); closeErr != nil {
				log.Errorf("failed to close config after decode error: %v", closeErr)
			}
			return nil, fmt.Errorf("failed to scan config into target: %w", err)
		}
		businessConfig = providerOpts.configTarget
	}

	// Priority 1: Custom transformation logic (overrides target)
	if providerOpts.configTransformer != nil {
		businessConfig, err = providerOpts.configTransformer.Transform(cfg)
		if err != nil {
			if closeErr := cfg.Close(); closeErr != nil {
				log.Errorf("failed to close config after transform error: %v", closeErr)
			}
			return nil, fmt.Errorf("failed to transform config: %w", err)
		}
	}

	// 4. Assemble and return the final result.
	res = &resultImpl{
		config:         cfg,
		bootstrap:      bootstrap,
		businessConfig: businessConfig,
		configPath:     bootstrapPath,
	}
	return res, nil
}
