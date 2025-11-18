package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/golang-jwt/jwt/v5"

	"github.com/origadmin/runtime/extension/optionutil"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/log"
)

// Options holds common options that have been resolved once at the top level.
// These options are then passed down to individual middleware factory functions.
type Options struct {
	Logger         log.Logger
	MatchFunc      selector.MatchFunc
	Carrier        *Carrier
	Options        []Option
	ClaimsFactory  func() jwt.Claims
	SubjectFactory func() string
	SigningMethod  jwt.SigningMethod
}

// Option is a functional option for configuring middleware options.
type Option = options.Option

// WithMatchFunc sets the MatchFunc for the selector middleware.
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

// WithClaimsFactory provides a function that generates JWT claims.
// This is the recommended way to create dynamic claims for each token.
func WithClaimsFactory(factory func() jwt.Claims) Option {
	return optionutil.Update(func(o *Options) {
		o.ClaimsFactory = factory
	})
}

// WithSigningMethod sets the signing method to be used for JWT tokens.
// If not provided, the default from the configuration will be used.
func WithSigningMethod(method jwt.SigningMethod) Option {
	return optionutil.Update(func(o *Options) {
		o.SigningMethod = method
	})
}

// WithSubjectFactory provides a function that generates the JWT 'subject' (sub) claim.
// This is the recommended way to provide a meaningful user identifier for the token.
func WithSubjectFactory(factory func() string) Option {
	return optionutil.Update(func(o *Options) {
		o.SubjectFactory = factory
	})
}

// FromOptions resolves common options from a slice of generic Option.
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
