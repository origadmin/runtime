/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package metav1 implements the functions, types, and interfaces for the module.
package metav1

import (
	"bytes"
	"io"
	"io/fs"
	"time"

	"github.com/origadmin/runtime/interfaces/storage/meta"
)

const Version = 1

type FileMetaV1 struct {
	Version   int32  `msgpack:"v"`   // File meta version
	Size_     int64  `msgpack:"s"`   // File size (renamed to avoid conflict with method)
	MimeType_ string `msgpack:"m"`   // MIME type (renamed to avoid conflict with method)
	ModTime_  int64  `msgpack:"t"`   // Modify time (renamed to avoid conflict with method)
	RefCount_ int32  `msgpack:"r"`   // Reference count (renamed to avoid conflict with method)
	BlobID    string `msgpack:"bid"` // Reference to the blob content

	extensions map[string]interface{} `msgpack:"ext,omitempty"` // Extension properties
}

func (f *FileMetaV1) CurrentVersion() int32 {
	return Version
}

func (f *FileMetaV1) Size() int64 {
	return f.Size_
}

func (f *FileMetaV1) ModTime() time.Time {
	return time.Unix(0, f.ModTime_)
}

func (f *FileMetaV1) MimeType() string {
	return f.MimeType_
}

func (f *FileMetaV1) RefCount() int32 {
	return f.RefCount_
}

func (f *FileMetaV1) BlobRef() ([]string, error) {
	return []string{f.BlobID}, nil
}

func (f *FileMetaV1) GetExtension(key string) (interface{}, bool) {
	if f.extensions == nil {
		return nil, false
	}
	v, ok := f.extensions[key]
	return v, ok
}

func (f *FileMetaV1) SetExtension(key string, value interface{}) {
	if f.extensions == nil {
		f.extensions = make(map[string]interface{})
	}
	f.extensions[key] = value
}

type file struct {
	storage meta.BlobStore
	meta    meta.FileMeta // Changed to interface
	reader  io.Reader
	offset  int64
	closed  bool
}

type metaFileInfo struct {
	meta meta.FileMeta // Changed to interface
}

func (m metaFileInfo) Name() string {
	// Name is not part of FileMeta interface, it belongs to DirEntry.
	// This will be handled by Index module.
	return ""
}

func (m metaFileInfo) Size() int64 {
	return m.meta.Size()
}

func (m metaFileInfo) Mode() fs.FileMode {
	// Mode is not part of FileMeta interface, it belongs to DirEntry.
	// This will be handled by Index module.
	return 0
}

func (m metaFileInfo) ModTime() time.Time {
	return m.meta.ModTime()
}

func (m metaFileInfo) IsDir() bool {
	return false // FileMeta always represents a file, not a directory
}

func (m metaFileInfo) Sys() any {
	return nil
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

func NewMetaFileV1(storage meta.BlobStore, fileMeta meta.FileMeta) (fs.File, error) {
	// We need to cast fileMeta to the concrete FileMetaData[FileMetaV1] type
	// to access its Data field (which is FileMetaV1) and Info field.
	// This is a temporary workaround until the MetaStore handles versioning and casting.
	// For now, we'll assume it's always the correct type for V1.
	actualFileMeta, err := fileMeta.BlobRef()
	if err != nil {
		return nil, err
	}
	data, err := storage.Read(actualFileMeta[0]) // Use BlobRef
	if err != nil {
		return nil, err
	}
	return &file{
		storage: storage,
		reader:  bytes.NewReader(data),
	}, nil
}
