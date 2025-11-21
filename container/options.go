// Package container implements the functions, types, and interfaces for the module.
package container

import (
	"github.com/origadmin/runtime/extension/optionutil"
	"github.com/origadmin/runtime/interfaces/options"
)

type containerOptions struct {
	componentFactories map[string]ComponentFactory
}

// WithComponentFactory is an option to register a ComponentFactory with the container.
// It returns a function that applies the factory to the containerImpl.
func WithComponentFactory(name string, factory ComponentFactory) options.Option {
	return optionutil.Update(func(c *containerOptions) {
		if c.componentFactories == nil {
			c.componentFactories = make(map[string]ComponentFactory)
		}
		c.componentFactories[name] = factory
	})
}

func ComponentFactoryFromOptions(opts ...options.Option) map[string]ComponentFactory {
	o := optionutil.NewT[containerOptions](opts...)
	return o.componentFactories
}
