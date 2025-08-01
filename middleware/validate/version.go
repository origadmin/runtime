/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package validate implements the functions, types, and interfaces for the module.
package validate

//go:generate stringer -type=Version version.go

// Version is the version of the module.
type Version int

const (
	// V1 is the first version of the module, used github.com/envoyproxy/protoc-gen-validate
	V1 Version = 1
	// V2 is the second version of the module, used buf.build/go/protovalidate
	V2 Version = 2
)
