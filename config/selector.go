/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type (
	// SelectorGRPCFunc is a function type that returns a gRPC client option.
	// It takes a service selector configuration as input and returns a client option and an error.
	SelectorGRPCFunc = func(cfg *configv1.Service_Selector) (transgrpc.ClientOption, error)

	// SelectorHTTPFunc is a function type that returns an HTTP client option.
	// It takes a service selector configuration as input and returns a client option and an error.
	SelectorHTTPFunc = func(cfg *configv1.Service_Selector) (transhttp.ClientOption, error)
)

// SelectorOption represents a configuration option for a selector.
type SelectorOption struct {
	// GRPC is a function that returns a gRPC client option.
	GRPC SelectorGRPCFunc
	// HTTP is a function that returns an HTTP client option.
	HTTP SelectorHTTPFunc
}

// SelectorOptionSetting is a function type that sets a selector option.
type SelectorOptionSetting = func(option *SelectorOption)
