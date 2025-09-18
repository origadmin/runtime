/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package service implements the functions, types, and interfaces for the module.
package service

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
	RegisterGRPC(srv *GRPCServer)
}

// GRPCRegisterFunc is a function that implements GRPCRegistrar.
type GRPCRegisterFunc func(ctx context.Context, srv *GRPCServer)

// RegisterGRPC implements the GRPCRegistrar interface for GRPCRegisterFunc.
func (f GRPCRegisterFunc) RegisterGRPC(ctx context.Context, srv *GRPCServer) {
	f(ctx, srv)
}

// Register implements the ServerRegistrar interface for GRPCRegisterFunc.
func (f GRPCRegisterFunc) Register(ctx context.Context, srv any) {
	if srv, ok := srv.(*GRPCServer); ok {
		f(ctx, srv)
	}
}

// HTTPRegistrar is a capability interface for services that can register HTTP endpoints.
type HTTPRegistrar interface {
	RegisterHTTP(srv *HTTPServer)
}

// HTTPRegisterFunc is a function that implements HTTPRegistrar.
type HTTPRegisterFunc func(ctx context.Context, srv *HTTPServer)

// RegisterHTTP implements the HTTPRegistrar interface for HTTPRegisterFunc.
func (f HTTPRegisterFunc) RegisterHTTP(ctx context.Context, srv *HTTPServer) {
	f(ctx, srv)
}

// Register implements the ServerRegistrar interface for HTTPRegisterFunc.
func (f HTTPRegisterFunc) Register(ctx context.Context, srv any) {
	if srv, ok := srv.(*HTTPServer); ok {
		f(ctx, srv)
	}
}

// ServerRegistrar defines the single, universal entry point for service registration.
// It is the responsibility of the transport-specific factories to pass the correct server type
// (e.g., *GRPCServer or *HTTPServer) to the Register method.
// The implementation of this interface, typically within the user's project, is expected
// to perform a type switch to handle the specific server type.
type ServerRegistrar interface {
	Register(ctx context.Context, srv any)
}

// ServerRegisterFunc is a function that implements ServerRegistrar.
type ServerRegisterFunc func(ctx context.Context, srv any)

// Register implements the ServerRegistrar interface for ServerRegisterFunc.
func (f ServerRegisterFunc) Register(ctx context.Context, srv any) {
	f(ctx, srv)
}
