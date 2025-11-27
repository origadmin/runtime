package runtime

import (
	"time"

	"github.com/origadmin/runtime/interfaces/options"
)

// Option is a functional option for configuring the App.
// It allows for applying configurations to the App instance at creation time.
type Option func(*App)

// AppInfoOption is a functional option for configuring the AppInfo.
// It allows for applying configurations to the appInfo instance at creation time.
type AppInfoOption func(a *appInfo)

// WithContainerOptions adds options that will be applied to the dependency injection container.
// These options are collected during the New() phase and applied during the Load() phase
// when the container is created.
func WithContainerOptions(opts ...options.Option) Option {
	return func(a *App) {
		a.containerOpts = append(a.containerOpts, opts...)
	}
}

// WithEnv sets the environment for the application.
func WithEnv(env string) Option {
	return func(a *App) {
		a.appInfo.env = env
	}
}

// WithID sets a custom instance ID.
func WithID(id string) Option {
	return func(a *App) {
		a.appInfo.id = id
	}
}

// WithStartTime sets a custom start time.
func WithStartTime(startTime time.Time) Option {
	return func(a *App) {
		a.appInfo.startTime = startTime
	}
}

// WithMetadata adds a key-value pair to the application's metadata.
func WithMetadata(key, value string) Option {
	return func(a *App) {
		if a.appInfo.metadata == nil {
			a.appInfo.metadata = make(map[string]string)
		}
		a.appInfo.metadata[key] = value
	}
}

// WithMetadataMap adds a map of key-value pairs to the application's metadata.
func WithMetadataMap(metadata map[string]string) Option {
	return func(a *App) {
		if a.appInfo.metadata == nil {
			a.appInfo.metadata = make(map[string]string)
		}
		for k, v := range metadata {
			a.appInfo.metadata[k] = v
		}
	}
}

// WithAppInfoOptions adds options that will be applied to the AppInfo.
func WithAppInfoOptions(opts ...AppInfoOption) Option {
	return func(a *App) {
		for _, opt := range opts {
			opt(a.appInfo)
		}
	}
}

// WithAppInfoID sets the application ID.
func WithAppInfoID(id string) AppInfoOption {
	return func(a *appInfo) {
		a.id = id
	}
}

// WithAppInfoName sets the application name.
func WithAppInfoName(name string) AppInfoOption {
	return func(a *appInfo) {
		a.name = name
	}
}

// WithAppInfoVersion sets the application version
func WithAppInfoVersion(version string) AppInfoOption {
	return func(a *appInfo) {
		a.version = version
	}
}

// WithAppInfoEnv sets the application environment.
func WithAppInfoEnv(env string) AppInfoOption {
	return func(a *appInfo) {
		a.env = env
	}
}

// WithAppInfoStartTime sets the application start time.
func WithAppInfoStartTime(startTime time.Time) AppInfoOption {
	return func(a *appInfo) {
		a.startTime = startTime
	}
}

// WithAppInfoMetadata adds a key-value pair to the application's metadata.
func WithAppInfoMetadata(key, value string) AppInfoOption {
	return func(a *appInfo) {
		if a.metadata == nil {
			a.metadata = make(map[string]string)
		}
		a.metadata[key] = value
	}
}

// WithAppInfoMetadataMap adds a map of key-value pairs
func WithAppInfoMetadataMap(metadata map[string]string) AppInfoOption {
	return func(a *appInfo) {
		if a.metadata == nil {
			a.metadata = make(map[string]string)
		}
		for k, v := range metadata {
			a.metadata[k] = v
		}
	}
}
