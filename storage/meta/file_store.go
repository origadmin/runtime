/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package meta

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/origadmin/runtime/interfaces/storage/meta"
	"github.com/origadmin/runtime/storage/layout"
	metav1 "github.com/origadmin/runtime/storage/meta/v1"
	metav2 "github.com/origadmin/runtime/storage/meta/v2"
)

// FileMetaStore implements the MetaStore interface using the local filesystem.
// It relies on a ShardedStorage layout to manage the physical files.
type FileMetaStore struct {
	layout layout.ShardedStorage
}

// Ensure FileMetaStore implements the MetaStore interface.
var _ meta.Store = (*FileMetaStore)(nil)

// NewFileMetaStore creates a new FileMetaStore.
func NewFileMetaStore(basePath string) (*FileMetaStore, error) {
	ls, err := layout.NewLocalShardedStorage(basePath)
	if err != nil {
		return nil, err
	}
	return &FileMetaStore{layout: ls}, nil
}

// Create serializes the FileMeta and stores it.
// It returns the ID (hash) of the stored meta.
func (s *FileMetaStore) Create(fileMeta meta.FileMeta) (string, error) {
	var fileMetaData interface{} // Use interface{} to hold FileMetaData[V1] or FileMetaData[V2]

	switch fileMeta.CurrentVersion() {
	case metav1.Version:
		actualFileMeta, ok := fileMeta.(*metav1.FileMetaV1)
		if !ok {
			return "", fmt.Errorf("expected FileMetaV1 for version %d, got %T", metav1.Version, fileMeta)
		}
		fileMetaData = &metaiface.FileMetaData[metav1.FileMetaV1]{
			Version: metav1.Version,
			Data:    actualFileMeta,
		}
	case metav2.Version:
		actualFileMeta, ok := fileMeta.(*metav2.FileMetaV2)
		if !ok {
			return "", fmt.Errorf("expected FileMetaV2 for version %d, got %T", metav2.Version, fileMeta)
		}
				fileMetaData = &metaiface.FileMetaData[metav2.FileMetaV2]{
			Version: metav2.Version,
			Data:    actualFileMeta,
		}
	default:
		return "", fmt.Errorf("unsupported FileMeta version for creation: %d", fileMeta.CurrentVersion())
	}

	data, err := msgpack.Marshal(fileMetaData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal FileMeta: %w", err)
	}

	id := calculateHash(data)

	err = s.layout.Write(id, data)
	if err != nil {
		return "", fmt.Errorf("failed to write FileMeta to layout: %w", err)
	}

	return id, nil
}

// Get retrieves and deserializes the FileMeta by its ID.
func (s *FileMetaStore) Get(id string) (metaiface.FileMeta, error) {
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
		var fileMetaData metaiface.FileMetaData[metav1.FileMetaV1]
		err = msgpack.Unmarshal(data, &fileMetaData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal FileMetaV1: %w", err)
		}
		return fileMetaData.Data, nil
	case metav2.Version:
		var fileMetaData metaiface.FileMetaData[metav2.FileMetaV2]
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
func (s *FileMetaStore) Update(id string, fileMeta metaiface.FileMeta) error {
	var fileMetaData interface{}

	switch fileMeta.CurrentVersion() {
	case metav1.Version:
		actualFileMeta, ok := fileMeta.(*metav1.FileMetaV1)
		if !ok {
			return fmt.Errorf("expected FileMetaV1 for version %d, got %T", metav1.Version, fileMeta)
		}
		fileMetaData = &metaiface.FileMetaData[metav1.FileMetaV1]{
			Version: metav1.Version,
			Data:    actualFileMeta,
		}
	case metav2.Version:
		actualFileMeta, ok := fileMeta.(*metav2.FileMetaV2)
		if !ok {
			return fmt.Errorf("expected FileMetaV2 for version %d, got %T", metav2.Version, fileMeta)
		}
				fileMetaData = &metaiface.FileMetaData[metav2.FileMetaV2]{
			Version: metav2.Version,
			Data:    actualFileMeta,
		}
	default:
		return fmt.Errorf("unsupported FileMeta version for update: %d", fileMeta.CurrentVersion())
	}

	data, err := msgpack.Marshal(fileMetaData)
	if err != nil {
		return fmt.Errorf("failed to marshal updated FileMeta: %w", err)
	}

	return s.layout.Write(id, data)
}

// Delete removes the FileMeta record.
func (s *FileMetaStore) Delete(id string) error {
	return s.layout.Delete(id)
}

// BatchGet is not yet implemented.
func (s *FileMetaStore) BatchGet(ids []string) (map[string]meta.FileMeta, error) {
	return nil, fmt.Errorf("method BatchGet not implemented")
}

// Migrate is not yet implemented.
func (s *FileMetaStore) Migrate(id string) (meta.FileMeta, error) {
	return nil, fmt.Errorf("method Migrate not implemented")
}

// SupportedVersions returns the supported meta versions.
func (s *FileMetaStore) SupportedVersions() []int {
	return []int{metav1.Version, metav2.Version}
}

// DefaultVersion returns the default meta version.
func (s *FileMetaStore) DefaultVersion() int {
	return metav2.Version // Default to the latest version
}

// calculateHash is a helper to generate a hash for the meta record.
// This should ideally be consistent with how BlobStore generates its IDs.
func calculateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
