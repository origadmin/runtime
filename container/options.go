package container

import (
	"github.com/origadmin/runtime/extension/optionutil"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
)

// containerOptions holds the configurable settings for a Container.
type containerOptions struct {
	appInfo            interfaces.AppInfo
	componentFactories map[string]ComponentFactory
}

// WithAppInfo sets the application's metadata for the container.
func WithAppInfo(info interfaces.AppInfo) options.Option {
	return optionutil.Update(func(o *containerOptions) {
		o.appInfo = info
	})
}

func WithComponentFactory(name string, factory ComponentFactory) options.Option {
	return optionutil.Update(func(o *containerOptions) {
		o.componentFactories[name] = factory
	})
}
