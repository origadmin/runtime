/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package http

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
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
	ClientOptions []transhttp.ClientOption
	ServerOptions []transhttp.ServerOption
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

func WithClientOptions(options ...transhttp.ClientOption) Option {
	return func(o *Options) {
		o.ClientOptions = append(o.ClientOptions, options...)
	}
}
func WithServerOptions(options ...transhttp.ServerOption) Option {
	return func(o *Options) {
		o.ServerOptions = append(o.ServerOptions, options...)
	}
}
