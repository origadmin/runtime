/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package metav2 implements the functions, types, and interfaces for the module.
package metav2

import (
	"io"
	"io/fs"
	"time"

	"github.com/origadmin/runtime/interfaces/storage/meta"
)

const Version = 2

type FileMetaV2 struct {
	Version  int32  `msgpack:"v"` // File meta version
	Size     int64  `msgpack:"s"` // File size
	ModTime  int64  `msgpack:"t"` // Modify time
	MimeType string `msgpack:"m"` // File mime type
	RefCount int32  `msgpack:"r"` // Reference count

	// if version > 1, then we have the following fields:
	BlockSize   int32    `msgpack:"bs"` // New field
	BlockHashes []string `msgpack:"bh"` // New field
}

func (f FileMetaV2) CurrentVersion() int32 {
	return Version
}

type FileMeta = meta.FileMetaData[FileMetaV2]

type file struct {
	storage meta.BlobStorage
	meta    *FileMeta
	reader  io.Reader
	offset  int64
	closed  bool
}

type metaFileInfo struct {
	meta *FileMeta
}

func (m metaFileInfo) Name() string {
	return m.meta.Info.Name
}

func (m metaFileInfo) Size() int64 {
	return m.meta.Data.Size
}

func (m metaFileInfo) Mode() fs.FileMode {
	return m.meta.Info.FileMode
}

func (m metaFileInfo) ModTime() time.Time {
	return time.Unix(m.meta.Data.ModTime, 0)
}

func (m metaFileInfo) IsDir() bool {
	return m.meta.Info.Hash == ""
}

func (m metaFileInfo) Sys() any {
	return m.meta.Info.Sys
}

func (f file) Stat() (fs.FileInfo, error) {
	if f.closed {
		return nil, fs.ErrClosed
	}
	return &metaFileInfo{meta: f.meta}, nil
}

func (f file) Read(p []byte) (int, error) {
	if f.closed {
		return 0, fs.ErrClosed
	}
	return f.reader.Read(p)
}

func (f file) Close() error {
	f.closed = true
	return nil
}

func NewMetaFileV2(storage meta.BlobStorage, fileMeta *FileMeta) (fs.File, error) {
	cr := &chunkReader{
		storage: storage,
		hashes:  fileMeta.Data.BlockHashes,
	}
	return &file{
		storage: storage,
		meta:    fileMeta,
		reader:  cr,
	}, nil
}
