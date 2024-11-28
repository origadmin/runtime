/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package data implements the functions, types, and interfaces for the module.
package data

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/toolkits/errors"
	"gith
	"github.com/origadmin/toolkits/storage/cache"
)

const (
	ErrCacheConfigNil = errors.String("cache: config is nil")
)

type (
	Cache = cache.Cache
)

func OpenCache(cfg *configv1.Data) (Cache, error) {
	if cfg == nil {
		return nil, ErrCacheConfigNil
	}

	if c := cfg.GetCache().GetMemory(); c != nil {

	}

	return nil, errors.String("cache: unknown cache type")
}
