package registry

import (
	"github.com/origadmin/runtime/extensions/optionutil"
	"github.com/origadmin/runtime/interfaces/options"
)

// options is a container for registry-related options.
type registryOptions struct {
	Registrar KRegistrar
}

// WithRegistrar sets the ServerRegistrar for the service.
func WithRegistrar(r KRegistrar) options.Option {
	return optionutil.Update(func(o *registryOptions) {
		o.Registrar = r
	})
}

// FromOptions creates a new Options struct by applying a slice of functional options.
func FromOptions(opts ...options.Option) KRegistrar {
	o := optionutil.NewT[registryOptions](opts...)
	return o.Registrar
}
