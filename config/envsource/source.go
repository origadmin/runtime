/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package envsource is a configuration source that loads environment variables.
package envsource

import (
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"

	sourcev1 "github.com/origadmin/runtime/api/gen/go/config/source/v1"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/interfaces/options"
)

type source struct {
	data []*config.KeyValue
}

func NewSource(prefixes ...string) config.Source {
	return &source{
		data: loadEnviron(os.Environ(), prefixes),
	}
}

func (e *source) Load() (kv []*config.KeyValue, err error) {
	return e.data, nil
}

func loadEnviron(data, prefixes []string) []*config.KeyValue {
	var ok bool
	kvs := make([]*config.KeyValue, 0)
	var k, v, prefix string
	for _, datum := range data {
		k, v, _ = strings.Cut(datum, "=") //nolint:mnd
		if len(prefixes) > 0 {
			prefix, ok = matchPrefix(prefixes, k)
			if !ok || len(prefix) == len(k) {
				continue
			}
			// trim prefix
			k = strings.TrimPrefix(k, prefix)
			k = strings.TrimPrefix(k, "_")
		}

		if len(k) > 0 {
			kvs = append(kvs, &config.KeyValue{
				Key:   k,
				Value: []byte(v),
			})
		}
	}
	return kvs
}

func (e *source) Watch() (config.Watcher, error) {
	return env.NewWatcher()
}

func matchPrefix(prefixes []string, v string) (string, bool) {
	for _, prefix := range prefixes {
		if strings.HasPrefix(v, prefix) {
			return prefix, true
		}
	}
	return "", false
}

func NewEnvSource(sourceCfg *sourcev1.SourceConfig, opts ...options.Option) (runtimeconfig.KSource, error) {
	envSrc := sourceCfg.GetEnv()
	prefixes := FromOptions(opts...)
	if envSrc == nil {
		// This can happen if the source type is "file" but the `file` oneof is not set.
		// Returning nil, nil is a safe default, allowing other sources to proceed.
		return NewSource(prefixes...), nil
	}
	if len(prefixes) == 0 {
		prefixes = envSrc.GetPrefixes()
	}
	return NewSource(prefixes...), nil
}

func init() {
	runtimeconfig.RegisterSourceFactory("env", runtimeconfig.SourceFunc(NewEnvSource))
}
