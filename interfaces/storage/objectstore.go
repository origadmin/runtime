/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package storage defines the interfaces for storage services.
package storage

import (
	"context"
	"io"
	"time"

	ossv1 "github.com/origadmin/runtime/api/gen/go/config/data/oss/v1"
)

// ObjectInfo describes a stored object. It serves as the standard
// data transfer object for metadata across all storage backends.
type ObjectInfo struct {
	// Path is the full path (or key) of the object.
	Path string
	// Size is the size of the object in bytes.
	Size int64
	// ModTime is the last modification time of the object.
	ModTime time.Time
	// Metadata contains backend-specific metadata (e.g., ETag, Content-Type).
	// It is not intended for general use but for backend-specific logic.
	Metadata map[string]interface{}
}

// ListOptions provides options for the List operation.
type ListOptions struct {
	// Prefix allows listing objects that start with a specific prefix.
	Prefix string
	// Recursive determines if the listing should be recursive.
	// If false, it mimics a single directory listing by using a delimiter.
	Recursive bool
}

// ObjectStore defines a standard, universal interface for object storage systems.
// It focuses on the core capabilities common to most object stores (like S3)
// and local file systems, abstracting away implementation details like multipart uploads.
type ObjectStore interface {
	// Put uploads an object to the store.
	// The path is the unique key for the object.
	// The data is read from the provided io.Reader until EOF.
	// Implementations are responsible for handling large files efficiently and transparently.
	Put(ctx context.Context, path string, data io.Reader, size int64) (*ObjectInfo, error)

	// Get retrieves an object from the store.
	// It returns an io.ReadCloser which must be closed by the caller.
	// If the object does not exist, implementations should return an error
	// that can be checked with errors.Is(err, os.ErrNotExist) or a similar sentinel error.
	Get(ctx context.Context, path string) (io.ReadCloser, error)

	// Stat retrieves metadata about an object without fetching the object itself.
	// If the object does not exist, implementations should return an error
	// that can be checked with errors.Is(err, os.ErrNotExist) or a similar sentinel error.
	Stat(ctx context.Context, path string) (*ObjectInfo, error)

	// Delete removes an object from the store.
	Delete(ctx context.Context, path string) error

	// List returns a slice of object info for objects in the store, filtered by options.
	List(ctx context.Context, opts ListOptions) ([]*ObjectInfo, error)
}

// ObjectStoreBuilder defines the interface for building an ObjectStore instance.
// Each storage provider (e.g., local, s3) must implement this interface
// and register itself using the RegisterObjectStore function.
type ObjectStoreBuilder interface {
	// New builds a new ObjectStore instance from the given configuration.
	New(cfg *ossv1.ObjectStoreConfig) (ObjectStore, error)
}
