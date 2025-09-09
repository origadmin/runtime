package grpc

import (
	context "github.com/origadmin/runtime/context/adapter" // Use project's context adapter

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
	tkerrors "github.com/origadmin/toolkits/errors"
	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/service"
)

// NewClient creates a new gRPC client.
// It is the recommended way to create a client when the protocol is known in advance.
func NewClient(ctx context.Context, cfg *configv1.Service, opts ...service.Option) (interfaces.Client, error) {
	// 1. Convert config to client options
	clientOpts, err := adaptClientConfig(cfg)
	if err != nil {
		// adaptClientConfig 返回的是 tkerrors，这里继续返回 tkerrors
		return nil, tkerrors.Wrapf(err, "failed to adapt client config for grpc client creation") // 修正为 Wrapf
	}

	// 2. Apply and extract options from context
	svcOpts := &service.Options{ContextOptions: interfaces.ContextOptions{Context: ctx}}
	for _, opt := range opts {
		opt(svcOpts)
	}
	if clientOptsFromCtx := FromClientOptions(svcOpts); len(clientOptsFromCtx) > 0 {
		clientOpts = append(clientOpts, clientOptsFromCtx...)
	}

	// 3. Create the underlying transport client
	client, err := transgrpc.NewClient(ctx, clientOpts...)
	if err != nil {
		// transgrpc.NewClient 返回的是内部错误，这里继续返回 tkerrors
		return nil, tkerrors.Wrapf(err, "failed to create grpc client") // 修正为 Wrapf
	}
	return client, nil
}
