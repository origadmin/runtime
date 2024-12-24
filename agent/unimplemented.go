/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package agent implements the functions, types, and interfaces for the module.
package agent

import (
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"
)

type unimplementedAgent struct{}

var UnimplementedAgent Agent = &unimplementedAgent{}

func (u unimplementedAgent) URI() string {
	return ""
}

func (u unimplementedAgent) HTTPServer() *http.Server {
	return nil
}

func (u unimplementedAgent) Route() *http.Router {
	return nil
}

func (u unimplementedAgent) Server() *grpc.Server {
	return nil
}

func (u unimplementedAgent) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	return
}
