/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package http implements the functions, types, and interfaces for the module.
package http

import (
	"time"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
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
func NewClient(ctx context.Context, service *configv1.Service, rc *config.RuntimeConfig) (*transhttp.Client, error) {
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
	if serviceHttp := service.GetHttp(); serviceHttp != nil {
		if serviceHttp.Timeout != nil {
			timeout = serviceHttp.Timeout.AsDuration()
		}
	}

	options := []transhttp.ClientOption{
		transhttp.WithTimeout(timeout),
		transhttp.WithMiddleware(ms...),
	}

	if serviceOption.Discovery != nil {
		endpoint := helpers.ServiceDiscoveryName(service.GetName())
		options = append(options,
			transhttp.WithEndpoint(endpoint),
			transhttp.WithDiscovery(serviceOption.Discovery),
		)
	}

	if selectorOption.HTTP == nil {
		selectorOption.HTTP = selector.DefaultHTTP
	}
	if serviceSelector := service.GetSelector(); serviceSelector != nil {
		if option, err := selectorOption.HTTP(serviceSelector); err == nil && option != nil {
			options = append(options, option)
		}
	}

	conn, err := transhttp.NewClient(ctx, options...)
	if err != nil {
		return nil, errors.Errorf("dial http client [%s] failed: %s", service.GetName(), err.Error())
	}

	return conn, nil
}
