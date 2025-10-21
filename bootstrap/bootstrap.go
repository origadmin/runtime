package bootstrap

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
)

// New creates a new component provider, which is the main entry point for application startup.
// It orchestrates the entire process of configuration loading and component initialization.
// It now returns the Result interface.
func New(bootstrapPath string, opts ...Option) (res Result, err error) {
	// 1. Apply options to get access to WithAppInfo.
	// Assuming options have been flattened as per our discussion.
	providerOpts := FromOptions(opts...)

	// 2. Load configuration.
	cfg, sc, err := LoadConfig(bootstrapPath, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err) // Early exit, no cleanup needed yet.
	}

	// Define the core cleanup logic once.
	cleanupFunc := func() {
		if cfg != nil {
			if closeErr := cfg.Close(); closeErr != nil {
				log.Errorf("failed to close config during cleanup: %v", closeErr)
			}
		}
	}

	// Defer the cleanup logic. It will only run if the function returns an error.
	defer func() {
		if err != nil {
			cleanupFunc()
		}
	}()

	// 3. Merge and validate application info from config and options.
	finalAppInfo, err := mergeAppInfo(sc, providerOpts.appInfo)
	if err != nil {
		return nil, err // mergeAppInfo already wraps the error
	}

	// 4. Build the component container.
	c, logger, err := buildContainer(sc, providerOpts.componentFactories, opts...)
	if err != nil {
		return nil, err // buildContainer already wraps the error
	}

	// 5. Assemble and return the final result.
	res = &resultImpl{
		config:           cfg,
		structuredConfig: sc,
		appInfo:          finalAppInfo,
		container:        c,
		logger:           logger,
		cleanup:          cleanupFunc,
	}
	return res, nil
}
