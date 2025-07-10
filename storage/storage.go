package storage

import (
	"fmt"
	"path/filepath"

	indexiface "github.com/origadmin/runtime/interfaces/storage/index"
	blobimpl "github.com/origadmin/runtime/storage/blob"
	contentimpl "github.com/origadmin/runtime/storage/content"
	indeximpl "github.com/origadmin/runtime/storage/index"
	layoutimpl "github.com/origadmin/runtime/storage/layout"
	metaimpl "github.com/origadmin/runtime/storage/meta"
)

// Config holds the configuration for the storage service.
type Config struct {
	BasePath         string
	DefaultChunkSize int64 // New: Default chunk size for file operations
}

// Storage represents the assembled storage service.
// In a full implementation, this might be a facade implementing a higher-level Storage interface.
type Storage struct {
	IndexManager indexiface.Manager
	MetaStore    *metaimpl.Meta // New: Expose the Meta service
}

// New creates a new Storage service instance based on the provided configuration.
// This function acts as the entry point for creating the storage system.
func New(cfg Config) (*Storage, error) {
	if cfg.BasePath == "" {
		return nil, fmt.Errorf("storage config: BasePath cannot be empty")
	}

	// Set default chunk size if not provided
	if cfg.DefaultChunkSize == 0 {
		cfg.DefaultChunkSize = 4 * 1024 * 1024 // 4MB default
	}

	// 1. Create base paths for each component
	blobBasePath := filepath.Join(cfg.BasePath, "blobs")
	metaBasePath := filepath.Join(cfg.BasePath, "meta")
	indexPath := filepath.Join(cfg.BasePath, "index")

	// 2. Instantiate Blob Store
	blobStore, err := blobimpl.New(blobBasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob store: %w", err)
	}

	// 3. Instantiate Content Assembler
	contentAssembler := contentimpl.New(blobStore)

	// 4. Instantiate Meta Store
	metaStore, err := metaimpl.NewFileMetaStore(metaBasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create meta store: %w", err)
	}

	// 5. Instantiate Meta Service (uses MetaStore, BlobStore, ContentAssembler)
	metaService, err := metaimpl.New(metaStore, blobStore, contentAssembler, cfg.DefaultChunkSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create meta service: %w", err)
	}

	// 6. Instantiate Index Manager
	indexManager, err := indeximpl.NewFileManager(indexPath, metaStore)
	if err != nil {
		return nil, fmt.Errorf("failed to create index manager: %w", err)
	}

	return &Storage{
		IndexManager: indexManager,
		MetaStore:    metaService,
	}, nil
}
