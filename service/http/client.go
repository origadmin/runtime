/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package http implements the functions, types, and interfaces for the module.
package http

import (
	"time"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/helpers"
)

const defaultTimeout = 5 * time.Second

// NewClient Creating an HTTP client instance.
func NewClient(ctx context.Context, cfg *configv1.Service, ss ...Option) (*transhttp.Client, error) {
	if cfg == nil {
		//bootstrap = config.DefaultRuntimeConfig
		return nil, errors.New("service config is nil")
	}
	option := settings.ApplyDefaultsOrZero(ss...)
	timeout := defaultTimeout
	if serviceHttp := cfg.GetHttp(); serviceHttp != nil {
		if serviceHttp.Timeout != 0 {
			timeout = time.Duration(serviceHttp.Timeout)
		}
	}
	clientOptions := []transhttp.ClientOption{
		transhttp.WithTimeout(timeout),
		transhttp.WithMiddleware(option.Middlewares...),
	}
	if len(option.ClientOptions) > 0 {
		clientOptions = append(clientOptions, option.ClientOptions...)
	}
	if option.Discovery != nil {
		endpoint := helpers.ServiceName(option.ServiceName)
		log.Debugf("http service [%s] discovery endpoint [%s]", option.ServiceName, endpoint)
		clientOptions = append(clientOptions,
			transhttp.WithEndpoint(endpoint),
			transhttp.WithDiscovery(option.Discovery),
		)
	}

	if serviceSelector := cfg.GetSelector(); serviceSelector != nil {
		if len(option.NodeFilters) > 0 {
			clientOptions = append(clientOptions, transhttp.WithNodeFilter(option.NodeFilters...))
		}
	}

	conn, err := transhttp.NewClient(ctx, clientOptions...)
	if err != nil {
		return nil, errors.Errorf("dial http client [%s] failed: %s", cfg.GetName(), err.Error())
	}

	return conn, nil
}
