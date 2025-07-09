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
