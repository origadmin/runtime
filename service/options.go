package service

import (
	"github.com/go-kratos/kratos/v2/selector" // Import Kratos selector interface

	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/optionutil"
)

type serviceOptions struct {
	registrar ServerRegistrar
	// New fields for client-specific options
	clientEndpoint       string
	clientSelectorFilter selector.NodeFilter
}

// WithRegistrar sets the ServerRegistrar for the service.
func WithRegistrar(r ServerRegistrar) interfaces.Option {
	return optionutil.Update(func(o *serviceOptions) {
		o.registrar = r
	})
}

// WithClientEndpoint sets the client's target endpoint (e.g., "discovery:///service-name").
func WithClientEndpoint(endpoint string) interfaces.Option {
	return optionutil.Update(func(o *serviceOptions) {
		o.clientEndpoint = endpoint
	})
}

// WithClientSelectorFilter sets the client's node filter for load balancing.
func WithClientSelectorFilter(filter selector.NodeFilter) interfaces.Option {
	return optionutil.Update(func(o *serviceOptions) {
		o.clientSelectorFilter = filter
	})
}
