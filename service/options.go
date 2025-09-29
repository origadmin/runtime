package service

import (
	"github.com/go-kratos/kratos/v2/selector" // Import Kratos selector interface

	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/optionutil"
)

// ServiceOptions is a container for all service-level options.
// It is configured via the options pattern and is intended to be used by transport factories.
type ServiceOptions struct {
	Registrar            ServerRegistrar
	ClientEndpoint       string
	ClientSelectorFilter selector.NodeFilter
	MiddlewareProvider   interfaces.MiddlewareProvider
}

// WithRegistrar sets the ServerRegistrar for the service.
func WithRegistrar(r ServerRegistrar) interfaces.Option {
	return optionutil.Update(func(o *ServiceOptions) {
		o.Registrar = r
	})
}

// WithClientEndpoint sets the client's target endpoint (e.g., "discovery:///service-name").
func WithClientEndpoint(endpoint string) interfaces.Option {
	return optionutil.Update(func(o *ServiceOptions) {
		o.ClientEndpoint = endpoint
	})
}

// WithClientSelectorFilter sets the client's node filter for load balancing.
func WithClientSelectorFilter(filter selector.NodeFilter) interfaces.Option {
	return optionutil.Update(func(o *ServiceOptions) {
		o.ClientSelectorFilter = filter
	})
}

// WithMiddlewareProvider sets the MiddlewareProvider for the service.
func WithMiddlewareProvider(provider interfaces.MiddlewareProvider) interfaces.Option {
	return optionutil.Update(func(o *ServiceOptions) {
		o.MiddlewareProvider = provider
	})
}
