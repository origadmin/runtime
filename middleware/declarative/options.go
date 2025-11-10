// Package declarative implements the functions, types, and interfaces for the module.
package declarative

import (
	"github.com/origadmin/runtime/extension/optionutil"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/interfaces/security/declarative"
)

type Options struct {
	policyProvider declarative.PolicyProvider
	defaultPolicy  string
}

// WithPolicyProvider sets the PolicyProvider for the middleware.
func WithPolicyProvider(pm declarative.PolicyProvider) options.Option {
	return optionutil.Update(func(o *Options) {
		o.policyProvider = pm
	})
}

// WithDefaultPolicy sets the default policy name to use if no policy is specified in the metadata.
func WithDefaultPolicy(policyName string) options.Option {
	return optionutil.Update(func(o *Options) {
		o.defaultPolicy = policyName
	})
}

func FromOptions(opts []options.Option) *Options {
	o := optionutil.NewT[Options](opts...)
	return o
}
