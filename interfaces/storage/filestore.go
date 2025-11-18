/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package storage defines the interfaces for storage services.
package storage

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	filev1 "github.com/origadmin/runtime/api/gen/go/config/data/file/v1"
)

var (
	// objectStoreBuilders is a global registry for ObjectStore builders.
	objectStoreBuilders = make(map[string]ObjectStoreBuilder)
	// mu protects the global registry.
	mu sync.RWMutex
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
	Get(ctx context.Context, path string) (io.ReadCloser, error)

	// Stat retrieves metadata about an object without fetching the object itself.
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
	New(cfg *filev1.FilestoreConfig) (ObjectStore, error)
	// Name returns the name of the builder (e.g., "local", "s3").
	Name() string
}

// RegisterObjectStore registers a new ObjectStoreBuilder.
// This function is typically called from the init() function of a storage provider package.
// If a builder with the same name is already registered, it will panic.
func RegisterObjectStore(b ObjectStoreBuilder) {
	mu.Lock()
	defer mu.Unlock()

	name := b.Name()
	if _, exists := objectStoreBuilders[name]; exists {
		panic(fmt.Sprintf("storage: ObjectStore builder named %s already registered", name))
	}
	objectStoreBuilders[name] = b
}

// GetObjectStoreBuilder retrieves a registered ObjectStoreBuilder by name.
// It returns the builder and true if found, otherwise nil and false.
func GetObjectStoreBuilder(name string) (ObjectStoreBuilder, bool) {
	mu.RLock()
	defer mu.RUnlock()
	b, ok := objectStoreBuilders[name]
	return b, ok
}
