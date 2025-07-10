/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"time"

	blobiface "github.com/origadmin/runtime/interfaces/storage/blob"
	contentiface "github.com/origadmin/runtime/interfaces/storage/content"
	metaiface "github.com/origadmin/runtime/interfaces/storage/meta"
	metav2 "github.com/origadmin/runtime/storage/meta/v2"
)

const (
	// DefaultBlockSize is a general constant, not directly tied to meta.go's logic now.
	DefaultBlockSize = 4 * 1024 * 1024 // 4MB
)

// chunkData reads from the reader and splits the content into blocks of a fixed size.
// For each block, it calls the provided store function.
func (m *Meta) chunkData(r io.Reader) ([]string, int64, error) {
	var hashes []string
	var totalSize int64
	buf := make([]byte, m.chunkSize)

	for {
		n, err := io.ReadFull(r, buf)
		if err != nil && err != io.EOF && !errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, 0, err
		}
		if n == 0 {
			break
		}

		data := buf[:n]
		hash, storeErr := m.blobStorage.Write(data)
		if storeErr != nil {
			return nil, 0, storeErr
		}
		hashes = append(hashes, hash)
		totalSize += int64(n)

		if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
			break
		}
	}

	return hashes, totalSize, nil
}

// Meta 结构体现在管理文件内容的元数据。
// 它不再管理目录结构，也不再直接存储文件系统层面的属性（如文件名、权限）。
type Meta struct {
	metaStore   metaiface.Store
	blobStorage blobiface.Store
	assembler   contentiface.Assembler
	files       map[string]metaiface.FileMeta // Stores metaFile metadata by its full path (for in-memory simulation)
	chunkSize   int64                         // Configurable chunk size for writing large files
}

// New creates a new Meta instance.
func New(metaStore metaiface.Store, blobStorage blobiface.Store, assembler contentiface.Assembler, chunkSize int64) (*Meta, error) {
	m := &Meta{
		metaStore:   metaStore,
		blobStorage: blobStorage,
		assembler:   assembler,
		files:       make(map[string]metaiface.FileMeta),
		chunkSize:   chunkSize,
	}
	return m, nil
}

// WriteFile writes content to a metaFile, creating or updating its metadata.
// It handles small files by embedding content directly into FileMetaV2.
func (m *Meta) WriteFile(name string, r io.Reader, perm fs.FileMode) error {
	// Path validation should ideally happen at a higher layer (e.g., Index)
	// if !path.IsAbs(name) {
	// 	return fmt.Errorf("path must be absolute: %s", name)
	// }

	var fileMeta metaiface.FileMeta

	// Use a bytes.Buffer to peek at the content and determine if it's embedded or sharded.
	// This allows us to read the content once.
	buf := new(bytes.Buffer)
	teeReader := io.TeeReader(r, buf) // Read from r, and also write to buf

	// Try to read up to the embedded size threshold + 1 byte to determine if it's larger
	_, err := io.CopyN(io.Discard, teeReader, metav2.EmbeddedFileSizeThreshold+1)

	// The actual content read so far is in buf.Bytes()
	contentBytes := buf.Bytes()

	// Determine if content should be embedded (for FileMetaV2)
	if int64(len(contentBytes)) <= metav2.EmbeddedFileSizeThreshold && err == io.EOF { // Check for EOF to ensure it's truly small
		// It's a small metaFile, embed the content
		_ = &metav2.FileMetaV2{
			Version:      metav2.Version,
			FileSize:     int64(len(contentBytes)),
			ModifyTime:   time.Now().Unix(),
			MimeType:     "application/octet-stream", // Default MIME type
			RefCount:     1,                          // Initial ref count
			EmbeddedData: contentBytes,
			BlobHashes:   nil, // No external blob reference for embedded data
		}
		//fileMeta = &metaiface.FileMetaData[metav2.FileMetaV2]{
		//	Data: v2Data,
		//	Info: &metaFileInfo{meta: v2Data}, // Placeholder for FileInfo
		//}
	} else {
		// It's a large metaFile, or we couldn't determine size easily, so chunk the rest of the stream.
		// The teeReader has already consumed some bytes, so chunkData will continue from there.
		hashes, totalSize, chunkErr := m.chunkData(io.MultiReader(bytes.NewReader(contentBytes), r))
		if chunkErr != nil {
			return fmt.Errorf("failed to chunk and store data: %w", chunkErr)
		}

		// Create a FileMetaV2 instance for blob-referenced content
		_ = &metav2.FileMetaV2{
			Version:    metav2.Version,
			FileSize:   totalSize,
			ModifyTime: time.Now().Unix(),
			MimeType:   "application/octet-stream", // Default MIME type
			RefCount:   1,                          // Initial ref count
			BlobHashes: hashes,
		}
		//fileMeta = &metaiface.FileMetaData[metav2.FileMetaV2]{
		//	Data: v2Data,
		//	Info: &metaFileInfo{meta: v2Data}, // Placeholder for FileInfo
		//}
	}

	// Store metaFile metadata using the injected MetaStore
	_, err = m.metaStore.Create(fileMeta)
	if err != nil {
		return fmt.Errorf("failed to create metaFile meta in store: %w", err)
	}

	// Store metaFile metadata directly by its full path (in-memory simulation) - This part might be removed later if Index handles it
	m.files[name] = fileMeta

	return nil

}

// Open opens a metaFile for reading.
func (m *Meta) Open(name string) (fs.File, error) {
	// Path validation should ideally happen at a higher layer (e.g., Index)
	// if !path.IsAbs(name) {
	// 	return nil, fmt.Errorf("path must be absolute: %s", name)
	// }

	fileMeta, err := m.metaStore.Get(name) // Assuming name can be used as ID
	if err != nil {
		return nil, err
	}

	reader, err := m.assembler.NewReader(fileMeta)
	if err != nil {
		return nil, err
	}

	return &metaFile{
		meta:   fileMeta,
		reader: reader,
	}, nil
}

// Stat returns metaFile information.
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

// metaFile implements fs.File for our custom metaFile entries.
type metaFile struct {
	storage blobiface.Store
	meta    metaiface.FileMeta // Now uses the interface
	reader  io.Reader
	offset  int64
	closed  bool
}

func (f *metaFile) Stat() (fs.FileInfo, error) {
	if f.closed {
		return nil, fs.ErrClosed
	}
	return &metaFileInfo{meta: f.meta}, nil
}

func (f *metaFile) Read(p []byte) (int, error) {
	if f.closed {
		return 0, fs.ErrClosed
	}
	return f.reader.Read(p)
}

func (f *metaFile) Close() error {
	f.closed = true
	return nil
}

// metaFileInfo implements fs.FileInfo for our custom file entries.
type metaFileInfo struct {
	meta metaiface.FileMeta // Now uses the interface
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
	return false // FileMeta always represents a metaFile, not a directory
}

func (m metaFileInfo) Sys() any {
	return nil
}

// ReadDir is not supported for metaFile-only meta.
func (m *Meta) ReadDir(p string) ([]fs.DirEntry, error) {
	return nil, fmt.Errorf("ReadDir not supported for metaFile-only meta")
}

// Mkdir is not supported for metaFile-only meta.
func (m *Meta) Mkdir(p string, perm fs.FileMode) error {
	return fmt.Errorf("Mkdir not supported for file-only meta")
}

// StartUpload initiates a new file upload session for chunked uploads.
func (m *Meta) StartUpload(fileName string, totalSize int64, mimeType string) (uploadID string, err error) {
	// Placeholder implementation
	return "", fmt.Errorf("StartUpload not yet implemented")
}

// UploadChunk uploads a single chunk of a file.
func (m *Meta) UploadChunk(uploadID string, chunkIndex int, chunkData []byte) (err error) {
	// Placeholder implementation
	return fmt.Errorf("UploadChunk not yet implemented")
}

// FinishUpload finalizes a file upload session.
func (m *Meta) FinishUpload(uploadID string) (fileMetaID string, err error) {
	// Placeholder implementation
	return "", fmt.Errorf("FinishUpload not yet implemented")
}

// CancelUpload cancels a file upload session and cleans up temporary data.
func (m *Meta) CancelUpload(uploadID string) (err error) {
	// Placeholder implementation
	return fmt.Errorf("CancelUpload not yet implemented")
}
