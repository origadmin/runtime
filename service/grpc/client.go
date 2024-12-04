/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package grpc implements the functions, types, and interfaces for the module.
package grpc

import (
	"time"

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/helpers"
	"google.golang.org/grpc"

	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/service/selector"
)

const defaultTimeout = 5 * time.Second

// NewClient Creating a GRPC client instance
func NewClient(ctx context.Context, service *configv1.Service, rc *config.RuntimeConfig) (*grpc.ClientConn, error) {
	if rc == nil {
		rc = config.DefaultRuntimeConfig
	}
	serviceOption := rc.Service()
	selectorOption := rc.Selector()
	var ms []middleware.Middleware
	ms = middleware.NewClient(service.GetMiddleware())
	if serviceOption.Middlewares != nil {
		ms = append(ms, serviceOption.Middlewares...)
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

	if serviceOption.Discovery != nil {
		endpoint := helpers.ServiceName(service.GetName())
		options = append(options, transgrpc.WithEndpoint(endpoint),
			transgrpc.WithDiscovery(serviceOption.Discovery))
	}

	if selectorOption.GRPC == nil {
		selectorOption.GRPC = selector.DefaultGRPC
	}
	if serviceSelector := service.GetSelector(); serviceSelector != nil {
		if option, err := selectorOption.GRPC(serviceSelector); err == nil && option != nil {
			options = append(options, option)
		}
	}

	conn, err := transgrpc.DialInsecure(ctx, options...)
	if err != nil {
		return nil, errors.Errorf("dial grpc client [%s] failed: %s", service.GetName(), err.Error())
	}

	return conn, nil
}
