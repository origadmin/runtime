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

func NewSource(conds ...string) config.Source {
	return &envars{
		data: loadEnviron(os.Environ(), conds),
	}
}

func (e *envars) Load() (kv []*config.KeyValue, err error) {
	return e.data, nil
}

func loadEnviron(data, conds []string) []*config.KeyValue {
	var kv []*config.KeyValue
	var ok bool
	for _, env := range data {
		var k, v string
		subs := strings.SplitN(env, "=", 2) //nolint:mnd
		k = subs[0]
		if len(subs) > 1 {
			v = subs[1]
		}

		k, ok = matchFold(conds, k)
		if ok && len(k) != 0 {
			kv = append(kv, &config.KeyValue{
				Key:   k,
				Value: []byte(v),
			})
		}
	}
	return kv
}

func (e *envars) Watch() (config.Watcher, error) {
	return env.NewWatcher()
}

func matchFold(envs []string, data string) (string, bool) {
	for _, env := range envs {
		if strings.EqualFold(env, data) {
			return data, true
		}
	}
	return "", false
}
