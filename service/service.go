/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package builder implements the functions, types, and interfaces for the module.
package service

import (
	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type (
	// Builder is an interface that defines a method for registering a builder builder.
	Builder interface {
		Factory
		RegisterServiceBuilder(name string, factory Factory)
	}
	// Factory is an interface that defines a method for creating a new builder.
	Factory interface {
		NewGRPCServer(*configv1.Service, ...OptionSetting) (*GRPCServer, error)
		NewHTTPServer(*configv1.Service, ...OptionSetting) (*HTTPServer, error)
		NewGRPCClient(context.Context, *configv1.Service, ...OptionSetting) (*GRPCClient, error)
		NewHTTPClient(context.Context, *configv1.Service, ...OptionSetting) (*HTTPClient, error)
	}
)

type Service struct{}
