/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package selector implements the functions, types, and interfaces for the module.
package service

import (
	"sync"

	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/go-kratos/kratos/v2/selector/p2c"
	"github.com/go-kratos/kratos/v2/selector/random"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	"github.com/origadmin/runtime/config"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

var (
	once    sync.Once
	builder selector.Builder
)

// DefaultSelectorOptionBuilder is the default instance of the service builder.
var DefaultSelectorOptionBuilder = &selectorOption{}

type selectorOption struct{}

func (s selectorOption) GRPC(cfg *configv1.Service_Selector) (transgrpc.ClientOption, error) {
	return WithGRPC(cfg)
}

func (s selectorOption) HTTP(cfg *configv1.Service_Selector) (transhttp.ClientOption, error) {
	return WithHTTP(cfg)
}

func WithHTTP(cfg *configv1.Service_Selector) (HTTPClientOption, error) {
	var options HTTPClientOption
	if cfg.GetVersion() != "" {
		v := filter.Version(cfg.Version)
		options = transhttp.WithNodeFilter(v)
	}
	SetGlobalSelector(cfg.GetBuilder())

	return options, nil
}

func WithGRPC(cfg *configv1.Service_Selector) (GRPCClientOption, error) {
	var options GRPCClientOption
	if cfg.GetVersion() != "" {
		v := filter.Version(cfg.Version)
		options = transgrpc.WithNodeFilter(v)
	}
	SetGlobalSelector(cfg.GetBuilder())

	return options, nil
}

// SetGlobalSelector sets the global selector.
func SetGlobalSelector(selectorType string) {
	if builder != nil {
		return
	}
	var b selector.Builder
	switch selectorType {
	case "random":
		b = random.NewBuilder()
	case "wrr":
		b = wrr.NewBuilder()
	case "p2c":
		b = p2c.NewBuilder()
	default:
		return
	}
	once.Do(func() {
		if b != nil {
			builder = b
			// Set global selector
			selector.SetGlobalSelector(builder)
		}
	})
}

var _ config.SelectorOption = selectorOption{}
