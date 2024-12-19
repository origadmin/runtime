/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package service implements the functions, types, and interfaces for the module.
package service

import (
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
	// RegisterGRPCServer register a gRPC server
	RegisterGRPCServer = func(s *GRPCServer)
	// RegisterHTTPServer register a HTTP server
	RegisterHTTPServer = func(s *HTTPServer)
	// RegisterGRPCClient register a gRPC client
	RegisterGRPCClient = func(c *GRPCClient)
	// RegisterHTTPClient register a HTTP client
	RegisterHTTPClient = func(c *HTTPClient)
)

var (
	ErrServiceNotFound = errors.New("service not found")
)
