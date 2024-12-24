/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package agent implements the functions, types, and interfaces for the module.
package agent

import (
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"
)

// ApiVersionV1 defines the version number of the API
const ApiVersionV1 = "/api/v1"

// DefaultPrefix defines the default API prefix
const DefaultPrefix = "/api"

// DefaultVersion defines the default API version
const DefaultVersion = "v1"

// Agent is an interface that combines the HTTPAgent and GRPCAgent interfaces
type Agent interface {
	HTTPAgent
	GRPCAgent
}

// HTTPAgent is an interface that defines the basic methods of an HTTP proxy
type HTTPAgent interface {
	// URI returns the URI of the HTTP service
	URI() string
	// HTTPServer returns an instance of the HTTP server
	HTTPServer() *http.Server
	// Route returns an instance of the HTTP router
	Route() *http.Router
}

// GRPCAgent is an interface that defines the basic methods of a gRPC proxy
type GRPCAgent interface {
	// Server returns an instance of the gRPC server
	Server() *grpc.Server
	// RegisterService registers a gRPC service
	RegisterService(desc *grpc.ServiceDesc, impl interface{})
}

// agent is an implementation of the Agent interface
type agent struct {
	GRPCAgent
	HTTPAgent
}

// NewAgent creates a new Agent instance that supports both HTTP and gRPC
func NewAgent(server *http.Server, grpcServer *grpc.Server) Agent {
	return &agent{
		GRPCAgent: NewGRPC(grpcServer),
		HTTPAgent: NewHTTP(server),
	}
}

// NewAgentWithGRPC creates a new Agent instance that only supports gRPC
func NewAgentWithGRPC(grpcServer *grpc.Server) Agent {
	return &agent{
		GRPCAgent: NewGRPC(grpcServer),
		HTTPAgent: UnimplementedAgent,
	}
}

// NewAgentWithHTTP creates a new Agent instance that only supports HTTP
func NewAgentWithHTTP(server *http.Server) Agent {
	return &agent{
		GRPCAgent: UnimplementedAgent,
		HTTPAgent: NewHTTP(server),
	}
}
