package runtime

import (
	"github.com/origadmin/runtime/interfaces/options"
)

// Option is a functional option for configuring the App.
// It allows for applying configurations to the App instance at creation time.
type Option func(*App)

// WithContainerOptions adds options that will be applied to the dependency injection container.
// These options are collected during the New() phase and applied during the Load() phase
// when the container is created.
func WithContainerOptions(opts ...options.Option) Option {
	return func(a *App) {
		a.containerOpts = append(a.containerOpts, opts...)
	}
}

func WithAppInfo(name, version string, opts ...AppInfoOption) Option {
	return func(a *App) {
		a.appInfo = NewAppInfo(name, version, opts...)
	}
}
