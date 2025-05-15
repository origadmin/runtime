/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package envsetup implements the functions, types, and interfaces for the module.
package envsetup

import (
	"os"
	"strings"
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
	for k, v := range env {
		err := os.Setenv(strings.Join([]string{prefix, k}, "_"), v)
		if err != nil {
			return err
		}
	}
	return nil
}
