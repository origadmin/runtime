/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package metav1 implements the functions, types, and interfaces for the module.
package metav1

import (
	"time"
)

// Version is the current version of the file meta.
const Version = 1

// FileMetaV1 represents the file meta information.
type FileMetaV1 struct {
	FileSize   int64  `msgpack:"s"` // File size
	MimeType   string `msgpack:"m"` // MIME type
	ModifyTime int64  `msgpack:"t"` // Modify time
	RefCount   int32  `msgpack:"r"` // Reference count
}

func (f FileMetaV1) CurrentVersion() int32 {
	return Version
}

func (f FileMetaV1) Size() int64 {
	return f.FileSize
}

func (f FileMetaV1) ModTime() time.Time {
	return time.Unix(f.ModifyTime, 0)
}

// GetEmbeddedData returns nil as V1 does not support embedded data.
func (f FileMetaV1) GetEmbeddedData() []byte {
	return nil
}

// GetShards returns nil as V1 does not have a concept of sharded blobs in its structure.
func (f FileMetaV1) GetShards() []string {
	return nil
}
