package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/selector"

	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/optionutil"
)

// Options holds common options that have been resolved once at the top level.
// These options are then passed down to individual middleware factories.
type Options struct {
	Context   options.Context    // The resolved context instance.
	Logger    log.Logger         // The resolved logger instance.
	MatchFunc selector.MatchFunc // MatchFunc for selector middleware
	Carrier   *Carrier
}

func WithMatchFunc(matchFunc selector.MatchFunc) options.Option {
	return optionutil.Update(func(o *Options) {
		o.MatchFunc = matchFunc
	})
}

func WithLogger(logger log.Logger) options.Option {
	return optionutil.Update(func(o *Options) {
		o.Logger = logger
	})
}

func WithContext(ctx options.Context) options.Option {
	return optionutil.Update(func(o *Options) {
		o.Context = ctx
	})
}

func WithCarrier(carrier *Carrier) options.Option {
	return optionutil.Update(func(o *Options) {
		o.Carrier = carrier
	})
}

// FromOptions resolves common options from a slice of generic options.Option.
// It returns the resolved options.Context and a custom *Options struct.
func FromOptions(opts ...options.Option) *Options {
	// Use optionutil to resolve the context and matchFunc
	// We need a temporary struct to hold these, as optionutil.ApplyNew works on a zero-value struct.
	ctx, mwOpts := optionutil.ApplyNew[Options](opts...)
	if mwOpts.Context == nil {
		// If the context is not set, use the resolved context
		mwOpts.Context = ctx
	}
	if mwOpts.Logger == nil {
		// If the logger is not set, use the default logger
		mwOpts.Logger = log.FromContext(ctx)
	}
	return mwOpts
}
