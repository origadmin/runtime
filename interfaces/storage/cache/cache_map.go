/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package cache implements the functions, types, and interfaces for the module.
package cache

import (
	"context"
	"sync"
	"time"
)

type element struct {
	value    string
	expireAt time.Time
}

type mapCache struct {
	elems *sync.Pool
	maps  *sync.Map
}

func (s *mapCache) Clear(ctx context.Context) error {
	s.maps.Clear()
	return nil
}

func (s *mapCache) Exists(ctx context.Context, key string) (bool, error) {
	_, ok := s.maps.Load(key)
	return ok, nil
}

func (s *mapCache) Close(ctx context.Context) error {
	return nil
}

func (s *mapCache) Get(ctx context.Context, key string) (string, error) {
	value, ok := s.maps.Load(key)
	if !ok {
		return "", ErrNotFound
	}
	elem, ok := value.(*element)
	if !ok {
		return "", ErrInvalidElement
	}
	if !elem.expireAt.IsZero() && elem.expireAt.Before(time.Now()) {
		s.maps.Delete(key)
		s.elems.Put(elem)
		return "", ErrExpired
	}
	return elem.value, nil
}

func (s *mapCache) GetAndDelete(ctx context.Context, key string) (string, error) {
	value, ok := s.maps.LoadAndDelete(key)
	if !ok {
		return "", ErrNotFound
	}
	elem, ok := value.(*element)
	if !ok {
		return "", ErrInvalidElement
	}
	if !elem.expireAt.IsZero() && elem.expireAt.Before(time.Now()) {
		s.maps.Delete(key)
		s.elems.Put(elem)
		return "", ErrExpired
	}
	v := elem.value
	s.elems.Put(elem)
	return v, nil
}

func (s *mapCache) Set(ctx context.Context, key string, value string, expiration ...time.Duration) error {
	var expireAt time.Time
	if len(expiration) > 0 && expiration[0].Seconds() > 0 {
		expireAt = time.Now().Add(expiration[0])
	}
	elem := s.getElement()
	elem.value = value
	elem.expireAt = expireAt
	s.maps.Store(key, elem)
	return nil
}

func (s *mapCache) Delete(ctx context.Context, key string) error {
	if v, ok := s.maps.LoadAndDelete(key); ok {
		s.elems.Put(v)
		return nil
	}
	return ErrNotFound
}

func (s *mapCache) getElement() *element {
	return s.elems.Get().(*element)
}

func (s *mapCache) putElement(elem *element) {
	s.elems.Put(elem)
}

func NewMemoryCache() Cache {
	return &mapCache{
		elems: &sync.Pool{
			New: func() any {
				return &element{}
			},
		},
		maps: &sync.Map{},
	}
}

var _ Cache = (*mapCache)(nil)
