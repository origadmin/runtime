/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
)

type Resolver interface {
	Resolve(config KConfig) (Resolved, error)
}

type Resolved interface {
	WithDecode(name string, v any, decode func([]byte, any) error) error
	Value(name string) (any, error)
	Middleware() *middlewarev1.Middleware
	Services() []*configv1.Service
	Logger() *configv1.Logger
	Discovery() *configv1.Discovery
}

type ResolveObserver interface {
	Observer(string, KValue)
}

type ResolveFunc func(config KConfig) (Resolved, error)

func (r ResolveFunc) Resolve(config KConfig) (Resolved, error) {
	return r(config)
}

type resolver struct {
	values map[string]any
}

func (r resolver) Services() []*configv1.Service {
	var ss []*configv1.Service
	if r.loadConfig("service", &ss) {
		return ss
	}
	return nil
}

func (r resolver) Discovery() *configv1.Discovery {
	v, ok := r.values["discovery"]
	if !ok {
		return nil
	}
	var discovery configv1.Discovery
	err := mapstructure.Decode(v, &discovery)
	if err != nil {
		return nil
	}
	return &discovery
}

func (r resolver) WithDecode(name string, v any, decode func([]byte, any) error) error {
	if v == nil {
		return fmt.Errorf("value %s is nil", name)
	}
	data, err := r.Value(name)
	if err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("value %s is nil", name)
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return decode(marshal, v)
}

func (r resolver) Value(name string) (any, error) {
	v, ok := r.values[name]
	if !ok {
		return nil, fmt.Errorf("value %s not found", name)
	}
	return v, nil
}

func (r resolver) Middleware() *middlewarev1.Middleware {
	var m middlewarev1.Middleware
	if r.loadConfig("middleware", &m) {
		return &m
	}
	return nil
}

func (r resolver) Logger() *configv1.Logger {
	var l configv1.Logger
	if r.loadConfig("logger", &l) {
		return &l
	}
	return nil
}

func (r resolver) loadConfig(key string, target interface{}) bool {
	v, ok := r.values[key]
	if !ok {
		return false
	}
	return mapstructure.Decode(v, target) == nil
}

var DefaultResolver Resolver = ResolveFunc(func(config KConfig) (Resolved, error) {
	var r resolver
	err := config.Scan(&r.values)
	if err != nil {
		return nil, err
	}
	return r, nil
})
