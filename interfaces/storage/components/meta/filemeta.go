/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"time"
)

type FileMeta interface {
	// CurrentVersion returns the version number of this metadata record.
	CurrentVersion() int32
	// Size returns the byte size of the file contents.
	Size() int64
	// ModTime returns when the contents of the file itself were last modified.
	ModTime() time.Time
	// GetEmbeddedData returns the raw byte data if the file content is embedded directly in the metadata.
	// It returns nil if the content is not embedded.
	GetEmbeddedData() []byte
	// GetShards returns a list of blob hashes if the file content is stored in external blobs (shards).
	// It returns nil if the content is embedded or the file is empty.
	GetShards() []string
}

type FileMetaVersion struct {
	Version int32 `msgpack:"v"`
}

func (f FileMetaVersion) CurrentVersion() int32 {
	return f.Version
}

type FileMetaData[T any] struct {
	Version int32 `json:"version" msgpack:"v"`
	Data    *T    `json:"data" msgpack:"d"`
}
