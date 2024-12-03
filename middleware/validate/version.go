/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package validate implements the functions, types, and interfaces for the module.
package validate

//go:generate stringer -type=Version version.go

// Version is the version of the module.
type Version int

const (
	V1 Version = 1
	V2 Version = 2
)
