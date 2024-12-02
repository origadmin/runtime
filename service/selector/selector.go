/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package selector implements the functions, types, and interfaces for the module.
package selector

import (
	"sync"

	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/go-kratos/kratos/v2/selector/p2c"
	"github.com/go-kratos/kratos/v2/selector/random"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

var (
	once = &sync.Once{}
)

func WithHTTP(cfg *configv1.Service_Selector) (transhttp.ClientOption, error) {
	var options transhttp.ClientOption
	if cfg.GetVersion() != "" {
		v := filter.Version(cfg.Version)
		options = transhttp.WithNodeFilter(v)
	}
	SetGlobalSelector(cfg.GetBuilder())

	return options, nil
}

func WithGRPC(cfg *configv1.Service_Selector) (transgrpc.ClientOption, error) {
	var options transgrpc.ClientOption
	if cfg.GetVersion() != "" {
		v := filter.Version(cfg.Version)
		options = transgrpc.WithNodeFilter(v)
	}
	SetGlobalSelector(cfg.GetBuilder())

	return options, nil
}

// SetGlobalSelector sets the global selector.
func SetGlobalSelector(selectorType string) {
	var builder selector.Builder
	switch selectorType {
	case "random":
		builder = random.NewBuilder()
	case "wrr":
		builder = wrr.NewBuilder()
	case "p2c":
		builder = p2c.NewBuilder()
	default:
		return
	}
	if builder != nil {
		once.Do(func() {
			selector.SetGlobalSelector(builder)
		})
	}
}
