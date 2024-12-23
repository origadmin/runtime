/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	selectorv1 "github.com/origadmin/runtime/gen/go/middleware/selector/v1"
	"github.com/origadmin/runtime/middleware/selector"
)

func Selector(cfg *selectorv1.Selector, matchFunc selector.MatchFunc) selector.Selector {
	if cfg == nil || !cfg.Enabled {
		return selector.Unfiltered()
	}
	return selector.NewSelectorFilter(cfg.GetNames(), matchFunc)
}
