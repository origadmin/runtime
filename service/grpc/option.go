/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package grpc

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type Option struct {
	Config       *configv1.Service
	ServiceName  string
	Discovery    registry.Discovery
	NodeFilters  []selector.NodeFilter
	Middlewares  []middleware.Middleware
	EndpointFunc func(scheme string, host string, addr string) (string, error)
}

type OptionSetting = func(o *Option)

func WithNodeFilter(filters ...selector.NodeFilter) OptionSetting {
	return func(o *Option) {
		o.NodeFilters = append(o.NodeFilters, filters...)
	}
}
