/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package agent implements the functions, types, and interfaces for the module.
package agent

import (
	"google.golang.org/grpc"
)

type grpcAgent struct {
	prefix  string
	version string
	server  *grpc.Server
}

func (g *grpcAgent) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	g.server.RegisterService(desc, impl)
}

func (g *grpcAgent) Server() *grpc.Server {
	return g.server
}

func NewGRPC(server *grpc.Server) GRPCAgent {
	return &grpcAgent{
		prefix:  DefaultPrefix,
		version: DefaultVersion,
		server:  server,
	}
}
