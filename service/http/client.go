/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package http implements the functions, types, and interfaces for the module.
package http

import (
	"time"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goexts/generic/settings"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service/selector"
	"github.com/origadmin/runtime/service/tls"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/helpers"
)

const defaultTimeout = 5 * time.Second

// NewClient Creating an HTTP client instance.
func NewClient(ctx context.Context, cfg *configv1.Service, options ...Option) (*transhttp.Client, error) {
	if cfg == nil {
		return nil, errors.New("service config is nil")
	}
	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/http"))
	option := settings.ApplyZero(options)
	timeout := defaultTimeout
	clientOptions := option.ClientOptions
	if serviceHttp := cfg.GetHttp(); serviceHttp != nil {
		if serviceHttp.Timeout != 0 {
			timeout = time.Duration(serviceHttp.Timeout * 1e6)
		}
		if serviceHttp.UseTls {
			tlsConfig, err := tls.NewClientTLSConfig(serviceHttp.GetTlsConfig())
			if err != nil {
				return nil, err
			}
			if tlsConfig != nil {
				option.ClientOptions = append(option.ClientOptions, transhttp.WithTLSConfig(tlsConfig))
			}
		}
	}
	clientOptions = append(clientOptions, transhttp.WithTimeout(timeout))
	if len(option.Middlewares) > 0 {
		clientOptions = append(clientOptions, transhttp.WithMiddleware(option.Middlewares...))
	}
	if option.Discovery != nil {
		endpoint := helpers.ServiceDiscovery(option.ServiceName)
		ll.Debugw("msg", "init with discovery", "service", "http", "name", option.ServiceName, "endpoint", endpoint)
		clientOptions = append(clientOptions,
			transhttp.WithEndpoint(endpoint),
			transhttp.WithDiscovery(option.Discovery),
		)
	}
	if serviceSelector := cfg.GetSelector(); serviceSelector != nil {
		filter, err := selector.NewFilter(cfg.GetSelector())
		if err == nil {
			option.NodeFilters = append(option.NodeFilters, filter)
		}
	}
	if len(option.NodeFilters) > 0 {
		clientOptions = append(clientOptions, transhttp.WithNodeFilter(option.NodeFilters...))
	}

	conn, err := transhttp.NewClient(ctx, clientOptions...)
	if err != nil {
		return nil, errors.Errorf("dial http client [%s] failed: %s", cfg.GetName(), err.Error())
	}

	return conn, nil
}
