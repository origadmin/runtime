package middleware

import (
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/optionutil"
)

// Options holds common options that have been resolved once at the top level.
// These options are then passed down to individual middleware factories.
type Options struct {
	Logger log.Logger // The resolved logger instance.
}

// FromOptions resolves common options from a slice of generic options.Option.
// It returns the resolved options.Context and a custom *Options struct.
func FromOptions(opts ...options.Option) (options.Context, *Options) {
	ctx, mwOpts := optionutil.ApplyNew[Options](opts...)
	if mwOpts == nil {
		mwOpts = &Options{}
	}
	if mwOpts.Logger == nil {
		mwOpts.Logger = log.FromContext(ctx)
	}
	return ctx, mwOpts
}

// withOptions is a helper to wrap a *middleware.Options struct into a generic options.Option.
// This is used when passing the resolved common options down to individual factories.
func withOptions(mwOpts *Options) options.Option {
	return func(c options.Context) options.Context {
		return c.With(optionutil.Key[Options]{}, mwOpts)
	}
}

func fromOptions(opts ...options.Option) *Options {
	ctx := optionutil.ApplyContext(optionutil.Empty(), opts...)
	mwOpts, ok := ctx.Value(optionutil.Key[Options]{}).(*Options)
	if !ok {
		mwOpts = &Options{}
	}
	if mwOpts.Logger == nil {
		mwOpts.Logger = log.FromContext(ctx)
	}
	return mwOpts
}
