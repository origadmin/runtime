package service

import (
	"github.com/origadmin/runtime/optionutil"
)

type serviceOptions struct {
	registrar ServerRegistrar
}

// Options contains the options for creating service components.
// It embeds interfaces.OptionValue for common context handling.
type Options = optionutil.Options[serviceOptions]

// Option is a function that configures service.Options.
type Option func(*Options)

// WithRegistrar sets the ServerRegistrar for the service.
func WithRegistrar(r ServerRegistrar) Option {
	return func(o *Options) {
		o.Update(func(v *serviceOptions) {
			v.registrar = r
		})
	}
}

func RegistrarFromOptions(o *Options) ServerRegistrar {
	opts := o.Unwrap()
	if opts == nil {
		return nil
	}
	return opts.registrar
}
