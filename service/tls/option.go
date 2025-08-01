/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package tls implements the functions, types, and interfaces for the module.
package tls

import (
	"crypto/tls"
)

type Option = func(*tls.Config)

func WithInsecureSkipVerify() Option {
	return func(c *tls.Config) {
		c.InsecureSkipVerify = true
	}
}
