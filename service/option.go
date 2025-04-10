/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package service implements the functions, types, and interfaces for the module.
package service

import (
	"github.com/origadmin/runtime/service/grpc"
	"github.com/origadmin/runtime/service/http"
)

const DefaultHostEnv = "HOST"

type (
	EndpointFunc = func(scheme string, host string, addr string) (string, error)
)

// HTTPOption is the type for HTTP option settings.
type HTTPOption = http.Options

// GRPCOption is the type for gRPC option settings.
type GRPCOption = grpc.Options
