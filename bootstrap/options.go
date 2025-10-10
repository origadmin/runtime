package bootstrap

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"

	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/optionutil"
)

// --- Options for LoadConfig ---

// ConfigTransformer defines an interface for custom transformation of kratosconfig.Config to interfaces.Config.
type ConfigTransformer interface {
	Transform(interfaces.Config) (interfaces.StructuredConfig, error)
}

// ConfigTransformFunc is a function type that implements the ConfigTransformer interface.
type ConfigTransformFunc func(kratosconfig.Config) (interfaces.Config, error)

// Transform implements the ConfigTransformer interface for ConfigTransformFunc.
func (f ConfigTransformFunc) Transform(config kratosconfig.Config) (interfaces.Config, error) {
	return f(config)
}

// ConfigLoadOptions holds configuration for the LoadConfig function.
type ConfigLoadOptions struct {
	defaultPaths      map[string]string
	config            interfaces.Config // Added: Custom interfaces.Config implementation
	configTransformer ConfigTransformer // Custom interface for transformation (now also handles function form)
}

// WithDefaultPaths provides a default path map for components.
// This map is used as a base and can be overridden by paths defined in the bootstrap file.
func WithDefaultPaths(paths map[string]string) Option {
	return optionutil.WithUpdate(func(o *ConfigLoadOptions) {
		o.defaultPaths = paths
	})
}

// WithConfig allows providing a custom interfaces.Config implementation.
// WithCond this option is used, LoadConfig will return the provided config directly,
// bypassing the default Kratos config creation and file loading.
func WithConfig(cfg interfaces.Config) Option {
	return optionutil.WithUpdate(func(o *ConfigLoadOptions) {
		o.config = cfg
	})
}

// WithConfigTransformer allows providing an object that implements the ConfigTransformer interface,
// or a function of type ConfigTransformFunc.
// This provides a flexible way to customize the creation of interfaces.Config from kratosconfig.Config.
func WithConfigTransformer(transformer ConfigTransformer) Option {
	return optionutil.WithUpdate(func(o *ConfigLoadOptions) {
		o.configTransformer = transformer
	})
}

func FromConfigLoadOptions(opts ...Option) *ConfigLoadOptions {
	var bootstrapOpt ConfigLoadOptions
	optionutil.Apply(&bootstrapOpt, opts...)
	return &bootstrapOpt
}

// ComponentFactory is a function that creates a component instance.
// It receives the configuration decoder and the component provider, allowing for
// configuration-driven instantiation and dependency resolution.
// This is an alias for interfaces.ComponentFactoryFunc for consistency.
type ComponentFactory = interfaces.ComponentFactory

// --- Options for New ---

// Options holds configuration for the New function.
type Options struct {
	appInfo            *interfaces.AppInfo
	componentFactories map[string]ComponentFactory
}

// Option configures the New function.
type Option = options.Option

// WithAppInfo provides the application's metadata to the provider.
// This is a required option for New.
func WithAppInfo(info *interfaces.AppInfo) Option {
	return optionutil.WithUpdate(func(o *Options) {
		o.appInfo = info
	})
}

// WithComponent registers a component factory to be used during bootstrap.
func WithComponent(key string, factory ComponentFactory) Option {
	return optionutil.WithUpdate(func(o *Options) {
		if o.componentFactories == nil {
			o.componentFactories = make(map[string]ComponentFactory)
		}
		o.componentFactories[key] = factory
	})
}

func FromOptions(opts ...Option) *Options {
	var bootstrapOpt Options
	optionutil.Apply(&bootstrapOpt, opts...)
	return &bootstrapOpt
}
