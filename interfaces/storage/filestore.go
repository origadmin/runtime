/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package storage implements the functions, types, and interfaces for the module.
package storage

import (
	"io"
	"time"
)

// FileInfo describes a file or directory. It serves as the standard
// data transfer object for metadata across all storage backends.
type FileInfo struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime time.Time
}

// CompletedPart holds information about a successfully uploaded part in a multipart upload.
type CompletedPart struct {
	// PartNumber is the sequence number of the part, starting from 1.
	PartNumber int
	// ETag is the entity tag returned by the storage backend for the uploaded part.
	// It's used to verify the part's integrity upon completion.
	ETag string
}

// MultipartUpload represents a stateful, in-progress multipart upload session,
// giving the caller full control over the upload lifecycle.
type MultipartUpload interface {
	// UploadPart uploads a single chunk of data as a part of the multipart upload.
	// The provided reader is consumed. partNumber must be between 1 and 10000.
	// The `size` parameter is a hint for the underlying driver, similar to the one in Write.
	UploadPart(partNumber int, reader io.Reader, size int64) (CompletedPart, error)

	// Complete finalizes the multipart upload, assembling the uploaded parts into a single object.
	// The parts must be sorted by PartNumber.
	Complete(parts []CompletedPart) error

	// Abort cancels the multipart upload, cleaning up any uploaded parts that have been
	// stored temporarily by the backend.
	Abort() error

	// UploadID returns the unique identifier for this upload session.
	UploadID() string
}

// FileStore defines a standard, universal interface for file and object storage systems.
// It provides both simple, one-shot operations and advanced, stateful multipart operations.
type FileStore interface {
	// --- Simple, Stateless Operations ---

	List(path string) ([]FileInfo, error)
	Stat(path string) (FileInfo, error)
	Read(path string) (io.ReadCloser, error)
	Mkdir(path string) error
	Delete(path string) error
	Rename(oldPath, newPath string) error

	// --- Write Operations ---

	// Write provides a simple, stream-based method for writing an object. It's ideal
	// for smaller files or when direct control over chunking is not required.
	// For large files, implementations should transparently handle multipart uploads internally.
	// The `size` parameter is a hint for the underlying driver. If the size is unknown, it should be -1.
	Write(path string, data io.Reader, size int64) error

	// InitiateMultipartUpload starts a new multipart upload session and returns a handler
	// for managing it. This is the entry point for advanced use cases like resumable
	// or parallel uploads.
	//InitiateMultipartUpload(path string) (MultipartUpload, error)
}
