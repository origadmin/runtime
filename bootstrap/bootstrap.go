package bootstrap

import (
	"fmt"

	bootstrapv1 "github.com/origadmin/runtime/api/gen/go/config/bootstrap/v1" // Import bootstrapv1
	bootstrapconfig "github.com/origadmin/runtime/bootstrap/internal/config"
	"github.com/origadmin/runtime/interfaces/constant"
	"github.com/origadmin/runtime/log"
)

// defaultComponentPaths provides the framework's default path map for core components.
var defaultComponentPaths = map[constant.ComponentKey]string{
	constant.ConfigApp:                "app",
	constant.ComponentLogger:          "logger",
	constant.ComponentData:            "data",
	constant.ComponentDatabases:       "databases",
	constant.ComponentCaches:          "caches",
	constant.ComponentObjectStores:    "object_stores",
	constant.ComponentRegistries:      "discoveries",
	constant.ComponentDefaultRegistry: "default_registry_name",
	constant.ComponentMiddlewares:     "middlewares",
	constant.ComponentServers:         "servers",
	constant.ComponentClients:         "clients",
}

// New creates a new component provider, which is the main entry point for application startup.
// It orchestrates the entire process of configuration loading.
// It now returns the Result interface, which contains configuration-related data and the raw App protobuf message.
func New(bootstrapPath string, opts ...Option) (res Result, err error) {
	// 1. Apply bootstrap options.
	providerOpts := FromOptions(opts...)

	// 2. Load configuration from all sources (local and remote).
	cfg, err := LoadConfig(bootstrapPath, providerOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 3. Decode the Bootstrap message to get the App field.
	var bootstrapCfg bootstrapv1.Bootstrap
	if err := cfg.Decode("", &bootstrapCfg); err != nil {
		// Log the error but continue, as App field in bootstrap config is optional.
		log.Warnf("failed to decode bootstrap config, App field might be missing: %v", err)
	}
	app := bootstrapCfg.GetApp()

	// 4. Create the final StructuredConfig.
	paths := make(map[constant.ComponentKey]string, len(defaultComponentPaths))
	for k, v := range defaultComponentPaths {
		paths[k] = v
	}
	if providerOpts.defaultPaths != nil {
		for component, path := range providerOpts.defaultPaths {
			paths[component] = path
		}
	}
	sc := bootstrapconfig.NewStructured(cfg, paths)
	if providerOpts.configTransformer != nil {
		sc, err = providerOpts.configTransformer.Transform(cfg, sc)
		if err != nil {
			if closeErr := cfg.Close(); closeErr != nil {
				log.Errorf("failed to close config after transform error: %v", closeErr)
			}
			return nil, fmt.Errorf("failed to transform config: %w", err)
		}
	}

	// 5. Assemble and return the final result.
	res = &resultImpl{
		config:           cfg,
		structuredConfig: sc,
		appConfig:        app,
	}
	return res, nil
}
