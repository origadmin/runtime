// Package container implements the functions, types, and interfaces for the module.
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

type Option = options.Option

// WithAppInfo sets the application's metadata for the container.
// This is the correct way to inject the definitive AppInfo.
func WithAppInfo(info interfaces.AppInfo) Option {
	return optionutil.Update(func(o *containerOptions) {
		o.appInfo = info
	})
}

// WithComponentFactory is an option to register a ComponentFactory with the container.
// It returns a function that applies the factory to the containerImpl.
func WithComponentFactory(name string, factory ComponentFactory) Option {
	return optionutil.Update(func(c *containerOptions) {
		if c.componentFactories == nil {
			c.componentFactories = make(map[string]ComponentFactory)
		}
		c.componentFactories[name] = factory
	})
}

func fromOptions(opts []options.Option) *containerOptions {
	return optionutil.NewT[containerOptions](opts...)
}
