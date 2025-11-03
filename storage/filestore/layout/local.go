/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package layout

import (
	"fmt"
	"os"
	"path/filepath"

	layoutiface "github.com/origadmin/runtime/interfaces/storage/components/layout"
	"github.com/origadmin/runtime/storage/filestore/internal/fileutil"
)

// localShardedStorage implements ShardedStorage for the local filesystem.
type localShardedStorage struct {
	basePath string // The root directory (e.g., "/var/data/blobs")
}

// NewLocalShardedStorage creates a new instance.
func NewLocalShardedStorage(basePath string) (layoutiface.ShardedStorage, error) {
	// Ensure the base path exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}
	return &localShardedStorage{basePath: basePath}, nil
}

// idToPath is the core logic that converts an ID to a file path.
func (s *localShardedStorage) idToPath(id string) (string, error) {
	if len(id) < 4 {
		return "", fmt.Errorf("id '%s' is too short for sharding", id)
	}
	// Our xx/yy/rest logic
	return filepath.Join(s.basePath, id[:2], id[2:4], id), nil
}

func (s *localShardedStorage) Write(id string, data []byte) error {
	path, err := s.idToPath(id)
	if err != nil {
		return err
	}

	// Ensure the parent directory (e.g., /base/xx/yy/) exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return fileutil.AtomicWrite(path, data)
}

func (s *localShardedStorage) Read(id string) ([]byte, error) {
	path, err := s.idToPath(id)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

func (s *localShardedStorage) Exists(id string) (bool, error) {
	path, err := s.idToPath(id)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *localShardedStorage) Delete(id string) error {
	path, err := s.idToPath(id)
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func (s *localShardedStorage) GetPath(id string) (string, error) {
	return s.idToPath(id)
}
