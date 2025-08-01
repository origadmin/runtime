/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package service implements the functions, types, and interfaces for the module.
package service

import (
	"github.com/go-kratos/kratos/v2/transport"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/interfaces/factory"
)

type ServerBuilder interface {
	factory.Registry[ServerFactory]
	DefaultBuild(*configv1.Service, ...ServerOption) (transport.Server, error)
	Build(string, *configv1.Service, ...ServerOption) (transport.Server, error)
}

type ServerFactory interface {
	New(*configv1.Service, ...ServerOption) (transport.Server, error)
}

type ClientGRPCFactory interface {
	New(context.Context, *configv1.Service, ...GRPCOption) (*GRPCClient, error)
}

type ClientHTTPFactory interface {
	New(context.Context, *configv1.Service, ...HTTPOption) (*HTTPClient, error)
}

type (
	// Factory is an interface that defines a method for creating a new buildImpl.
	Factory interface {
		NewGRPCServer(*configv1.Service, ...GRPCOption) (*GRPCServer, error)
		NewHTTPServer(*configv1.Service, ...HTTPOption) (*HTTPServer, error)
		NewGRPCClient(context.Context, *configv1.Service, ...GRPCOption) (*GRPCClient, error)
		NewHTTPClient(context.Context, *configv1.Service, ...HTTPOption) (*HTTPClient, error)
	}
)

type Service struct{}
