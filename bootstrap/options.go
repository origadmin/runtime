package bootstrap

import (
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/interfaces"
)

// --- Options for NewDecoder ---

// DecoderOption configures the NewDecoder function.
type DecoderOption func(*decoderOptions)

// decoderOptions holds configuration for the NewDecoder function.
type decoderOptions struct {
	defaultPaths  map[string]string
	kratosOptions []runtimeconfig.Option
}

// WithDefaultPaths provides a default path map for components.
// This map is used as a base and can be overridden by paths defined in the bootstrap file.
func WithDefaultPaths(paths map[string]string) DecoderOption {
	return func(o *decoderOptions) {
		o.defaultPaths = paths
	}
}

// WithKratosOption passes Kratos-specific config options to the underlying config creation.
func WithKratosOption(opts ...runtimeconfig.Option) DecoderOption {
	return func(o *decoderOptions) {
		o.kratosOptions = append(o.kratosOptions, opts...)
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
	appInfo               interfaces.AppInfo // Modified: Now holds an interfaces.AppInfo
	decoderOptions        []DecoderOption
	componentsToConfigure []configurableComponent
}

// WithAppInfo provides the application's metadata to the provider.
// This is a required option for NewProvider.
// Modified to accept an interfaces.AppInfo directly.
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
