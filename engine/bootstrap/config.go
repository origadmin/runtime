package bootstrap

import (
	"fmt"
	"path/filepath"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"

	bootstrapv1 "github.com/origadmin/runtime/api/gen/go/config/bootstrap/v1"
	sourcev1 "github.com/origadmin/runtime/api/gen/go/config/source/v1"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/config/envsource"
	"github.com/origadmin/runtime/log"
)

// PathResolverFunc defines the signature for a function that resolves a configuration path.
type PathResolverFunc func(path string) string

// ConfigTransformer defines an interface for custom transformation of KConfig to a business configuration object.
type ConfigTransformer interface {
	Transform(runtimeconfig.KConfig) (any, error)
}

// ConfigTransformFunc is a function type that implements the ConfigTransformer interface.
type ConfigTransformFunc func(runtimeconfig.KConfig) (any, error)

// Transform implements the ConfigTransformer interface for ConfigTransformFunc.
func (f ConfigTransformFunc) Transform(config runtimeconfig.KConfig) (any, error) {
	return f(config)
}

// LoadConfig creates a new configuration decoder instance.
func LoadConfig(bootstrapPath string, providerOpts *ProviderOptions) (*bootstrapv1.Bootstrap, runtimeconfig.KConfig, error) {
	logger := log.NewHelper(log.DefaultLogger)

	var baseConfig runtimeconfig.KConfig
	var bootstrapConfig *bootstrapv1.Bootstrap
	var err error

	if providerOpts.config != nil {
		if cfg, ok := providerOpts.config.(runtimeconfig.KConfig); ok {
			baseConfig = cfg
		} else {
			return nil, nil, fmt.Errorf("engine: provided config is not a valid runtimeconfig.KConfig")
		}
	} else if providerOpts.directly {
		logger.Infof("Loading config directly from: %s", bootstrapPath)
		sources := []*sourcev1.SourceConfig{SourceWithFile(bootstrapPath)}
		baseConfig, err = runtimeconfig.New(&sourcev1.Sources{Configs: sources}, providerOpts.frameworkOptions...)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create base config for direct loading: %w", err)
		}
	} else {
		bootstrapConfig = loadBootstrapFromFile(bootstrapPath, providerOpts, logger)

		if bootstrapConfig == nil || len(bootstrapConfig.GetSources()) == 0 {
			return nil, nil, fmt.Errorf("no configuration sources found in bootstrap file: %s", bootstrapPath)
		}

		baseConfig, err = runtimeconfig.New(&sourcev1.Sources{Configs: bootstrapConfig.GetSources()}, providerOpts.frameworkOptions...)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create base config from bootstrap sources: %w", err)
		}
	}

	logger.Info("All configuration sources prepared, starting final load...")
	if err := baseConfig.Load(); err != nil {
		return nil, nil, err
	}

	return bootstrapConfig, baseConfig, nil
}

func isFileSource(source *sourcev1.SourceConfig) bool {
	return source != nil && source.Type == "file" && source.File != nil
}

func loadBootstrapFromFile(bootstrapPath string, providerOpts *ProviderOptions, logger *log.Helper) *bootstrapv1.Bootstrap {
	bootstrapCfg, err := LoadBootstrapConfig(bootstrapPath, providerOpts.prefixes...)
	if err != nil || bootstrapCfg == nil || len(bootstrapCfg.GetSources()) == 0 {
		if err != nil {
			logger.Warnf("Failed to load or parse sources from bootstrap file: %v", err)
		}
		return nil
	}

	bootstrapDir := filepath.Dir(bootstrapPath)
	for i, source := range bootstrapCfg.Sources {
		if !isFileSource(source) {
			continue
		}
		path := source.File.Path
		resolvedPath := path

		if providerOpts.pathResolver != nil {
			resolvedPath = providerOpts.pathResolver(path)
		} else if !filepath.IsAbs(path) {
			resolvedPath = filepath.Join(bootstrapDir, path)
		}

		source.File.Path = resolvedPath
		bootstrapCfg.Sources[i] = source
		logger.Infof("Load bootstrap file: %s", resolvedPath)
	}
	return bootstrapCfg
}

func bootstrapSources(path string, prefixes ...string) kratosconfig.Option {
	return kratosconfig.WithSource(
		file.NewSource(path),
		envsource.NewSource(prefixes...),
	)
}

func LoadBootstrapConfig(bootstrapPath string, prefixes ...string) (*bootstrapv1.Bootstrap, error) {
	configOpts := []kratosconfig.Option{bootstrapSources(bootstrapPath, prefixes...)}
	bootConfig := kratosconfig.New(configOpts...)

	defer func() {
		if err := bootConfig.Close(); err != nil {
			log.Errorf("failed to close temporary bootstrap config: %v", err)
		}
	}()

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
