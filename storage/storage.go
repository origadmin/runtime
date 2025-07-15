package storage

import (
	"fmt"
	"path/filepath"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	indexiface "github.com/origadmin/runtime/interfaces/storage/index"
	blobimpl "github.com/origadmin/runtime/storage/blob"
	contentimpl "github.com/origadmin/runtime/storage/content"
	indeximpl "github.com/origadmin/runtime/storage/index"
	metaimpl "github.com/origadmin/runtime/storage/meta"
)

const (
	// DefaultChunkSize specifies the default block size for file splitting operations
	// when it's not provided in the configuration. The value is 4MB.
	DefaultChunkSize = 4 * 1024 * 1024
)

// Storage defines the high-level interface for the storage service.
type Storage interface {
	GetIndexManager() indexiface.Manager
	GetMetaStore() *metaimpl.Meta
}

// storage represents the assembled storage service.
// It implements the Storage interface.
type storage struct {
	IndexManager indexiface.Manager
	MetaStore    *metaimpl.Meta
}

// New creates a new Storage service instance based on the provided protobuf configuration.
// This function acts as the entry point for creating the storage system.
func New(cfg *configv1.Storage) (Storage, error) {
	// 3. Add validation logic for the new protobuf config structure
	if cfg == nil {
		return nil, fmt.Errorf("storage config cannot be nil")
	}
	if cfg.GetType() != "filestore" {
		return nil, fmt.Errorf("this New function only supports 'filestore' type, got '%s'", cfg.GetType())
	}

	fsCfg := cfg.GetFilestore()
	if fsCfg == nil {
		return nil, fmt.Errorf("filestore config block is missing")
	}
	if fsCfg.GetDriver() != "local" {
		return nil, fmt.Errorf("this New function only supports 'local' filestore driver, got '%s'", fsCfg.GetDriver())
	}

	localCfg := fsCfg.GetLocal()
	if localCfg == nil {
		return nil, fmt.Errorf("local config block is missing for filestore driver 'local'")
	}

	basePath := localCfg.GetRoot()
	if basePath == "" {
		return nil, fmt.Errorf("storage config: filestore.local.root (BasePath) cannot be empty")
	}

	// NOTE: DefaultChunkSize is not part of the proto config.
	// Using a hardcoded default. Consider adding this to the FileLocal message if needed.
	// Use the chunk size from the proto config, with a fallback to a sensible default.
	defaultChunkSize := fsCfg.GetChunkSize()
	if defaultChunkSize == 0 {
		defaultChunkSize = 4 * 1024 * 1024 // 4MB default
	}

	// 1. Create base paths for each component
	blobBasePath := filepath.Join(basePath, "blobs")
	metaBasePath := filepath.Join(basePath, "meta")
	indexPath := filepath.Join(basePath, "index")

	// 2. Instantiate Blob Store
	blobStore := blobimpl.New(blobBasePath)

	// 3. Instantiate Content Assembler
	contentAssembler := contentimpl.New(blobStore)

	// 4. Instantiate Meta Store
	metaStore, err := metaimpl.NewFileMetaStore(metaBasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create meta store: %w", err)
	}

	// 5. Instantiate Meta Service (uses MetaStore, BlobStore, ContentAssembler)
	metaService, err := metaimpl.New(metaStore, blobStore, contentAssembler, defaultChunkSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create meta service: %w", err)
	}

	// 6. Instantiate Index Manager
	indexManager, err := indeximpl.NewFileManager(indexPath, metaStore)
	if err != nil {
		return nil, fmt.Errorf("failed to create index manager: %w", err)
	}

	return &storage{
		IndexManager: indexManager,
		MetaStore:    metaService,
	}, nil
}

// GetIndexManager returns the index manager component.
func (s *storage) GetIndexManager() indexiface.Manager {
	return s.IndexManager
}

// GetMetaStore returns the meta store component.
func (s *storage) GetMetaStore() *metaimpl.Meta {
	return s.MetaStore
}
