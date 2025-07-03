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
	// MimeType returns the MIME type of the file content.
	MimeType() string
	// RefCount returns the reference count of the file content.
	RefCount() int32
	// BlobRef returns the IDs of the blobs that constitute this file's content.
	BlobRef() ([]string, error)
	// GetExtension is used to obtain extension properties related to the content of a file.
	GetExtension(key string) (interface{}, bool)
	// SetExtension is used to set extension properties related to the contents of the file.
	SetExtension(key string, value interface{})
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
