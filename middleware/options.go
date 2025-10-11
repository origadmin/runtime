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
	Logger    log.Logger         // The resolved logger instance.
	MatchFunc selector.MatchFunc // MatchFunc for selector middleware
	Carrier   *Carrier
}

type Option = options.Option

func WithMatchFunc(matchFunc selector.MatchFunc) Option {
	return optionutil.Update(func(o *Options) {
		o.MatchFunc = matchFunc
	})
}

func WithLogger(logger log.Logger) Option {
	return optionutil.Update(func(o *Options) {
		o.Logger = logger
	})
}

func WithCarrier(carrier *Carrier) Option {
	return optionutil.Update(func(o *Options) {
		o.Carrier = carrier
	})
}

// FromOptions resolves common options from a slice of generic Option.
// It returns the resolved options.Context and a custom *Options struct.
func FromOptions(opts ...Option) *Options {
	// Use optionutil to resolve the context and matchFunc
	// We need a temporary struct to hold these, as optionutil.New works on a zero-value struct.
	ctx, mwOpts := optionutil.New[Options](opts...)
	if mwOpts.Logger == nil {
		// WithCond the logger is not set, use the default logger
		mwOpts.Logger = log.FromContext(ctx)
	}
	return mwOpts
}
