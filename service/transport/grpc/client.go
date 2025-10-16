package grpc

import (
	"context"

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"

	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
	runtimeerrors "github.com/origadmin/runtime/errors"
)

// NewClient creates a new concrete gRPC client connection based on the provided configuration.
// It returns *transgrpc.ClientConn.
func NewClient(ctx context.Context, grpcConfig *transportv1.GrpcClientConfig, clientOpts *ClientOptions) (*grpc.ClientConn, error) {
	// Initialize the Kratos gRPC client options using the adapter function.
	kratosOpts, err := initGrpcClientOptions(ctx, grpcConfig, clientOpts)
	if err != nil {
		return nil, err
	}

	// Configure TLS.
	useSecureDial := grpcConfig.GetTlsConfig().GetEnabled()

	// Create the Kratos gRPC client connection.
	var conn *grpc.ClientConn
	if useSecureDial {
		conn, err = transgrpc.Dial(ctx, kratosOpts...)
	} else {
		conn, err = transgrpc.DialInsecure(ctx, kratosOpts...)
	}

	if err != nil {
		return nil, runtimeerrors.WrapStructured(err, Module, "failed to create gRPC client")
	}

	return conn, nil
}
