/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package token provides token caching functionality for security module
package token

import (
	"context"
	"fmt"
	"time"

	"github.com/goexts/generic/settings"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/runtime/storage"
)

const (
	CacheAccess  = "security:token:access"
	CacheRefresh = "security:token:refresh"
)

type StorageOption = func(*tokenCacheStorage)

func WithCache(c storageiface.Cache) StorageOption {
	return func(o *tokenCacheStorage) {
		o.c = c
	}
}

// tokenCacheStorage is the implementation of CacheStorage interface
type tokenCacheStorage struct {
	c storageiface.Cache
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
		defaultCacheConfig := &configv1.Cache{
			Driver: "memory",
			Memory: &configv1.Memory{},
		}
		c, err := storage.New(defaultCacheConfig)
		if err != nil {
			// Handle error, perhaps log it or panic if cache is critical
			panic(fmt.Sprintf("failed to create default memory cache: %v", err))
		}
		service.c = c
	}
	return service
}
