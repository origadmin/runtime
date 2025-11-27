// Package service implements the functions, types, and interfaces for the module.
package transport

import (
	"github.com/origadmin/runtime/extensions/optionutil"
	"github.com/origadmin/runtime/interfaces/options"
)

type transportOptions struct {
	registrar ServerRegistrar
}

func WithServerRegistrar(registrar ServerRegistrar) options.Option {
	return optionutil.Update(func(o *transportOptions) {
		o.registrar = registrar
	})
}

func ServerRegistrarFromOptions(opts []options.Option) ServerRegistrar {
	o := optionutil.NewT[transportOptions](opts...)
	return o.registrar
}

func ServerRegistrarFromContext(ctx options.Context) ServerRegistrar {
	v := optionutil.ValueCond(ctx, func(l *transportOptions) bool { return l != nil && l.registrar != nil }, &transportOptions{})
	return v.registrar
}
