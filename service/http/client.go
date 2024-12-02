/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package grpc implements the functions, types, and interfaces for the module.
package grpc

import (
	"time"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goexts/generic/settings"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/helpers"

	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/service/selector"
)

const defaultTimeout = 5 * time.Second

// NewClient Creating an HTTP client instance.
func NewClient(ctx context.Context, service *configv1.Service, opts ...config.ServiceSetting) (*transhttp.Client, error) {
	option := settings.Apply(&config.ServiceOption{}, opts)
	var ms []middleware.Middleware
	ms = middleware.NewClient(service.GetMiddleware())
	if option.Middlewares != nil {
		ms = append(ms, option.Middlewares...)
	}

	timeout := defaultTimeout
	if serviceHttp := service.GetHttp(); serviceHttp != nil {
		if serviceHttp.Timeout != nil {
			timeout = serviceHttp.Timeout.AsDuration()
		}
	}

	options := []transhttp.ClientOption{
		transhttp.WithTimeout(timeout),
		transhttp.WithMiddleware(ms...),
	}

	if option.Discovery != nil {
		endpoint := helpers.ServiceDiscoveryName(service.GetName())
		options = append(options,
			transhttp.WithEndpoint(endpoint),
			transhttp.WithDiscovery(option.Discovery),
		)
	}

	if option, err := selector.NewHTTP(service.GetSelector()); err == nil {
		options = append(options, option)
	}
	conn, err := transhttp.NewClient(ctx, options...)

	if err != nil {
		return nil, errors.Errorf("dial grpc client [%s] failed: %s", service.GetName(), err.Error())
	}

	return conn, nil
}
