// Package log implements the functions, types, and contracts for the module.
package log

import (
	"github.com/origadmin/runtime/helpers/optionutil"
	"github.com/origadmin/runtime/contracts/options"
)

type loggerContext struct {
	Logger Logger
}

func WithLogger(logger Logger) options.Option {
	return optionutil.Update(func(l *loggerContext) {
		l.Logger = logger
	})
}

func FromOptions(opts []options.Option) Logger {
	l := optionutil.NewT[loggerContext](opts...)
	if l.Logger == nil {
		l.Logger = DefaultLogger
	}
	return l.Logger
}

func FromContext(ctx options.Context) Logger {
	v := optionutil.ValueCond(ctx, func(l *loggerContext) bool { return l != nil && l.Logger != nil }, &loggerContext{
		Logger: DefaultLogger,
	})
	return v.Logger
}
