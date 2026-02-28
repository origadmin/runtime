package bootstrap

import (
	"fmt"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	bootstrapconfig "github.com/origadmin/runtime/bootstrap/internal/config"
	"github.com/origadmin/runtime/contracts/constant"
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

	// 2. Load full configuration using the sources from bootstrap config
	bootstrap, cfg, err := LoadConfig(bootstrapPath, providerOpts) // Now providerOpts contains preloaded sources
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	log.Debugf("Load bootstrap config : %+v", bootstrap)
	appConfig := bootstrap.GetApp()
	if appConfig == nil {
		appConfig = &appv1.App{}
	}
	// 3. Create the final StructuredConfig.
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
	var businessConfig any = sc.DecodedConfig()
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
		config:           cfg,
		structuredConfig: sc,
		bootstrap:        bootstrap,
		businessConfig:   businessConfig,
	}
	return res, nil
}
