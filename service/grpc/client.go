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

	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service/selector"
	"github.com/origadmin/runtime/service/tls"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/helpers"
)

const defaultTimeout = 5 * time.Second

// NewClient Creating a GRPC client instance
func NewClient(ctx context.Context, cfg *configv1.Service, options ...Option) (*grpc.ClientConn, error) {
	if cfg == nil {
		//bootstrap = config.DefaultRuntimeConfig
		return nil, errors.New("service config is nil")
	}
	option := settings.ApplyDefaultsOrZero(options...)
	timeout := defaultTimeout
	clientOptions := []transgrpc.ClientOption{
		transgrpc.WithTimeout(timeout),
		transgrpc.WithMiddleware(option.Middlewares...),
	}
	if serviceGrpc := cfg.GetGrpc(); serviceGrpc != nil {
		if serviceGrpc.Timeout != 0 {
			timeout = time.Duration(serviceGrpc.Timeout * 1e6)
		}
		if serviceGrpc.UseTls {
			tlsConfig, err := tls.NewClientTLSConfig(serviceGrpc.GetTlsConfig())
			if err != nil {
				return nil, err
			}
			if tlsConfig != nil {
				option.ClientOptions = append(option.ClientOptions, transgrpc.WithTLSConfig(tlsConfig))
			}
		}
	}
	if len(option.ClientOptions) > 0 {
		clientOptions = append(clientOptions, option.ClientOptions...)
	}

	if option.Discovery != nil {
		endpoint := helpers.ServiceDiscovery(option.ServiceName)
		log.Debugf("grpc service [%s] discovery endpoint [%s]", option.ServiceName, endpoint)
		clientOptions = append(clientOptions,
			transgrpc.WithEndpoint(endpoint),
			transgrpc.WithDiscovery(option.Discovery))
	}
	if serviceSelector := cfg.GetSelector(); serviceSelector != nil {
		filter, err := selector.NewFilter(cfg.GetSelector())
		if err == nil {
			option.NodeFilters = append(option.NodeFilters, filter)
		}

	}
	if len(option.NodeFilters) > 0 {
		clientOptions = append(clientOptions, transgrpc.WithNodeFilter(option.NodeFilters...))
	}

	conn, err := transgrpc.DialInsecure(ctx, clientOptions...)
	if err != nil {
		return nil, errors.Errorf("dial grpc client [%s] failed: %s", cfg.GetName(), err.Error())
	}

	return conn, nil
}
