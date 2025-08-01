/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package meta

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"

	layoutiface "github.com/origadmin/runtime/interfaces/storage/components/layout"
	metaiface "github.com/origadmin/runtime/interfaces/storage/components/meta"
	"github.com/origadmin/runtime/storage/filestore/layout"
	metav1 "github.com/origadmin/runtime/storage/filestore/meta/v1"
	metav2 "github.com/origadmin/runtime/storage/filestore/meta/v2"
)

// Store implements the metaiface.Store interface using the local filesystem.
// It relies on a ShardedStorage layout to manage the physical files.
type Store struct {
	layout layoutiface.ShardedStorage
}

// Ensure Store implements the metaiface.Store interface.
var _ metaiface.Store = (*Store)(nil)

// NewStore creates a new Store.
func NewStore(basePath string) (*Store, error) {
	ls, err := layout.NewLocalShardedStorage(basePath)
	if err != nil {
		return nil, err
	}
	return &Store{layout: ls}, nil
}

// Create serializes the FileMeta and stores it using a pre-determined, content-derived ID.
// The caller is responsible for generating the correct ID.
func (s *Store) Create(id string, fileMeta metaiface.FileMeta) error {
	// For new writes, we always expect the current version (V2) to be passed.
	v2Meta, ok := fileMeta.(*metav2.FileMetaV2)
	if !ok {
		return fmt.Errorf("unsupported FileMeta type for creation, expected *metav2.FileMetaV2, got %T", fileMeta)
	}

	// Serialize the metadata object for storage.
	fileMetaData := &metaiface.StoreMeta[metav2.FileMetaV2]{
		Version: metav2.Version,
		Data:    *v2Meta,
	}
	data, err := msgpack.Marshal(fileMetaData)
	if err != nil {
		return fmt.Errorf("failed to marshal FileMetaV2: %w", err)
	}

	// Write to the sharded layout using the content-derived ID provided by the caller.
	return s.layout.Write(id, data)
}

// Get retrieves and deserializes the FileMeta by its ID.
func (s *Store) Get(id string) (metaiface.FileMeta, error) {
	data, err := s.layout.Read(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read FileMeta from layout: %w", err)
	}

	var versionOnly metaiface.FileMetaVersion
	err = msgpack.Unmarshal(data, &versionOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal FileMeta version: %w", err)
	}

	switch versionOnly.Version {
	case metav1.Version:
		var fileMetaData metaiface.StoreMeta[metav1.FileMetaV1]
		err = msgpack.Unmarshal(data, &fileMetaData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal FileMetaV1: %w", err)
		}
		return fileMetaData.Data, nil
	case metav2.Version:
		var fileMetaData metaiface.StoreMeta[metav2.FileMetaV2]
		err = msgpack.Unmarshal(data, &fileMetaData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal FileMetaV2: %w", err)
		}
		return fileMetaData.Data, nil
	default:
		return nil, fmt.Errorf("unsupported FileMeta version: %d", versionOnly.Version)
	}
}

// Update serializes the updated FileMeta and overwrites the existing record.
func (s *Store) Update(id string, fileMeta metaiface.FileMeta) error {
	// For updates, we also assume we are working with the latest version format.
	// If an older version needs updating, it should be migrated first.
	v2Meta, ok := fileMeta.(*metav2.FileMetaV2)
	if !ok {
		return fmt.Errorf("unsupported FileMeta type for update, expected *metav2.FileMetaV2, got %T", fileMeta)
	}

	fileMetaData := &metaiface.StoreMeta[metav2.FileMetaV2]{
		Version: metav2.Version,
		Data:    *v2Meta,
	}

	data, err := msgpack.Marshal(fileMetaData)
	if err != nil {
		return fmt.Errorf("failed to marshal updated FileMeta: %w", err)
	}

	return s.layout.Write(id, data)
}

// Delete removes the FileMeta record.
func (s *Store) Delete(id string) error {
	return s.layout.Delete(id)
}

// Migrate is not yet implemented.
func (s *Store) Migrate(id string) (metaiface.FileMeta, error) {
	return nil, fmt.Errorf("method Migrate not implemented")
}

// CurrentVersion returns the version number used for writing new metadata.
// The system maintains backward compatibility for reading older, supported versions.
func (s *Store) CurrentVersion() int {
	return metav2.Version // Always write the latest version for new metadata.
}
