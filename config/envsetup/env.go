/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package envsetup implements the functions, types, and interfaces for the module.
package envsetup

import (
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

func Set(env map[string]string) error {
	for k, v := range env {
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetWithPrefix(prefix string, env map[string]string) error {
	ll := log.NewHelper(log.With(log.GetLogger(), "module", "envsetup"))
	for k, v := range env {
		pk := strings.Join([]string{prefix, k}, "_")
		//pk = strings.ToUpper(pk)
		ll.Infow("msg", "set env", "key", pk, "value", v)
		err := os.Setenv(pk, v)
		if err != nil {
			return err
		}
	}
	return nil
}
