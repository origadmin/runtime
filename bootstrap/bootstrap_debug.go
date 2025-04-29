//go:build debug

/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package bootstrap is a package that provides the bootstrap information for the service.
package bootstrap

func init() {
	buildEnv = EnvDebug
}
