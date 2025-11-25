package container

import (
	"github.com/origadmin/runtime/extensions/optionutil"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
)

// containerOptions holds the configurable settings for a Container.
type containerOptions struct {
	appInfo                interfaces.AppInfo
	componentFactories     map[string]ComponentFactory
	defaultCacheName       string
	defaultDatabaseName    string
	defaultObjectStoreName string
	defaultRegistrarName   string
}

// WithAppInfo sets the application's metadata for the container.
func WithAppInfo(info interfaces.AppInfo) options.Option {
	return optionutil.Update(func(o *containerOptions) {
		o.appInfo = info
	})
}

// WithComponentFactory registers a component factory with the container.
func WithComponentFactory(name string, factory ComponentFactory) options.Option {
	return optionutil.Update(func(o *containerOptions) {
		if o.componentFactories == nil {
			o.componentFactories = make(map[string]ComponentFactory)
		}
		o.componentFactories[name] = factory
	})
}

// WithDefaultCacheName sets the global default cache name.
func WithDefaultCacheName(name string) options.Option {
	return optionutil.Update(func(o *containerOptions) {
		o.defaultCacheName = name
	})
}

// WithDefaultDatabaseName sets the global default database name.
func WithDefaultDatabaseName(name string) options.Option {
	return optionutil.Update(func(o *containerOptions) {
		o.defaultDatabaseName = name
	})
}

// WithDefaultObjectStoreName sets the global default object store name.
func WithDefaultObjectStoreName(name string) options.Option {
	return optionutil.Update(func(o *containerOptions) {
		o.defaultObjectStoreName = name
	})
}

// WithDefaultRegistrarName sets the global default registrar name.
func WithDefaultRegistrarName(name string) options.Option {
	return optionutil.Update(func(o *containerOptions) {
		o.defaultRegistrarName = name
	})
}
