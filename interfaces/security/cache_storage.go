/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security is a toolkit for security check and authorization
package security

import (
	"context"
	"time"

	"github.com/goexts/generic/settings"

	"github.com/origadmin/toolkits/storage/cache"
)

const (
	TokenCacheAccess  = "security:token:access"
	TokenCacheRefresh = "security:token:refresh"
)

type StorageSetting = func(*cacheStorage)

func WithCache(c cache.Cache) StorageSetting {
	return func(o *cacheStorage) {
		o.c = c
	}
}

// cacheStorage is the implementation of cacheStorage
type cacheStorage struct {
	c cache.Cache
}

func (obj *cacheStorage) Store(ctx context.Context, tokenStr string, duration time.Duration) error {
	return obj.c.Set(ctx, tokenStr, tokenStr, duration)
}

func (obj *cacheStorage) Exist(ctx context.Context, tokenStr string) (bool, error) {
	ok, err := obj.c.Exists(ctx, tokenStr)
	switch {
	case ok:
		return true, nil
	default:
		return false, err
	}
}

func (obj *cacheStorage) Remove(ctx context.Context, tokenStr string) error {
	return obj.c.Delete(ctx, tokenStr)
}

func (obj *cacheStorage) Close(ctx context.Context) error {
	return obj.c.Close(ctx)
}

// NewCacheStorage creates a new cacheStorage with a c and optional StoreOptions
func NewCacheStorage(ss ...StorageSetting) CacheStorage {
	service := settings.Apply(&cacheStorage{}, ss)
	if service.c == nil {
		service.c = cache.NewMemoryCache()
	}
	return service
}
