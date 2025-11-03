package blob

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	blobiface "github.com/origadmin/runtime/interfaces/storage/components/blob"
	layoutiface "github.com/origadmin/runtime/interfaces/storage/components/layout"
	"github.com/origadmin/runtime/storage/filestore/layout"
)

type blobStore struct {
	layout layoutiface.ShardedStorage // Use the concrete layout type
}

// New creates a new BlobStore implementation that uses a ShardedStorage for persistence.
func New(basePath string) (blobiface.Store, error) { // Accept basePath
	ls, err := layout.NewLocalShardedStorage(basePath) // Instantiate layout internally
	if err != nil {
		// Handle error: This should ideally be a panic or a fatal log,
		// as a failure to initialize the underlying storage is critical.
		// For now, we'll just panic for simplicity.
		return nil, fmt.Errorf("failed to create blob store: %w", err)
	}
	return &blobStore{
		layout: ls,
	}, nil
}

// Ensure blobStore implements the BlobStore interface.
var _ blobiface.Store = (*blobStore)(nil)

// Write calculates the SHA256 hash of the data and uses it as the ID.
// It then delegates the writing to the sharded layout manager.
func (s *blobStore) Write(data []byte) (string, error) {
	hashBytes := sha256.Sum256(data)
	hashString := hex.EncodeToString(hashBytes[:])

	err := s.layout.Write(hashString, data)
	if err != nil {
		return "", fmt.Errorf("failed to write blob to layout: %w", err)
	}
	return hashString, nil
}

// Read delegates reading to the sharded layout manager.
func (s *blobStore) Read(id string) ([]byte, error) {
	return s.layout.Read(id)
}

// Exists delegates existence check to the sharded layout manager.
func (s *blobStore) Exists(id string) (bool, error) {
	return s.layout.Exists(id)
}

// Delete delegates deletion to the sharded layout manager.
func (s *blobStore) Delete(id string) error {
	return s.layout.Delete(id)
}
