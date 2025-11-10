// Package service implements the functions, types, and interfaces for the module.
package service

import (
	"github.com/origadmin/runtime/extension/optionutil"
	"github.com/origadmin/runtime/interfaces/options"
)

type serviceOptions struct {
	registrar ServerRegistrar
}

func WithRegistrar(registrar ServerRegistrar) options.Option {
	return optionutil.Update(func(o *serviceOptions) {
		o.registrar = registrar
	})
}

func ServerRegistrarFromOptions(opts []options.Option) ServerRegistrar {
	o := optionutil.NewT[serviceOptions](opts...)
	return o.registrar
}

func ServerRegistrarFromContext(ctx options.Context) ServerRegistrar {
	v := optionutil.ValueCond(ctx, func(l *serviceOptions) bool { return l != nil && l.registrar != nil }, &serviceOptions{})
	return v.registrar
}
