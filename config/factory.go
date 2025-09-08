package config

import (
	"fmt"

	configenv "github.com/go-kratos/kratos/v2/config/env"
	"github.com/goexts/generic/configure"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/log"
)

var (
	// defaultConfigFactory is the default config factory.
	defaultConfigFactory = NewFactory()
)

// configFactory is a config factory that implements interfaces.ConfigBuilder.
type configFactory struct {
	factory.Registry[interfaces.ConfigFactory]
}

// RegisterConfigFunc registers a new ConfigBuilder with the given name and function.
func (f *configFactory) RegisterConfigFunc(name string, buildFunc BuildFunc) {
	f.Register(name, buildFunc)
}

// BuildFunc is a function type that takes a KConfig and a list of Options and returns a Selector and an error.
type BuildFunc func(*configv1.SourceConfig, *interfaces.Options) (kratosconfig.Source, error)

// NewSource is a method that implements the ConfigBuilder interface for ConfigBuildFunc.
func (fn BuildFunc) NewSource(cfg *configv1.SourceConfig, opts *interfaces.Options) (kratosconfig.Source, error) {
	// Call the function with the given KConfig and a list of Options.
	return fn(cfg, opts)
}

// NewConfig creates a new Selector object based on the given KConfig and options.
func (f *configFactory) NewConfig(cfg *configv1.SourceConfig, opts ...interfaces.Option) (kratosconfig.Config, error) {
	options := configure.Apply(&interfaces.Options{}, opts) // Corrected: Use settings.Apply with a new interfaces.Options{}
	sources := options.Sources
	if sources == nil {
		sources = make([]kratosconfig.Source, 0)
	}

	for _, t := range cfg.Types {
		bld, ok := f.Get(t)
		if !ok {
			return nil, fmt.Errorf("unknown type: %s", t)
		}
		log.Infof("registering type: %s", t)
		source, err := bld.NewSource(cfg, options)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	if options.EnvPrefixes != nil {
		sources = append(sources, configenv.NewSource(options.EnvPrefixes...))
	}

	options.ConfigOptions = append(options.ConfigOptions, kratosconfig.WithSource(sources...))
	return kratosconfig.New(options.ConfigOptions...), nil
}

func (f *configFactory) SyncConfig(cfg *configv1.SourceConfig, v any, opts ...interfaces.Option) error {
	// This method is a placeholder. Actual synchronization logic would go here.
	// For now, we'll just return nil or an error if needed.
	return nil
}

// NewFactory creates a new config factory.
func NewFactory() interfaces.ConfigBuilder {
	return &configFactory{
		Registry: factory.New[interfaces.ConfigFactory](),
	}
}
