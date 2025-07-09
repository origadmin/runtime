/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package blob

import (
	"crypto/sha256"
	"encoding/hex"

	layoutiface "github.com/origadmin/runtime/interfaces/storage/layout"
)

// FileStore implements the BlobStore interface using the local filesystem.
// It relies on a ShardedStorage layout to manage the physical files.
type FileStore struct {
	layout layout.ShardedStorage
}

// NewFileStore creates a new FileStore.
func NewFileStore(ls layoutiface.ShardedStorage) (*FileStore, error) {
	return &FileStore{layout: ls}, nil
}

// Write calculates the SHA256 hash of the data and uses it as the ID.
// It then delegates the writing to the sharded layout manager.
func (s *FileStore) Write(data []byte) (string, error) {
	hashBytes := sha256.Sum256(data)
	hashString := hex.EncodeToString(hashBytes[:])

	err := s.layout.Write(hashString, data)
	if err != nil {
		return "", err
	}
	return hashString, nil
}

// Read delegates reading to the sharded layout manager.
func (s *FileStore) Read(id string) ([]byte, error) {
	return s.layout.Read(id)
}

// Exists delegates existence check to the sharded layout manager.
func (s *FileStore) Exists(id string) (bool, error) {
	return s.layout.Exists(id)
}

// Delete delegates deletion to the sharded layout manager.
func (s *FileStore) Delete(id string) error {
	return s.layout.Delete(id)
}
