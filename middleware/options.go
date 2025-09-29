package middleware

import (
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/optionutil"
)

// Options holds common options that have been resolved once at the top level.
// These options are then passed down to individual middleware factories.
type Options struct {
	Logger    log.Logger      // The resolved logger instance.
	Context   options.Context // The resolved context instance.
	MatchFunc options.MatchFunc // MatchFunc for selector middleware
}

// FromOptions resolves common options from a slice of generic options.Option.
// It returns the resolved options.Context and a custom *Options struct.
func FromOptions(opts ...options.Option) (options.Context, *Options) {
	// Use optionutil to resolve the context and matchFunc
	// We need a temporary struct to hold these, as optionutil.ApplyNew works on a zero-value struct.
	type tempCommonOptions struct {
		Context   options.Context
		MatchFunc options.MatchFunc
	}

	ctx, commonOpts := optionutil.ApplyNew[tempCommonOptions](opts...)

	// Resolve the logger using log.FromOptions
	resolvedLogger := log.FromOptions(opts...)

	// Create the middleware.Options struct
	mwOpts := &Options{
		Logger:    resolvedLogger,
		Context:   ctx, // Use the resolved context
		MatchFunc: commonOpts.MatchFunc,
	}
	return ctx, mwOpts
}

// withOptions is a helper to wrap a *middleware.Options struct into a generic options.Option.
// This is used when passing the resolved common options down to individual factories.
func withOptions(mwOpts *Options) options.Option {
	return optionutil.Update(func(o *Options) {
		// This updates a middleware.Options struct. This is used when a factory
		// needs to receive a middleware.Options struct via functional options.
		*o = *mwOpts // Copy the resolved options
	})
}
