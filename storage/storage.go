package storage

import (
	"fmt"
	"path/filepath"

	blob_impl "github.com/origadmin/runtime/storage/blob"
	content_impl "github.com/origadmin/runtime/storage/content"
	index_impl "github.com/origadmin/runtime/storage/index"
	layout_impl "github.com/origadmin/runtime/storage/layout"


eta_impl "github.com/origadmin/runtime/storage/meta"

	indexiface "gi
	metaiface "github.com/origadmin/runtime/interfaces/storage/meta"
	indexiface "github.com/origadmin/runtime/interfaces/storage/index"
)

// Config holds the configuration for the storage service.
type Config struct {
	BasePath string
	// Add more configuration options here for different blob/meta/content implementations
	// e.g., BlobStoreType string, MetaStoreType string
}

// Storage represents the assembled storage service.
// In a full implementation, this might be a facade implementing a higher-level Storage interface.
type Storage struct {
	IndexManager indexiface.Manager
}

// New creates a new Storage service instance based on the provided configuration.
// This function acts as the entry point for creating the storage system.
func New(cfg Config) (*Storage, error) {
	if cfg.BasePath == "" {
		return nil, fmt.Errorf("storage config: BasePath cannot be empty")
	}

	// Create base paths for each component
	blobBasePath := filepath.Join(cfg.BasePath, "blobs")
	metaBasePath := filepath.Join(cfg.BasePath, "meta")
	indexPath := filepath.Join(cfg.BasePath, "index")

	// 1. Instantiate Layouts for each component
	blobLayout, err := layout_impl.NewLocalShardedStorage(blobBasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob layout: %w", err)
	}
	metaLayout, err := layout_impl.NewLocalShardedStorage(metaBasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create meta layout: %w", err)
	}
	indexLayout, err := layout_impl.NewLocalShardedStorage(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create index layout: %w", err)
	}

	// 2. Instantiate Blob Store
	blobStore, err := blob_impl.NewFileStore(blobLayout)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob store: %w", err)
	}

	// 3. Instantiate Content Assembler
	contentAssembler := content_impl.New(blobStore)

	// 4. Instantiate Meta Store
	metaStore, err := meta_impl.NewFileMetaStore(metaLayout)
	if err != nil {
		return nil, fmt.Errorf("failed to create meta store: %w", err)
	}

	// 5. Instantiate Meta Service (uses MetaStore, BlobStore, ContentAssembler)
	metaService, err := meta_impl.New(metaStore, blobStore, contentAssembler)
	if err != nil {
		return nil, fmt.Errorf("failed to create meta service: %w", err)
	}

	// 6. Instantiate Index Manager
	indexManager, err := index_impl.NewFileManager(indexPath, metaService, indexLayout)
	if err != nil {
		return nil, fmt.Errorf("failed to create index manager: %w", err)
	}

	return &Storage{
		IndexManager: indexManager,
	}, nil
}
