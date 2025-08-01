/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package service implements the functions, types, and interfaces for the module.
package service

import (
	"context"
	"errors"
	"time"

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"
)

const DefaultTimeout = 5 * time.Second

type (
	// GRPCServer define the gRPC server interface
	GRPCServer = transgrpc.Server
	// HTTPServer define the HTTP server interface
	HTTPServer = transhttp.Server
	// GRPCClient define the gRPC client interface
	GRPCClient = grpc.ClientConn
	// HTTPClient define the HTTP client interface
	HTTPClient = transhttp.Client
)

type (
	// GRPCServerOption define the gRPC server options
	GRPCServerOption = transgrpc.ServerOption
	// HTTPServerOption define the HTTP server options
	HTTPServerOption = transhttp.ServerOption
	// GRPCClientOption define the gRPC client options
	GRPCClientOption = transgrpc.ClientOption
	// HTTPClientOption define the HTTP client options
	HTTPClientOption = transhttp.ClientOption

	Option interface {
		GRPCServerOption | HTTPServerOption | GRPCClientOption | HTTPClientOption
	}
)

type (
	// RegisterFunc register a service
	RegisterFunc[T any] func(context.Context, *T)
	// RegisterGRPCFunc register a gRPC server
	RegisterGRPCFunc = RegisterFunc[GRPCServer]
	// RegisterHTTPFunc register a HTTP server
	RegisterHTTPFunc = RegisterFunc[HTTPServer]
)

func (f RegisterFunc[T]) Register(ctx context.Context, t *T) {
	f(ctx, t)
}

type Registrar[T any] interface {
	Register(context.Context, *T)
}

type HTTPRegistrar interface {
	RegisterHTTP(context.Context, *HTTPServer)
}

type GRPCRegistrar interface {
	RegisterGRPC(context.Context, *GRPCServer)
}

type ServerRegistrar interface {
	Register(context.Context, any)
	RegisterHTTP(context.Context, *HTTPServer)
	RegisterGRPC(context.Context, *GRPCServer)
}

var (
	ErrServiceNotFound = errors.New("service not found")
)
