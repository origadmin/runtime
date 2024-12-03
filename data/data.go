/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package data implements the functions, types, and interfaces for the module.
package data

import (
	"github.com/origadmin/toolkits/errors"

	"github.com/origadmin/runtime/config"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type Service interface {
	GetCache(name string) Cache
}

type serviceImpl struct {
	cache  Cache // Default cache
	caches map[string]Cache
	// Other fields and methods
}

func (s serviceImpl) GetCache(name string) Cache {
	if cache, ok := s.caches[name]; ok {
		return cache
	}
	return s.cache
}

func NewService(data *configv1.Data, rc *config.RuntimeConfig) (Service, error) {
	cacheConfigs := data.GetCaches()
	caches := make(map[string]Cache, len(cacheConfigs))
	var defaultCache Cache
	for _, dataCache := range cacheConfigs {
		cache, err := OpenCache(dataCache)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open cache")
		}
		if dataCache.GetName() == "default" || dataCache.GetName() == "" {
			defaultCache = cache
		}
		caches[dataCache.GetName()] = cache
	}

	return &serviceImpl{
		cache:  defaultCache,
		caches: caches,
	}, nil
}

var _ Service = (*serviceImpl)(nil)
