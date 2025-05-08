/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package service implements the functions, types, and interfaces for the module.
package service

import (
	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces/builder"
)

type (
	// Builder is an interface that defines a method for registering a buildImpl.
	Builder interface {
		builder.Builder[Factory]
		Factory
	}
	// Factory is an interface that defines a method for creating a new buildImpl.
	Factory interface {
		NewGRPCServer(*configv1.Service, ...GRPCOption) (*GRPCServer, error)
		NewHTTPServer(*configv1.Service, ...HTTPOption) (*HTTPServer, error)
		NewGRPCClient(context.Context, *configv1.Service, ...GRPCOption) (*GRPCClient, error)
		NewHTTPClient(context.Context, *configv1.Service, ...HTTPOption) (*HTTPClient, error)
	}
)

type Service struct{}
