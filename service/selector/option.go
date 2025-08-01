/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package selector implements the functions, types, and interfaces for the module.
package selector

import (
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
)

type (
	// GRPCFunc is a function type that returns a gRPC client option.
	// It takes a service selector configuration as input and returns a client option and an error.
	GRPCFunc = func(cfg *configv1.Service_Selector) (transgrpc.ClientOption, error)

	// HTTPFunc is a function type that returns an HTTP client option.
	// It takes a service selector configuration as input and returns a client option and an error.
	HTTPFunc = func(cfg *configv1.Service_Selector) (transhttp.ClientOption, error)
)

// Options represents a configuration option for a selector.
type Options struct {
	// GRPC is a function that returns a gRPC client option.
	GRPC GRPCFunc
	// HTTP is a function that returns an HTTP client option.
	HTTP HTTPFunc
}

// Option is a function type that sets a selector option.
type Option = func(config *Options)
