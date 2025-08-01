/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package pagination implements the functions, types, and interfaces for the module.
package pagination

import (
	"github.com/goexts/generic/settings"
	"github.com/goexts/generic/types"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 15
	MaxPage         = 100
	MaxPageSize     = 1000
)

type PageLimiter interface {
	Current(in int32) int32
	PerPage(in int32) int32
}

type Limiter struct {
	MaxPage         int32
	DefaultPage     int32
	MaxPageSize     int32
	DefaultPageSize int32
}

type LimiterOption = func(*Limiter)

func WithMaxPageSize(size int32) LimiterOption {
	return func(l *Limiter) {
		l.MaxPageSize = size
	}
}

func WithMaxPage(page int32) LimiterOption {
	return func(l *Limiter) {
		l.MaxPage = page
	}
}

func WithDefaultPageSize(size int32) LimiterOption {
	return func(l *Limiter) {
		l.DefaultPageSize = size
	}
}

func WithDefaultPage(page int32) LimiterOption {
	return func(l *Limiter) {
		l.DefaultPage = page
	}
}

func (obj Limiter) Current(in int32) int32 {
	in = types.ZeroOr(in, obj.DefaultPage)
	if in > obj.MaxPage {
		in = obj.DefaultPage
	}
	return in
}

func (obj Limiter) PerPage(in int32) int32 {
	in = types.ZeroOr(in, obj.DefaultPageSize)
	if in > (obj.MaxPageSize) {
		in = obj.MaxPageSize
	}
	return in
}

func NewLimiter(page, defaultPage, pageSize, defaultPageSize int32) *Limiter {
	return settings.ApplyWithZero(
		WithMaxPageSize(pageSize),
		WithMaxPage(page),
		WithDefaultPageSize(defaultPageSize),
		WithDefaultPage(defaultPage))
}

func DefaultLimiter() *Limiter {
	return NewLimiter(MaxPage, DefaultPage, MaxPageSize, DefaultPageSize)
}

var _ PageLimiter = (*Limiter)(nil)
