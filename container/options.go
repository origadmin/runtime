package container

import (
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/optionutil"
)

// containerContext holds the application container instance.
type containerContext struct {
	Container interfaces.Container
}

// WithContainer creates an Option that sets the application container.
// This option is typically used during the initialization of components
// that require access to the dependency injection container.
func WithContainer(c interfaces.Container) options.Option {
	return optionutil.Update(func(cc *containerContext) {
		cc.Container = c
	})
}

// FromOptions extracts the application container from a list of options.
// It applies the provided options to a new containerContext and returns
// the contained interfaces.Container. If no container is set via options,
// it returns nil.
func FromOptions(opts []options.Option) interfaces.Container {
	cc := optionutil.NewT[containerContext](opts...)
	return cc.Container
}

// FromContext extracts the application container from an options.Context.
// This function is useful when the container needs to be retrieved from
// a context object, typically in scenarios where options are propagated
// through a context. If no container is found in the context, it returns nil.
func FromContext(ctx options.Context) interfaces.Container {
	v := optionutil.ValueCond(ctx, func(cc *containerContext) bool { return cc != nil && cc.Container != nil }, &containerContext{})
	return v.Container
}
