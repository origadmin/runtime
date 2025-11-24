// Package service implements the functions, types, and interfaces for the module.
package service

import (
	"github.com/origadmin/runtime/extensions/optionutil"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/service/transport"
)

type serviceOptions struct {
	registrar transport.ServerRegistrar
}

func WithRegistrar(registrar transport.ServerRegistrar) options.Option {
	return optionutil.Update(func(o *serviceOptions) {
		o.registrar = registrar
	})
}

func ServerRegistrarFromOptions(opts []options.Option) transport.ServerRegistrar {
	o := optionutil.NewT[serviceOptions](opts...)
	return o.registrar
}

func ServerRegistrarFromContext(ctx options.Context) transport.ServerRegistrar {
	v := optionutil.ValueCond(ctx, func(l *serviceOptions) bool { return l != nil && l.registrar != nil }, &serviceOptions{})
	return v.registrar
}
