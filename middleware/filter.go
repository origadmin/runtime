/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"strings"
)

type Filter interface {
	Filter(key string, m KMiddleware) Filter
	All() []KMiddleware
	Filtered() []KMiddleware
	Total() int
}

type filter struct {
	keys     []string
	filtered []KMiddleware
	all      []KMiddleware
}

func (f *filter) Total() int {
	return len(f.all)
}

func (f *filter) All() []KMiddleware {
	return f.all
}

func (f *filter) Filtered() []KMiddleware {
	return f.filtered
}

func (f *filter) Filter(key string, m KMiddleware) Filter {
	f.all = append(f.all, m)
	if len(f.keys) == 0 {
		return f
	}
	var kee string
	for _, kee = range f.keys {
		if strings.EqualFold(kee, key) {
			f.filtered = append(f.filtered, m)
			return f
		}
	}
	return f
}

func MakeFilter(keys []string) Filter {
	return &filter{
		keys: keys,
	}
}
