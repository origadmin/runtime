/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

type FileMeta interface {
	CurrentVersion() int32
}

type BlobStorage interface {
	Store(content []byte) (string, error)
	Retrieve(hash string) ([]byte, error)
}
