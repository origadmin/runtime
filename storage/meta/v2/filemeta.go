/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package metav2 implements the functions, types, and interfaces for the module.
package metav2

import (
	"bytes"
	"io"
	"time"

	"github.com/origadmin/runtime/interfaces/storage/blob"
)

const Version = 2

// EmbeddedFileSizeThreshold 定义了小文件直接嵌入元数据的最大大小 (256KB)
const EmbeddedFileSizeThreshold = 256 * 1024

type FileMetaV2 struct {
	Version    int32  `msgpack:"v"` // File meta version
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

func (f FileMetaV2) ContentReader(storage blob.BlobStore) (io.Reader, error) {
	if len(f.EmbeddedData) > 0 {
		return bytes.NewReader(f.EmbeddedData), nil
	} else {
		return &chunkReader{
				storage: storage,
				hashes:  f.BlobHashes,
			},
			nil
	}
}
