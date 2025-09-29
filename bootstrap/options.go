package bootstrap

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"

	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/interfaces"
)

// --- Options for NewDecoder ---

// ConfigTransformer defines an interface for custom transformation of kratosconfig.Config to interfaces.Config.
type ConfigTransformer interface {
	Transform(kratosconfig.Config) (interfaces.Config, error)
}

// ConfigTransformFunc is a function type that implements the ConfigTransformer interface.
type ConfigTransformFunc func(kratosconfig.Config) (interfaces.Config, error)

// Transform implements the ConfigTransformer interface for ConfigTransformFunc.
func (f ConfigTransformFunc) Transform(config kratosconfig.Config) (interfaces.Config, error) {
	return f(config)
}

// DecoderOption configures the NewDecoder function.
type DecoderOption func(*decoderOptions)

// decoderOptions holds configuration for the NewDecoder function.
type decoderOptions struct {
	defaultPaths      map[string]string
	configOptions     []runtimeconfig.Option // Changed from runtimeconfig.Empty to kratosconfig.Empty
	customConfig      interfaces.Config      // Added: Custom interfaces.Config implementation
	kratosConfig      kratosconfig.Config    // Added: Direct Kratos config instance
	configTransformer ConfigTransformer      // Custom interface for transformation (now also handles function form)
}

// WithDefaultPaths provides a default path map for components.
// This map is used as a base and can be overridden by paths defined in the bootstrap file.
func WithDefaultPaths(paths map[string]string) DecoderOption {
	return func(o *decoderOptions) {
		o.defaultPaths = paths
	}
}

// WithConfigOption passes Kratos-specific config options to the underlying config creation.
func WithConfigOption(opts ...runtimeconfig.Option) DecoderOption {
	return func(o *decoderOptions) {
		o.configOptions = append(o.configOptions, opts...)
	}
}

// WithCustomConfig allows providing a custom interfaces.Config implementation.
// If this option is used, NewDecoder will return the provided config directly,
// bypassing the default Kratos config creation and file loading.
func WithCustomConfig(cfg interfaces.Config) DecoderOption {
	return func(o *decoderOptions) {
		o.customConfig = cfg
	}
}

// WithKratosConfig allows providing a direct Kratos config.Config instance.
// If this option is used, NewDecoder will use the provided Kratos config directly,
// bypassing the default Kratos config creation and file loading from bootstrap.yaml sources.
func WithKratosConfig(kc kratosconfig.Config) DecoderOption {
	return func(o *decoderOptions) {
		o.kratosConfig = kc
	}
}

// WithConfigTransformer allows providing an object that implements the ConfigTransformer interface,
// or a function of type ConfigTransformFunc.
// This provides a flexible way to customize the creation of interfaces.Config from kratosconfig.Config.
func WithConfigTransformer(transformer ConfigTransformer) DecoderOption {
	return func(o *decoderOptions) {
		o.configTransformer = transformer
	}
}

// --- Options for NewProvider ---

// Option configures the NewProvider function.
type Option func(*options)

// configurableComponent holds the details for a user-defined component that needs to be decoded from config.
type configurableComponent struct {
	// Key is the top-level configuration key for this component (e.g., "custom_settings").
	Key string
	// Target is a pointer to the struct that the configuration should be decoded into.
	Target interface{}
}

// options holds configuration for the NewProvider function.
type options struct {
	appInfo               interfaces.AppInfo
	decoderOptions        []DecoderOption
	componentsToConfigure []configurableComponent
}

// Options contains the options for creating registry components.
// It embeds interfaces.ContextOptions for common context handling.
type Options = func(o *options)

// WithAppInfo provides the application's metadata to the provider.
// This is a required option for NewProvider.
func WithAppInfo(info interfaces.AppInfo) Option {
	return func(o *options) {
		o.appInfo = info
	}
}

// WithDecoderOptions allows passing DecoderOption functions to the internal NewDecoder call.
// This enables the caller of NewProvider to configure the decoding process.
func WithDecoderOptions(opts ...DecoderOption) Option {
	return func(o *options) {
		o.decoderOptions = append(o.decoderOptions, opts...)
	}
}

// WithComponent registers a custom component to be decoded from the configuration.
// The `key` specifies the top-level configuration key (e.g., "custom_settings").
// The `target` must be a pointer to a struct, which will be populated with the configuration values.
// After successful decoding, the populated struct will be available via `runtime.Component(key)`.
func WithComponent(key string, target interface{}) Option {
	return func(o *options) {
		o.componentsToConfigure = append(o.componentsToConfigure, configurableComponent{
			Key:    key,
			Target: target,
		})
	}
}

// WithOption is a generic way for any module to contribute its options.
// It takes a function that knows how to apply a module's specific options
// to a given options.Option and return the modified options.Option.
// This allows for type-safe module-specific configuration without bootstrap
// needing to know the internal types of each module's options.
func WithOption(opts ...func(options.Option)) Option {
	return func(o *options) {
		o.options = append(o.options, opts...)
	}
}
