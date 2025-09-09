package grpc

import (
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc" // Import Kratos gRPC transport

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/context" // Use project's context
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/service"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// NewClient creates a new gRPC client.
// It is the recommended way to create a client when the protocol is known in advance.
func NewClient(ctx context.Context, cfg *configv1.Service, opts ...service.Option) (interfaces.Client, error) {
	// 1. Get gRPC config and endpoint
	if cfg == nil || cfg.GetGrpc() == nil {
		return nil, tkerrors.Errorf("grpc config is required for client creation")
	}
	grpcCfg := cfg.GetGrpc()
	endpoint := grpcCfg.GetEndpoint()
	if endpoint == "" {
		return nil, tkerrors.Errorf("grpc endpoint is required for client creation")
	}

	// 2. Convert config to client options (these are []transgrpc.ClientOption)
	clientOptions, err := adaptClientConfig(cfg)
	if err != nil {
		return nil, tkerrors.Wrapf(err, "failed to adapt client config for grpc client creation")
	}

	// 3. Apply and extract options from context (these are []transgrpc.ClientOption)
	svcOpts := &service.Options{ContextOptions: interfaces.ContextOptions{Context: ctx}}
	for _, opt := range opts {
		opt(svcOpts)
	}
	if clientOptsFromCtx := FromClientOptions(svcOpts); len(clientOptsFromCtx) > 0 {
		clientOptions = append(clientOptions, clientOptsFromCtx...)
	}

	// 4. Create the underlying transport client using transgrpc.Dial
	// transgrpc.Dial expects transgrpc.ClientOption
	kratosClientOptions := []transgrpc.ClientOption{
		transgrpc.WithEndpoint(endpoint), // Endpoint is passed here
	}
	// Append the collected clientOptions (which are already transgrpc.ClientOption)
	kratosClientOptions = append(kratosClientOptions, clientOptions...)

	conn, err := transgrpc.Dial(ctx, kratosClientOptions...)
	if err != nil {
		return nil, tkerrors.Wrapf(err, "failed to create grpc client connection to %s", endpoint)
	}
	return conn, nil
}
