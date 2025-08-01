/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package envars

import (
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
)

type envars struct {
	data []*config.KeyValue
}

func NewSource(prefixes ...string) config.Source {
	return &envars{
		data: loadEnviron(os.Environ(), prefixes),
	}
}

func (e *envars) Load() (kv []*config.KeyValue, err error) {
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

func (e *envars) Watch() (config.Watcher, error) {
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
