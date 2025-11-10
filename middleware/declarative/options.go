// Package declarative implements the functions, types, and interfaces for the module.
package declarative

import (
	"github.com/origadmin/runtime/extension/optionutil"
	"github.com/origadmin/runtime/interfaces/options"
)

type Options struct {
	policyManager PolicyManager
	defaultPolicy string
}

// WithPolicyManager sets the PolicyManager for the middleware.
func WithPolicyManager(pm PolicyManager) options.Option {
	return optionutil.Update(func(o *Options) {
		o.policyManager = pm
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
