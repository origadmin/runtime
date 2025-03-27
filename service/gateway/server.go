/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package grpc implements the functions, types, and interfaces for the module.
package grpc

import (
	"net/url"
	"time"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goexts/generic/settings"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service/endpoint"
	"github.com/origadmin/toolkits/env"
	"github.com/origadmin/toolkits/errors"
)

const (
	Scheme   = "gateway"
	hostName = "HOST"
)

// NewServer Create a GRPC server instance
func NewServer(cfg *configv1.Service, ss ...OptionSetting) (*transhttp.Server, error) {
	log.Debugf("Creating new GRPC server instance with config: %+v", cfg)
	if cfg == nil {
		return nil, errors.New("service config is nil")
	}
	option := settings.ApplyDefaultsOrZero(ss...)
	serverOptions := []transhttp.ServerOption{
		transhttp.Middleware(option.Middlewares...),
	}
	if len(option.ServerOptions) > 0 {
		serverOptions = append(serverOptions, option.ServerOptions...)
	}
	if serviceHttp := cfg.GetHttp(); serviceHttp != nil {
		if serviceHttp.Network != "" {
			serverOptions = append(serverOptions, transhttp.Network(serviceHttp.Network))
		}
		if serviceHttp.Addr != "" {
			serverOptions = append(serverOptions, transhttp.Address(serviceHttp.Addr))
		}
		if serviceHttp.Timeout != 0 {
			serverOptions = append(serverOptions, transhttp.Timeout(time.Duration(serviceHttp.Timeout)))
		}
		if cfg.DynamicEndpoint && serviceHttp.Endpoint == "" {
			hostEnv := hostName
			if option.Prefix != "" {
				hostEnv = env.Var(option.Prefix, hostName)
			}
			dynamic, err := endpoint.GenerateDynamic(&endpoint.Option{
				EnvVar:       hostEnv,
				HostIP:       option.HostIP,
				EndpointFunc: nil,
			}, serviceHttp.Addr)
			if err != nil {
				return nil, err
			}
			serviceHttp.Endpoint = dynamic
		}
		log.Debugf("GRPC endpoint: %s", serviceHttp.Endpoint)
		if serviceHttp.Endpoint != "" {
			endpoint, err := url.Parse(serviceHttp.Endpoint)
			if err == nil {
				serverOptions = append(serverOptions, transhttp.Endpoint(endpoint))
			} else {
				log.Errorf("Failed to parse endpoint: %v", err)
			}
		}
	}

	srv := transhttp.NewServer(serverOptions...)
	return srv, nil
}
