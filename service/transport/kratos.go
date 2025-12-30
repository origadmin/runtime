/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package transport implements the functions, types, and interfaces for the module.
package transport

import (
	"context"
	"time"

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"
)

//go:generate adptool .
//go:adapter:package github.com/go-kratos/kratos/v2/transport/grpc transgrpc
//go:adapter:ignore NewWrappedStream
//go:adapter:package:type *
//go:adapter:package:type:suffix GRPC
//go:adapter:package:func *
//go:adapter:package:func:suffix GRPC
////go:adapter:package google.golang.org/grpc grpc
////go:adapter:package:type *
////go:adapter:package:type:suffix GRPC
////go:adapter:package:func *
////go:adapter:package:func:suffix GRPC
//go:adapter:package github.com/go-kratos/kratos/v2/transport/http transhttp
//go:adapter:package:type *
//go:adapter:package:type:suffix HTTP
//go:adapter:package:func *
//go:adapter:package:func:suffix HTTP
//go:adapter:package github.com/go-kratos/kratos/v2/transport transport

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

// GRPCRegistrar is a capability interface for services that can register gRPC endpoints.
type GRPCRegistrar interface {
	RegisterGRPC(ctx context.Context, srv *GRPCServer) error
}

// GRPCRegisterFunc is a function that implements GRPCRegistrar.
type GRPCRegisterFunc func(ctx context.Context, srv *GRPCServer) error

// RegisterGRPC implements the GRPCRegistrar interface for GRPCRegisterFunc.
func (f GRPCRegisterFunc) RegisterGRPC(ctx context.Context, srv *GRPCServer) error {
	return f(ctx, srv)
}

// Register implements the ServerRegistrar interface for GRPCRegisterFunc.
func (f GRPCRegisterFunc) Register(ctx context.Context, srv any) error {
	if srv, ok := srv.(*GRPCServer); ok {
		return f(ctx, srv)
	}
	return nil
}

// HTTPRegistrar is a capability interface for services that can register HTTP endpoints.
type HTTPRegistrar interface {
	RegisterHTTP(ctx context.Context, srv *HTTPServer) error
}

// HTTPRegisterFunc is a function that implements HTTPRegistrar.
type HTTPRegisterFunc func(ctx context.Context, srv *HTTPServer) error

// RegisterHTTP implements the HTTPRegistrar interface for HTTPRegisterFunc.
func (f HTTPRegisterFunc) RegisterHTTP(ctx context.Context, srv *HTTPServer) error {
	return f(ctx, srv)
}

// Register implements the ServerRegistrar interface for HTTPRegisterFunc.
func (f HTTPRegisterFunc) Register(ctx context.Context, srv any) error {
	if srv, ok := srv.(*HTTPServer); ok {
		return f(ctx, srv)
	}
	return nil
}
