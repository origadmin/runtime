package bootstrap

import (
	appv1 "github.com/origadmin/runtime/api/gen/go/app/v1"
	runtimeconfig "github.com/origadmin/runtime/config"
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

// options holds configuration for the NewProvider function.
type options struct {
	appInfo        *appv1.AppInfo
	decoderOptions []DecoderOption
}

// WithAppInfo provides the application's metadata to the provider.
// This is a required option for NewProvider.
func WithAppInfo(info *appv1.AppInfo) Option {
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
