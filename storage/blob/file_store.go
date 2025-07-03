/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package blob

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/origadmin/runtime/storage/layout"
)

// FileBlobStore implements the BlobStore interface using the local filesystem.
// It relies on a ShardedStorage layout to manage the physical files.
type FileBlobStore struct {
	layout layout.ShardedStorage
}

// NewFileBlobStore creates a new FileBlobStore.
func NewFileBlobStore(basePath string) (*FileBlobStore, error) {
	// Create the sharded storage layout manager
	ls, err := layout.NewLocalShardedStorage(basePath)
	if err != nil {
		return nil, err
	}
	return &FileBlobStore{layout: ls}, nil
}

// Write calculates the SHA256 hash of the data and uses it as the ID.
// It then delegates the writing to the sharded layout manager.
func (s *FileBlobStore) Write(data []byte) (string, error) {
	hashBytes := sha256.Sum256(data)
	hashString := hex.EncodeToString(hashBytes[:])

	err := s.layout.Write(hashString, data)
	if err != nil {
		return "", err
	}
	return hashString, nil
}

// Read delegates reading to the sharded layout manager.
func (s *FileBlobStore) Read(id string) ([]byte, error) {
	return s.layout.Read(id)
}

// Exists delegates existence check to the sharded layout manager.
func (s *FileBlobStore) Exists(id string) (bool, error) {
	return s.layout.Exists(id)
}

// Delete delegates deletion to the sharded layout manager.
func (s *FileBlobStore) Delete(id string) error {
	return s.layout.Delete(id)
}
