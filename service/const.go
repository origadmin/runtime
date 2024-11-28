/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package service implements the functions, types, and interfaces for the module.
package service

import (
	"time"

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"
)

const DefaultTimeout = 5 * time.Second

type (
	GRPCServer = transgrpc.Server
	HTTPServer = transhttp.Server
	GRPCClient = grpc.ClientConn
	HTTPClient = transhttp.Client
)
