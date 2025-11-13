// Package declarative implements the functions, types, and interfaces for the module.
package declarative

import (
	"github.com/origadmin/runtime/extension/optionutil"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/interfaces/security/declarative"
)

// Options holds the configuration for the declarative security middleware.
type Options struct {
	policyProvider      declarative.PolicyProvider
	credentialExtractor declarative.CredentialExtractor // Added CredentialExtractor
	defaultPolicy       string
}

// WithPolicyProvider sets the PolicyProvider for the middleware.
func WithPolicyProvider(pm declarative.PolicyProvider) options.Option {
	return optionutil.Update(func(o *Options) {
		o.policyProvider = pm
	})
}

// WithCredentialExtractor sets the CredentialExtractor for the middleware.
func WithCredentialExtractor(ce declarative.CredentialExtractor) options.Option {
	return optionutil.Update(func(o *Options) {
		o.credentialExtractor = ce
	})
}

// WithDefaultPolicy sets the default policy name to use if no policy is specified in the metadata.
func WithDefaultPolicy(policyName string) options.Option {
	return optionutil.Update(func(o *Options) {
		o.defaultPolicy = policyName
	})
}

// FromOptions converts a slice of options.Option to a declarative.Options struct.
func FromOptions(opts []options.Option) *Options {
	o := optionutil.NewT[Options](opts...)
	return o
}
