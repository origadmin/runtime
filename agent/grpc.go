/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package agent implements the functions, types, and interfaces for the module.
package agent

import (
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
)

type grpcAgent struct {
	prefix  string
	version string
	server  *transgrpc.Server
}

func (g *grpcAgent) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	g.server.RegisterService(desc, impl)
}

func (g *grpcAgent) Server() *transgrpc.Server {
	return g.server
}

func NewGRPC(server *transgrpc.Server) GRPCAgent {
	return &grpcAgent{
		prefix:  DefaultPrefix,
		version: DefaultVersion,
		server:  server,
	}
}
