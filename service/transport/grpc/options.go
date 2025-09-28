package grpc

import (
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/origadmin/runtime/optionutil"
	"github.com/origadmin/runtime/service"
)

// optionsKey is a private key type to avoid collisions in context.
var (
	serverOptionsKey = optionutil.Key[[]transgrpc.ServerOption]{}
	clientOptionsKey = optionutil.Key[[]transgrpc.ClientOption]{}
)

type grpcServerOptions struct {
	grpcServer []transgrpc.ServerOption
}

// WithServerOption is an option to add a Kratos transgrpc.ServerOption to the context.
func WithServerOption(opt ...transgrpc.ServerOption) service.Option {
	return func(o *service.Options) {
		optionutil.Append(o, serverOptionsKey, opt...)
	}
}

// WithClientOption is an option to add a transgrpc.ClientOption to the context.
func WithClientOption(opt ...transgrpc.ClientOption) service.Option { // Change parameter type
	return func(o *service.Options) {
		optionutil.Append(o, clientOptionsKey, opt...)
	}
}

// FromServerOptions returns the collected Kratos transgrpc.ServerOption from the service.Options' emptyContext.
func FromServerOptions(o *service.Options) []transgrpc.ServerOption {
	return optionutil.Slice(o, serverOptionsKey)
}

// FromClientOptions returns the collected transgrpc.ClientOption from the service.Options' emptyContext.
func FromClientOptions(o *service.Options) []transgrpc.ClientOption { // Change return type
	return optionutil.Slice(o, clientOptionsKey)
}
