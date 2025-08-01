/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package service implements the functions, types, and interfaces for the module.
package service

import (
	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/service/grpc"
	"github.com/origadmin/runtime/service/http"
)

const DefaultHostEnv = "HOST"

type (
	EndpointFunc = func(scheme string, host string, addr string) (string, error)
)

// HTTPOption is the type for HTTP option settings.
type HTTPOption = http.Option

// GRPCOption is the type for gRPC option settings.
type GRPCOption = grpc.Option

type Options struct {
	grpc []GRPCOption
	http []HTTPOption
}

type ServerOption = func(o *Options)

func (o Options) ToGRPC() GRPCOption {
	return func(opts *grpc.Options) {
		settings.Apply(opts, o.grpc)
	}
}

func WithGRPC(opts ...GRPCOption) ServerOption {
	return func(o *Options) {
		o.grpc = opts
	}
}

func (o Options) ToHTTP() http.Option {
	return func(opts *http.Options) {
		settings.Apply(opts, o.http)
	}
}

func WithHTTP(opts ...HTTPOption) ServerOption {
	return func(o *Options) {
		o.http = opts
	}
}
