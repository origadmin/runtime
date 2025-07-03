/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"time"

	meta "github.com/origadmin/runtime/interfaces/storage/meta" // Alias for clarity
	metav1 "github.com/origadmin/runtime/storage/meta/v1"
	metav2 "github.com/origadmin/runtime/storage/meta/v2"
)

const (
	// DefaultBlockSize is a general constant, not directly tied to meta.go's logic now.
	DefaultBlockSize = 1024 * 1024 // 1MB
)

// Meta 结构体现在管理文件内容的元数据。
// 它不再管理目录结构，也不再直接存储文件系统层面的属性（如文件名、权限）。
type Meta struct {
	blobStorage meta.BlobStore
	path        string                   // Base path for meta storage (e.g., where meta blobs are stored)
	files       map[string]meta.FileMeta // Stores file metadata by its full path (for in-memory simulation)
}

// New creates a new Meta instance.
func New(path string, blobStorage meta.BlobStore) (*Meta, error) {
	m := &Meta{
		path:        path,
		blobStorage: blobStorage,
		files:       make(map[string]meta.FileMeta),
	}
	return m, nil
}

// WriteFile writes content to a file, creating or updating its metadata.
// It handles small files by embedding content directly into FileMetaV2.
func (m *Meta) WriteFile(name string, r io.Reader, perm fs.FileMode) error {
	if !path.IsAbs(name) {
		return fmt.Errorf("path must be absolute: %s", name)
	}

	content, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	var fileMeta meta.FileMeta

	// Determine if content should be embedded (for FileMetaV2)
	if int64(len(content)) <= metav2.EmbeddedFileSizeThreshold {
		// Create a FileMetaV2 instance for embedded content
		v2Data := &metav2.FileMetaV2{
			Version:      metav2.Version,
			FileSize:     int64(len(content)),
			ModifyTime:   time.Now().Unix(),
			MimeType:     "application/octet-stream", // Default MIME type
			RefCount:     1,                          // Initial ref count
			EmbeddedData: content,
			BlobRef:      "", // No external blob reference for embedded data
		}
		fileMeta = &meta.FileMetaData[metav2.FileMetaV2]{
			Data: v2Data,
			Info: &metaFileInfo{meta: v2Data}, // Placeholder for FileInfo
		}
	} else {
		// Store content in blob storage for large files
		contentHash, err := m.blobStorage.Store(content)
		if err != nil {
			return err
		}

		// Create a FileMetaV2 instance for blob-referenced content
		v2Data := &metav2.FileMetaV2{
			Version:    metav2.Version,
			FileSize:   int64(len(content)),
			ModifyTime: time.Now().Unix(),
			MimeType:   "application/octet-stream", // Default MIME type
			RefCount:   1,                          // Initial ref count
			BlobRef:    contentHash,
			// BlockSize and BlockHashes would be populated here if chunking was handled by meta.go
			// For now, assuming blobStorage.Store handles it as a single blob.
		}
		fileMeta = &meta.FileMetaData[metav2.FileMetaV2]{
			Data: v2Data,
			Info: &metaFileInfo{meta: v2Data}, // Placeholder for FileInfo
		}
	}

	// Store file metadata directly by its full path (in-memory simulation)
	m.files[name] = fileMeta

	return nil
}

// Open opens a file for reading.
func (m *Meta) Open(name string) (fs.File, error) {
	if !path.IsAbs(name) {
		return nil, fmt.Errorf("path must be absolute: %s", name)
	}

	fileMeta, ok := m.files[name]
	if !ok {
		return nil, os.ErrNotExist
	}

	// Based on the version of the FileMeta, call the appropriate NewMetaFile function
	switch fileMeta.Version() {
	case metav1.Version:
		// We need to pass the underlying data to NewMetaFileV1
		if fm, ok := fileMeta.(*meta.FileMetaData[metav1.FileMetaV1]); ok {
			return metav1.NewMetaFileV1(m.blobStorage, fm)
		}
		return nil, fmt.Errorf("invalid FileMeta type for V1: %T", fileMeta)
	case metav2.Version:
		// We need to pass the underlying data to NewMetaFileV2
		if fm, ok := fileMeta.(*meta.FileMetaData[metav2.FileMetaV2]); ok {
			return metav2.NewMetaFileV2(m.blobStorage, fm)
		}
		return nil, fmt.Errorf("invalid FileMeta type for V2: %T", fileMeta)
	default:
		return nil, fmt.Errorf("unsupported file meta version: %d", fileMeta.Version())
	}
}

// Stat returns file information.
func (m *Meta) Stat(name string) (fs.FileInfo, error) {
	if !path.IsAbs(name) {
		return nil, fmt.Errorf("path must be absolute: %s", name)
	}

	fileMeta, ok := m.files[name]
	if !ok {
		return nil, os.ErrNotExist
	}

	// Return a metaFileInfo that wraps the FileMeta interface
	return &metaFileInfo{meta: fileMeta}, nil
}

// file implements fs.File for our custom file entries.
type file struct {
	storage meta.BlobStore
	meta    meta.FileMeta // Now uses the interface
	reader  io.Reader
	offset  int64
	closed  bool
}

func (f *file) Stat() (fs.FileInfo, error) {
	if f.closed {
		return nil, fs.ErrClosed
	}
	return &metaFileInfo{meta: f.meta}, nil
}

func (f *file) Read(p []byte) (int, error) {
	if f.closed {
		return 0, fs.ErrClosed
	}
	return f.reader.Read(p)
}

func (f *file) Close() error {
	f.closed = true
	return nil
}

// metaFileInfo implements fs.FileInfo for our custom file entries.
type metaFileInfo struct {
	meta meta.FileMeta // Now uses the interface
}

func (m metaFileInfo) Name() string {
	// Name is not part of FileMeta interface, it belongs to DirEntry.
	// For now, return a placeholder. This will be handled by Index module.
	return ""
}

func (m metaFileInfo) Size() int64 {
	return m.meta.Size()
}

func (m metaFileInfo) Mode() fs.FileMode {
	// Mode is not part of FileMeta interface, it belongs to DirEntry.
	// For now, return a placeholder. This will be handled by Index module.
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

// ReadDir is not supported for file-only meta.
func (m *Meta) ReadDir(p string) ([]fs.DirEntry, error) {
	return nil, fmt.Errorf("ReadDir not supported for file-only meta")
}

// Mkdir is not supported for file-only meta.
func (m *Meta) Mkdir(p string, perm fs.FileMode) error {
	return fmt.Errorf("Mkdir not supported for file-only meta")
}
