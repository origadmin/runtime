/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"strings"
)

type Filter interface {
	Filter(key string, m Middleware) Filter
	All() []Middleware
	Filtered() []Middleware
	Total() int
}

type filter struct {
	keys     []string
	filtered []Middleware
	all      []Middleware
}

func (f *filter) Total() int {
	return len(f.all)
}

func (f *filter) All() []Middleware {
	return f.all
}

func (f *filter) Filtered() []Middleware {
	return f.filtered
}

func (f *filter) Filter(key string, m Middleware) Filter {
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
