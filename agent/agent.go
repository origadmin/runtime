/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package agent implements the functions, types, and interfaces for the module.
package agent

import (
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"
)

const (
	ApiVersionV1   = "/api/v1"
	DefaultPrefix  = "/api"
	DefaultVersion = "v1"
)

type Agent interface {
	HTTPAgent
	GRPCAgent
}

type HTTPAgent interface {
	URI() string
	HTTPServer() *http.Server
	Route() *http.Router
}

type GRPCAgent interface {
	Server() *grpc.Server
	RegisterService(desc *grpc.ServiceDesc, impl interface{})
}

type agent struct {
	GRPCAgent
	HTTPAgent
}

func NewAgent(server *http.Server, grpcServer *grpc.Server) Agent {
	return &agent{
		GRPCAgent: NewGRPC(grpcServer),
		HTTPAgent: NewHTTP(server),
	}
}

func NewAgentWithGRPC(grpcServer *grpc.Server) Agent {
	return &agent{
		GRPCAgent: NewGRPC(grpcServer),
		HTTPAgent: UnimplementedAgent,
	}
}

func NewAgentWithHTTP(server *http.Server) Agent {
	return &agent{
		GRPCAgent: UnimplementedAgent,
		HTTPAgent: NewHTTP(server),
	}
}
