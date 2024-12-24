/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package agent implements the functions, types, and interfaces for the module.
package agent

import (
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"
)

type unimplementedAgent struct{}

var UnimplementedAgent Agent = &unimplementedAgent{}

func (u unimplementedAgent) URI() string {
	return ""
}

func (u unimplementedAgent) HTTPServer() *transhttp.Server {
	return nil
}

func (u unimplementedAgent) Route() *transhttp.Router {
	return nil
}

func (u unimplementedAgent) Server() *transgrpc.Server {
	return nil
}

func (u unimplementedAgent) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	return
}
