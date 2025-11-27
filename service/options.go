package service

import (
	"context" // 导入标准库的 context，因为 runtime/context 可能是它的别名或封装

	"github.com/origadmin/runtime/extensions/optionutil"
	"github.com/origadmin/runtime/interfaces/options"
)

type serviceOptions struct {
	registrar []ServerRegistrar
	ctx       context.Context // 添加 context.Context 字段
}

func WithServerRegistrar(registrar ...ServerRegistrar) options.Option {
	return optionutil.Update(func(o *serviceOptions) {
		o.registrar = registrar
	})
}

// WithContext sets the context.Context for the service.
func WithContext(ctx context.Context) options.Option {
	return optionutil.Update(func(o *serviceOptions) {
		o.ctx = ctx
	})
}

// WithContextRegistrar sets both the context.Context and ServerRegistrar for the service.
func WithContextRegistrar(ctx context.Context, registrar ...ServerRegistrar) options.Option {
	return optionutil.Update(func(o *serviceOptions) {
		o.ctx = ctx
		o.registrar = registrar
	})
}

func ServerRegistrarFromOptions(opts []options.Option) []ServerRegistrar {
	o := fromOptions(opts)
	return o.registrar
}

func ServerRegistrarFromContext(ctx options.Context) []ServerRegistrar {
	v := optionutil.ValueCond(ctx, func(l *serviceOptions) bool { return l != nil && l.registrar != nil }, &serviceOptions{})
	return v.registrar
}

// ContextFromOptions retrieves the context.Context from a slice of options.
// If no context is found, context.Background() is returned.
func ContextFromOptions(opts []options.Option) context.Context {
	o := fromOptions(opts)
	return o.ctx
}

func fromOptions(opts []options.Option) *serviceOptions {
	o := optionutil.NewT[serviceOptions](opts...)
	if o.ctx == nil {
		o.ctx = context.Background()
	}
	return o
}
