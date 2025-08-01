/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package metav2 implements the functions, types, and interfaces for the module.
package metav2

import (
	"time"
)

// Version defines the version of the metadata format.
const Version = 2

// EmbeddedFileSizeThreshold Defines the maximum size of metadata for small files to be embedded directly (256KB)
const EmbeddedFileSizeThreshold = 256 * 1024

// FileMetaV2 represents the metadata of a file.
type FileMetaV2 struct {
	FileSize   int64  `msgpack:"s"` // File size
	ModifyTime int64  `msgpack:"t"` // Modify time
	MimeType   string `msgpack:"m"` // File mime type
	RefCount   int32  `msgpack:"r"` // Reference count

	BlobSize   int32    `msgpack:"bs"` // Blob size
	BlobHashes []string `msgpack:"bh"` // Reference to the blob content

	EmbeddedData []byte `msgpack:"ed,omitempty"` // Used to store file content that is less than the EmbeddedFileSizeThreshold
}

func (f FileMetaV2) CurrentVersion() int32 {
	return Version
}

func (f FileMetaV2) Size() int64 {
	return f.FileSize
}

func (f FileMetaV2) ModTime() time.Time {
	return time.Unix(f.ModifyTime, 0)
}

// GetEmbeddedData returns the embedded data of the file.
func (f FileMetaV2) GetEmbeddedData() []byte {
	return f.EmbeddedData
}

// GetShards returns the blob hashes (shards) of the file.
func (f FileMetaV2) GetShards() []string {
	return f.BlobHashes
}
