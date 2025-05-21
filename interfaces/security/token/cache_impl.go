/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package token provides token caching functionality for security module
package token

import (
	"context"
	"time"

	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/interfaces/storage/cache"
)

const (
	CacheAccess  = "security:token:access"
	CacheRefresh = "security:token:refresh"
)

type StorageOption = func(*tokenCacheStorage)

func WithCache(c cache.Cache) StorageOption {
	return func(o *tokenCacheStorage) {
		o.c = c
	}
}

// tokenCacheStorage is the implementation of CacheStorage interface
type tokenCacheStorage struct {
	c cache.Cache
}

func (obj *tokenCacheStorage) Store(ctx context.Context, tokenStr string, duration time.Duration) error {
	return obj.c.Set(ctx, tokenStr, "", duration)
}

func (obj *tokenCacheStorage) Exist(ctx context.Context, tokenStr string) (bool, error) {
	ok, err := obj.c.Exists(ctx, tokenStr)
	switch {
	case ok:
		return true, nil
	default:
		return false, err
	}
}

func (obj *tokenCacheStorage) Remove(ctx context.Context, tokenStr string) error {
	return obj.c.Delete(ctx, tokenStr)
}

func (obj *tokenCacheStorage) Close(ctx context.Context) error {
	return obj.c.Close(ctx)
}

// New creates a new CacheStorage instance
func New(ss ...StorageOption) CacheStorage {
	service := settings.ApplyZero(ss)
	if service.c == nil {
		service.c = cache.NewMemoryCache()
	}
	return service
}
