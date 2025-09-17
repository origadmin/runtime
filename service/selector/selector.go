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

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/toolkits/errors"
)

const (
	Random = "random"
	WRR    = "wrr"
	P2C    = "p2c"
)

var (
	once    sync.Once
	builder selector.Builder
)

func NewFilter(cfg *configv1.Service_Selector) (selector.NodeFilter, error) {
	// Check if the version is specified in the configuration
	if cfg.GetVersion() != "" {
		// Return the version filter and no error
		return filter.Version(cfg.Version), nil
	}
	// Return the node filter and no error
	return nil, errors.New("version is nil")
}

// SetSelectorGlobalSelector sets the global selector.
func SetSelectorGlobalSelector(selectorType string) {
	if builder != nil {
		return
	}
	var b selector.Builder
	switch selectorType {
	case Random:
		b = random.NewBuilder()
	case WRR:
		b = wrr.NewBuilder()
	case P2C:
		b = p2c.NewBuilder()
	default:
		log.Warnf("selector type %s is not supported", selectorType)
		return
	}
	once.Do(func() {
		if b != nil {
			builder = b
			// Set global selector
			SetGlobalSelector(builder)
		}
	})
}
