/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package grpc implements the functions, types, and interfaces for the module.
package grpc

import (
	"time"

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/goexts/generic/settings"
	"google.golang.org/grpc"

	"github.com/origadmin/runtime/config"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/service/selector"
	"github.com/origadmin/toolkits/context"
)

const defaultTimeout = 5 * time.Second

// NewClient Creating a GRPC client instance
func NewClient(ctx context.Context, service *configv1.Service, opts ...config.ServiceSetting) (*grpc.ClientConn, error) {
	option := settings.Apply(&config.ServiceOption{}, opts)
	var ms []middleware.Middleware
	ms = middleware.NewClient(service.GetMiddleware())
	if option.Middlewares != nil {
		ms = append(ms, option.Middlewares...)
	}

	timeout := defaultTimeout
	if serviceGrpc := service.GetGrpc(); serviceGrpc != nil {
		if serviceGrpc.Timeout != nil {
			timeout = serviceGrpc.Timeout.AsDuration()
		}
	}

	options := []transgrpc.ClientOption{
		transgrpc.WithTimeout(timeout),
		transgrpc.WithMiddleware(ms...),
	}

	if option.Discovery != nil {
		endpoint := helpers.ServiceDiscoveryName(service.GetName())
		options = append(options,
			transgrpc.WithEndpoint(endpoint),
			transgrpc.WithDiscovery(option.Discovery),
		)
	}

	if option, err := selector.NewGRPC(service.GetSelector()); err == nil {
		options = append(options, option)
	}
	conn, err := transgrpc.DialInsecure(ctx, options...)

	if err != nil {
		return nil, errors.Errorf("dial grpc client [%s] failed: %s", service.GetName(), err.Error())
	}

	return conn, nil
}
