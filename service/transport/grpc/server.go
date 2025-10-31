package grpc

import (
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	grpcv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/grpc/v1"
)

// NewServer creates a new concrete gRPC server instance based on the provided configuration.
// It returns *transgrpc.Server, not the generic interfaces.Server.
func NewServer(grpcConfig *grpcv1.Server, serverOpts *ServerOptions) (*transgrpc.Server, error) {
	// Initialize the Kratos gRPC server options using the adapter function.
	kratosOpts, err := initGrpcServerOptions(grpcConfig, serverOpts)
	if err != nil {
		return nil, err
	}

	// Create the Kratos gRPC server instance.
	srv := transgrpc.NewServer(kratosOpts...)

	return srv, nil
}
