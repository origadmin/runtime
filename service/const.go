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
)

type (
	// RegisterGRPCServerFunc register a gRPC server
	RegisterGRPCServerFunc = func(context.Context, *GRPCServer)
	// RegisterHTTPServerFunc register a HTTP server
	RegisterHTTPServerFunc = func(context.Context, *HTTPServer)
	// RegisterGRPCClientFunc register a gRPC client
	RegisterGRPCClientFunc = func(context.Context, *GRPCClient)
	// RegisterHTTPClientFunc register a HTTP client
	RegisterHTTPClientFunc = func(context.Context, *HTTPClient)
)

type HTTPServerRegister interface {
	HTTPServer(context.Context, *HTTPServer)
}

type GRPCServerRegister interface {
	GRPCServer(context.Context, *GRPCServer)
}

type ServerRegister interface {
	GRPCServerRegister
	HTTPServerRegister
	Server(context.Context, *GRPCServer, *HTTPServer)
}

var (
	ErrServiceNotFound = errors.New("service not found")
)

type httpCtx struct{}

func NewHTTPContext(ctx context.Context, c transhttp.Context) context.Context {
	return context.WithValue(ctx, httpCtx{}, c)
}

func FromHTTPContext(ctx context.Context) (transhttp.Context, bool) {
	v, ok := ctx.Value(httpCtx{}).(transhttp.Context)
	return v, ok
}
