/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package grpc

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
)

type Option struct {
	Prefix        string
	HostIP        string
	ServiceName   string
	Discovery     registry.Discovery
	NodeFilters   []selector.NodeFilter
	Middlewares   []middleware.Middleware
	EndpointFunc  func(scheme string, host string, addr string) (string, error)
	ServerOptions []transhttp.ServerOption
}

type OptionSetting = func(o *Option)

func WithNodeFilter(filters ...selector.NodeFilter) OptionSetting {
	return func(o *Option) {
		o.NodeFilters = append(o.NodeFilters, filters...)
	}
}
func WithDiscovery(serviceName string, discovery registry.Discovery) OptionSetting {
	return func(o *Option) {
		o.ServiceName = serviceName
		o.Discovery = discovery
	}
}

func WithMiddlewares(middlewares ...middleware.Middleware) OptionSetting {
	return func(o *Option) {
		o.Middlewares = append(o.Middlewares, middlewares...)
	}
}

func WithEndpointFunc(endpointFunc func(scheme string, host string, addr string) (string, error)) OptionSetting {
	return func(o *Option) {
		o.EndpointFunc = endpointFunc
	}
}
func WithPrefix(prefix string) OptionSetting {
	return func(o *Option) {
		o.Prefix = prefix
	}
}

func WithHostIp(hostIp string) OptionSetting {
	return func(o *Option) {
		o.HostIP = hostIp
	}
}

func WithServerOptions(opts ...transhttp.ServerOption) OptionSetting {
	return func(o *Option) {
		o.ServerOptions = append(o.ServerOptions, opts...)
	}
}
