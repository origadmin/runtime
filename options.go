package runtime

import (
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/extensions/optionutil"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
)

// appOptions holds the configurable settings for a App.
type appOptions struct {
	appInfo         interfaces.AppInfo
	bootstrapOpts   []options.Option
	containerOpts   []options.Option
	kratosAppOpts   []options.Option
	structuredCfg   interfaces.StructuredConfig
	config          interfaces.Config
	bootstrapResult bootstrap.Result
}

type Option = options.Option

// WithAppInfo sets the application's metadata.
func WithAppInfo(info interfaces.AppInfo) Option {
	return optionutil.Update(func(o *appOptions) {
		o.appInfo = info
	})
}

// WithContainerOptions applies options to the underlying container.
func WithContainerOptions(opts ...options.Option) Option {
	return optionutil.Update(func(o *appOptions) {
		o.containerOpts = append(o.containerOpts, opts...)
	})
}

// WithBootstrapOptions applies options to the underlying bootstrap process.
func WithBootstrapOptions(opts ...options.Option) Option {
	return optionutil.Update(func(o *appOptions) {
		o.bootstrapOpts = append(o.bootstrapOpts, opts...)
	})
}
