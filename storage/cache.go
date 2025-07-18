/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package storage implements the functions, types, and interfaces for the module.
package storage

import (
	"fmt"

	storageiface "github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/toolkits/errors"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
)

const (
	ErrCacheConfigNil = errors.String("cache: config is nil")
)

type (
	Cache = cache.Cache
)

func OpenCache(cfg *configv1.Storage) (storageiface.Cache, error) {
	if cfg == nil {
		return nil, ErrCacheConfigNil
	}
	cacheCfg := cfg.GetCache()
	fmt.Println(cacheCfg)
	return nil, errors.String("cache: unknown cache type")
}
