/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package gateway

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
)

type Options struct {
	Prefix        string
	HostIP        string
	ServiceName   string
	Discovery     registry.Discovery
	NodeFilters   []selector.NodeFilter
	Middlewares   []middleware.Middleware
	EndpointFunc  func(scheme string, host string, addr string) (string, error)
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

func WithEndpointFunc(endpointFunc func(scheme string, host string, addr string) (string, error)) Option {
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
		o.HostIP = hostIp
	}
}

func WithServerOptions(opts ...transhttp.ServerOption) Option {
	return func(o *Options) {
		o.ServerOptions = append(o.ServerOptions, opts...)
	}
}
