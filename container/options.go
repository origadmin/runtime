package container

import (
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/helpers/optionutil"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/log"
)

// containerOptions holds the configurable settings for a Container.
type containerOptions struct {
	appInfo                *appv1.App
	middlewareConfig       *middlewarev1.Middlewares
	loggerConfig           *loggerv1.Logger
	componentFactories     map[string]ComponentFactory
	defaultCacheName       string
	defaultDatabaseName    string
	defaultObjectStoreName string
	defaultRegistrarName   string
	Logger                 log.Logger
}

// WithAppInfo sets the application's metadata for the container.
func WithAppInfo(info *appv1.App) options.Option {
	return optionutil.Update(func(o *containerOptions) {
		o.appInfo = info
	})
}

// WithMiddlewareConfig sets the middleware configuration for the container.
func WithMiddlewareConfig(cfg *middlewarev1.Middlewares) options.Option {
	return optionutil.Update(func(o *containerOptions) {
		o.middlewareConfig = cfg
	})
}

// WithLoggerConfig sets the logger configuration for the container.
func WithLoggerConfig(cfg *loggerv1.Logger) options.Option {
	return optionutil.Update(func(o *containerOptions) {
		o.loggerConfig = cfg
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

// WithLogger sets the logger for the container.
func WithLogger(logger log.Logger) options.Option {
	return log.WithLogger(logger)
}

type containerContext struct {
	Container Container
}

func WithContainer(container Container) options.Option {
	return optionutil.Update(func(o *containerContext) {
		o.Container = container
	})
}

func FromOptions(opts []options.Option) Container {
	l := optionutil.NewT[containerContext](opts...)
	if l.Container != nil {
		return l.Container
	}
	return nil
}
