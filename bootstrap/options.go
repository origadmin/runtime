package bootstrap

import (
	"github.com/origadmin/runtime/optionutil"

	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
)

// Option defines the function signature for configuration options used in the bootstrap process.
type Option = options.Option

// ProviderOptions holds all the configurable settings for the bootstrap provider.
// It is populated by applying a series of Option functions.
type ProviderOptions struct {
	appInfo            *interfaces.AppInfo
	componentFactories map[string]interfaces.ComponentFactory
	config             interfaces.Config
	configTransformer  ConfigTransformer
	defaultPaths       map[string]string
	directory          string
	directly           bool
	pathResolver       PathResolverFunc
	bootstrapPrefix    string
	rawOptions         []Option
}

// WithAppInfo provides application metadata to the bootstrap process.
// This information is merged with any metadata loaded from the configuration source.
func WithAppInfo(info *interfaces.AppInfo) Option {
	return optionutil.Update(func(o *ProviderOptions) {
		o.appInfo = info
	})
}

// WithComponent registers a component factory to be used during bootstrap.
func WithComponent(key string, factory interfaces.ComponentFactory) Option {
	return optionutil.Update(func(o *ProviderOptions) {
		if o.componentFactories == nil {
			o.componentFactories = make(map[string]interfaces.ComponentFactory)
		}
		o.componentFactories[key] = factory
	})
}

// WithConfig provides a pre-initialized configuration instance.
// If this option is used, the bootstrap process will skip loading configuration from files.
func WithConfig(cfg interfaces.Config) Option {
	return optionutil.Update(func(o *ProviderOptions) {
		o.config = cfg
	})
}

// WithConfigTransformer allows providing an object that implements the ConfigTransformer interface.
// This provides a flexible way to customize the creation of the StructuredConfig.
func WithConfigTransformer(transformer ConfigTransformer) Option {
	return optionutil.Update(func(o *ProviderOptions) {
		o.configTransformer = transformer
	})
}

// WithDefaultPaths provides a default path map for components.
// This map is merged with the framework's defaults, with these values taking precedence.
func WithDefaultPaths(paths map[string]string) Option {
	return optionutil.Update(func(o *ProviderOptions) {
		if o.defaultPaths == nil {
			o.defaultPaths = make(map[string]string)
		}
		for k, v := range paths {
			o.defaultPaths[k] = v
		}
	})
}

// WithDirectory sets the base directory for resolving relative paths in configuration files.
// If not set, paths are resolved relative to the current working directory.
func WithDirectory(dir string) Option {
	return optionutil.Update(func(o *ProviderOptions) {
		o.directory = dir
	})
}

// WithDirectly sets whether to treat the bootstrapPath as a direct configuration source,
// ignoring any `sources` defined within it.
func WithDirectly(directly bool) Option {
	return optionutil.Update(func(o *ProviderOptions) {
		o.directly = directly
	})
}

// WithPathResolver provides a custom function to resolve relative paths in configuration sources.
func WithPathResolver(resolver PathResolverFunc) Option {
	return optionutil.Update(func(o *ProviderOptions) {
		o.pathResolver = resolver
	})
}

// WithBootstrapPrefix sets the prefix for environment variables that can override
// settings in the initial bootstrap configuration file (e.g., `bootstrap.yaml`).
func WithBootstrapPrefix(prefix string) Option {
	return optionutil.Update(func(o *ProviderOptions) {
		o.bootstrapPrefix = prefix
	})
}

// FromOptions creates a ProviderOptions struct from a slice of Option functions.
// This is the single entry point for processing all bootstrap options.
func FromOptions(opts ...Option) *ProviderOptions {
	po := optionutil.NewT[ProviderOptions]()
	po.rawOptions = opts
	return po
}
