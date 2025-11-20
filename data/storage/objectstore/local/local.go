/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package local

import (
	"io"
	"os"
	"path/filepath"
	"strings" // Added import

	ossv1 "github.com/origadmin/runtime/api/gen/go/config/data/oss/v1"
	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/data/storage/objectstore"
	runtimeerrors "github.com/origadmin/runtime/errors"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
)

const (
	DriverName = "local"
)

// localFactory implements the objectstore.Factory interface for the local object store.
type localFactory struct{}

func init() {
	// Register the local object store factory.
	objectstore.Register(DriverName, &localFactory{})
}

// New creates a new local object store.
func (f *localFactory) New(cfg *ossv1.ObjectStoreConfig) (storageiface.ObjectStore, error) {
	if cfg == nil || cfg.GetLocal() == nil {
		return nil, runtimeerrors.NewStructured(objectstore.Module, "object store config is nil").WithCaller()
	}
	localCfg := cfg.GetLocal()
	if localCfg == nil || localCfg.GetRoot() == "" {
		return nil, runtimeerrors.NewStructured(objectstore.Module, "root directory for local objectstore is not configured").WithCaller()
	}
	if err := os.MkdirAll(localCfg.GetRoot(), 0750); err != nil {
		return nil, runtimeerrors.WrapStructured(err, objectstore.Module, "failed to create root directory")
	}
	return &localStore{rootDir: localCfg.GetRoot()}, nil
}

// localStore implements the storageiface.ObjectStore for the local filesystem.
type localStore struct {
	rootDir string
}

func (l *localStore) Put(ctx context.Context, path string, data io.Reader, size int64) (*storageiface.ObjectInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (l *localStore) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	//TODO implement me
	panic("implement me")
}

func (l *localStore) Stat(ctx context.Context, path string) (*storageiface.ObjectInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (l *localStore) Delete(ctx context.Context, path string) error {
	//TODO implement me
	panic("implement me")
}

func (l *localStore) List(ctx context.Context, opts storageiface.ListOptions) ([]*storageiface.ObjectInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (l *localStore) resolvePath(path string) (string, error) {
	fullPath := filepath.Join(l.rootDir, path)
	// Ensure the fullPath is clean (e.g., no /./, /../ resolved)
	cleanedPath := filepath.Clean(fullPath)

	// Check if the cleaned path starts with the root directory.
	// If it doesn't, or if the root directory is not exactly the same and
	// the path after the root directory doesn't start with a separator,
	// it means the path attempts to escape the root directory.
	if !strings.HasPrefix(cleanedPath, l.rootDir) ||
		(len(cleanedPath) > len(l.rootDir) && cleanedPath[len(l.rootDir)] != os.PathSeparator && cleanedPath != l.rootDir) { // Added cleanedPath != l.rootDir check
		return "", runtimeerrors.NewStructured(objectstore.Module, "path %q attempts to escape root directory %q", path, l.rootDir).WithCaller()
	}

	return cleanedPath, nil
}
