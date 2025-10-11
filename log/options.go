// Package log implements the functions, types, and interfaces for the module.
package log

import (
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/optionutil"
)

type loggerContext struct {
	Logger Logger
}

func WithLogger(logger Logger) options.Option {
	return optionutil.Update(func(l *loggerContext) {
		l.Logger = logger
	})
}

func FromOptions(opts ...options.Option) Logger {
	_, l := optionutil.New[loggerContext](opts...)
	if l.Logger == nil {
		l.Logger = DefaultLogger
	}
	return l.Logger
}

func FromContext(ctx options.Context) Logger {
	v := optionutil.ValueCond(ctx, func(l *loggerContext) bool { return l != nil }, &loggerContext{
		Logger: DefaultLogger,
	})
	return v.Logger
}
