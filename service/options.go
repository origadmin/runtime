package service

import (
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/optionutil"
)

// Options is a container for all service-level options.
// It is configured via the options pattern and is intended to be used by transport factories.
// These options are typically provided during application bootstrap.
type Options struct {
	// Registrar is used to register server information with a service discovery registry.
	Registrar ServerRegistrar

	// Container provides access to various application components, including middleware.
	Container interfaces.Container
}

// FromOptions creates a new Options struct by applying a slice of functional options.
func FromOptions(opts []options.Option) *Options {
	o := &Options{}
	optionutil.Apply(o, opts...)
	return o
}

// WithRegistrar sets the ServerRegistrar for the service.
func WithRegistrar(r ServerRegistrar) options.Option {
	return optionutil.WithUpdate(func(o *Options) {
		o.Registrar = r
	})
}

// WithContainer sets the application component container.
func WithContainer(c interfaces.Container) options.Option {
	return optionutil.WithUpdate(func(o *Options) {
		o.Container = c
	})
}
