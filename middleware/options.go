package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/selector"

	"github.com/origadmin/runtime/extension/optionutil"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/log"
)

// Options holds common options that have been resolved once at the top level.
// These options are then passed down to individual middleware factories.
type Options struct {
	Logger    log.Logger
	MatchFunc selector.MatchFunc // MatchFunc for selector middleware
	Carrier   *Carrier
	Options   []Option
}

type Option = options.Option

func WithMatchFunc(matchFunc selector.MatchFunc) Option {
	return optionutil.Update(func(o *Options) {
		o.MatchFunc = matchFunc
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
	mwOpts := optionutil.NewT[Options](opts...)
	mwOpts.Logger = log.FromOptions(opts)
	if mwOpts.Carrier == nil {
		// WithCond the carrier is not set, use a new Carrier instance
		mwOpts.Carrier = &Carrier{
			Clients: make(map[string]KMiddleware),
			Servers: make(map[string]KMiddleware),
		}
	}
	mwOpts.Options = opts
	return mwOpts
}
