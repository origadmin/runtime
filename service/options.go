package service

import (
	"github.com/go-kratos/kratos/v2/selector" // Import Kratos selector interface

	"github.com/origadmin/runtime/optionutil"
)

type serviceOptions struct {
	registrar ServerRegistrar
	// New fields for client-specific options
	clientEndpoint       string
	clientSelectorFilter selector.NodeFilter
}

// Options contains the options for creating service components.
// It embeds interfaces.OptionValue for common context handling.
type Options = optionutil.Options[serviceOptions]

// Option is a function that configures service.Options.
type Option func(*Options)

// WithRegistrar sets the ServerRegistrar for the service.
func WithRegistrar(r ServerRegistrar) Option {
	return func(o *Options) {
		o.Update(func(so *serviceOptions) {
			so.registrar = r
		})
	}
}

// WithClientEndpoint sets the client's target endpoint (e.g., "discovery:///service-name").
func WithClientEndpoint(endpoint string) Option {
	return func(o *Options) {
		o.Update(func(so *serviceOptions) {
			so.clientEndpoint = endpoint
		})
	}
}

// WithClientSelectorFilter sets the client's node filter for load balancing.
func WithClientSelectorFilter(filter selector.NodeFilter) Option {
	return func(o *Options) {
		o.Update(func(so *serviceOptions) {
			so.clientSelectorFilter = filter
		})
	}
}
