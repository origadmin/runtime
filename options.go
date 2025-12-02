package runtime

import (
	"time"

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

// WithAppInfo sets the App's AppInfo using a concrete *AppInfo instance.
// This option has priority over the name and version parameters in New().
func WithAppInfo(info *AppInfo) Option {
	return func(a *App) {
		a.appInfo = info
	}
}

// WithEnv sets the environment for the application.
func WithEnv(env string) Option {
	return func(a *App) {
		// This option will be applied after the appInfo is initialized.
		a.appInfo.SetEnv(env)
	}
}

// WithID sets a custom instance ID.
func WithID(id string) Option {
	return func(a *App) {
		// This option will be applied after the appInfo is initialized.
		a.appInfo.SetID(id)
	}
}

// WithStartTime sets a custom start time.
func WithStartTime(startTime time.Time) Option {
	return func(a *App) {
		// This option will be applied after the appInfo is initialized.
		a.appInfo.SetStartTime(startTime)
	}
}

// WithMetadata adds a key-value pair to the application's metadata.
func WithMetadata(key, value string) Option {
	return func(a *App) {
		// This option will be applied after the appInfo is initialized.
		a.appInfo.AddMetadata(key, value)
	}
}

// WithMetadataMap adds a map of key-value pairs to the application's metadata.
func WithMetadataMap(metadata map[string]string) Option {
	return func(a *App) {
		// This option will be applied after the appInfo is initialized.
		a.appInfo.SetMetadata(metadata)
	}
}
