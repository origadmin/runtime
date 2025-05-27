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

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/toolkits/errors"
)

var (
	once    sync.Once
	builder selector.Builder
)

// DefaultHTTP returns a default HTTP client option based on the provided service selector configuration.
//
// If the version is specified in the configuration, it sets a node filter with the version.
func DefaultHTTP(cfg *configv1.Service_Selector) (transhttp.ClientOption, error) {
	// Initialize an empty client option
	var options transhttp.ClientOption

	// Check if the version is specified in the configuration
	if cfg.GetVersion() != "" {
		// Create a version filter based on the configuration version
		v := filter.Version(cfg.Version)
		// Set the node filter option
		options = transhttp.WithNodeFilter(v)
	}

	// Set the global selector with the provided builder
	SetGlobalSelector(cfg.GetBuilder())

	// Return the client option and no error
	return options, nil
}

// DefaultGRPC returns a default gRPC client option based on the provided service selector configuration.
//
// If the version is specified in the configuration, it sets a node filter with the version.
func DefaultGRPC(cfg *configv1.Service_Selector) (transgrpc.ClientOption, error) {
	// Initialize an empty client option
	var options transgrpc.ClientOption

	// Check if the version is specified in the configuration
	if cfg.GetVersion() != "" {
		// Create a version filter based on the configuration version
		v := filter.Version(cfg.Version)
		// Set the node filter option
		options = transgrpc.WithNodeFilter(v)
	}

	// Set the global selector with the provided builder
	SetGlobalSelector(cfg.GetBuilder())

	// Return the client option and no error
	return options, nil
}

func NewFilter(cfg *configv1.Service_Selector) (selector.NodeFilter, error) {
	// Check if the version is specified in the configuration
	if cfg.GetVersion() != "" {
		// Create a version filter based on the configuration version
		// Set the global selector with the provided builder
		SetGlobalSelector(cfg.GetBuilder())
		// Return the version filter and no error
		return filter.Version(cfg.Version), nil
	}
	// Return the node filter and no error
	return nil, errors.New("version is nil")
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
