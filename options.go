package runtime

import (
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
)

// appOptions holds the configurable parameters for the runtime App.
// It is unexported to keep it as an internal detail of the Option pattern.
type appOptions struct {
	bootstrapOpts []bootstrap.Option
	appInfo       interfaces.AppInfo // Now an interface
}

// Option defines a function that configures the runtime App.
type Option func(*appOptions)

// WithBootstrapOptions passes bootstrap-specific options to the underlying bootstrap process.
func WithBootstrapOptions(opts ...bootstrap.Option) Option {
	return func(o *appOptions) {
		o.bootstrapOpts = append(o.bootstrapOpts, opts...)
	}
}

// WithAppInfo provides application metadata programmatically by accepting an
// implementation of the interfaces.AppInfo interface.
func WithAppInfo(info interfaces.AppInfo) Option {
	return func(o *appOptions) {
		o.appInfo = info
	}
}
