/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package grpc

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/origadmin/runtime/service/endpoint"
	"github.com/origadmin/toolkits/env"
)

type EndpointFunc = func(scheme string, host string, addr string) (string, error)

type Options struct {
	Prefix        string
	HostIp        string
	ServiceName   string
	Discovery     registry.Discovery
	NodeFilters   []selector.NodeFilter
	Middlewares   []middleware.Middleware
	EndpointFunc  EndpointFunc
	ClientOptions []transgrpc.ClientOption
	ServerOptions []transgrpc.ServerOption
}

type Option = func(o *Options)

func WithNodeFilter(filters ...selector.NodeFilter) Option {
	return func(o *Options) {
		o.NodeFilters = append(o.NodeFilters, filters...)
	}
}
func WithDiscovery(serviceName string, discovery registry.Discovery) Option {
	return func(o *Options) {
		o.ServiceName = serviceName
		o.Discovery = discovery
	}
}

func WithMiddlewares(middlewares ...middleware.Middleware) Option {
	return func(o *Options) {
		o.Middlewares = append(o.Middlewares, middlewares...)
	}
}

func WithEndpointFunc(endpointFunc EndpointFunc) Option {
	return func(o *Options) {
		o.EndpointFunc = endpointFunc
	}
}
func WithPrefix(prefix string) Option {
	return func(o *Options) {
		o.Prefix = prefix
	}
}

func WithHostIp(hostIp string) Option {
	return func(o *Options) {
		o.HostIp = hostIp
	}
}

func WithClientOptions(opts ...transgrpc.ClientOption) Option {
	return func(o *Options) {
		o.ClientOptions = append(o.ClientOptions, opts...)
	}
}

func WithServerOptions(opts ...transgrpc.ServerOption) Option {
	return func(o *Options) {
		o.ServerOptions = append(o.ServerOptions, opts...)
	}
}

func parseEndpointOption(opt *Options) *endpoint.Options {
	hostEnv := hostName
	if opt == nil {
		return &endpoint.Options{
			EnvVar:       hostEnv,
			EndpointFunc: endpoint.ExtractIP,
		}
	}
	if opt.Prefix != "" {
		hostEnv = env.Var(opt.Prefix, hostName)
	}
	ep := &endpoint.Options{
		EnvVar:       hostEnv,
		HostIP:       opt.HostIp,
		EndpointFunc: opt.EndpointFunc,
	}
	return ep
}
