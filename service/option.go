/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package builder implements the functions, types, and interfaces for the module.
package service

import (
	"github.com/origadmin/runtime/service/grpc"
	"github.com/origadmin/runtime/service/http"
)

type (
	EndpointFunc = func(scheme string, host string, addr string) (string, error)
)

// Option represents a set of configuration options for a builder.
type Option struct {
	http []http.OptionSetting
	grpc []grpc.OptionSetting
}

type OptionSetting = func(option *Option)

func WithGRPC(option ...grpc.OptionSetting) OptionSetting {
	return func(o *Option) {
		o.grpc = append(o.grpc, option...)
	}
}

func WithHTTP(option ...http.OptionSetting) OptionSetting {
	return func(o *Option) {
		o.http = append(o.http, option...)
	}
}
