package filestore

import (
	"fmt"
	"io"
	"path"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/origadmin/toolkits/errors"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	indexiface "github.com/origadmin/runtime/interfaces/storage/components/index"
	"github.com/origadmin/runtime/log"
	blobimpl "github.com/origadmin/runtime/storage/filestore/blob"
	contentimpl "github.com/origadmin/runtime/storage/filestore/content"
	indeximpl "github.com/origadmin/runtime/storage/filestore/index"
	layoutimpl "github.com/origadmin/runtime/storage/filestore/layout"
	
)

const (
	// DefaultChunkSize specifies the default block size for file splitting operations
	// when it's not provided in the configuration. The value is 4MB.
	DefaultChunkSize = 4 * 1024 * 1024
)

// storage represents the assembled Store service.
// It implements the Storage interface.
type storage struct {
	index     indexiface.Manager
	metaStore *metaimpl.Service
}

// List returns a list of files and directories at the given path.
func (s *storage) List(path string) ([]storageiface.FileInfo, error) {
	node, err := s.index.GetNodeByPath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list path '%s': %w", path, err)
	}
	if node.NodeType != indexiface.Directory {
		return nil, fmt.Errorf("path '%s' is not a directory", path)
	}

	nodes, err := s.index.ListChildren(node.NodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to list children of path '%s': %w", path, err)
	}
	infos := make([]storageiface.FileInfo, 0, len(nodes))
	for _, node := range nodes {
		info := storageiface.FileInfo{
			Name:    node.Name,
			Path:    filepath.Join(path, node.Name),
			IsDir:   node.NodeType == indexiface.Directory,
			ModTime: node.Mtime,
		}

		// For files, we fetch more accurate metadata from the metaStore.
		if !info.IsDir && node.MetaHash != "" {
			fileMeta, err := s.metaStore.Get(node.MetaHash)
			if err != nil {
				// Log the error and skip this file, so one bad file doesn't break the whole listing.
				log.Warnf("Failed to get metadata for file '%s' (id: %s): %v. Skipping.", info.Path, node.MetaHash, err)
				continue
			}
			info.Size = fileMeta.Size()
			info.ModTime = fileMeta.ModTime() // Use the more accurate content modification time.
		}
		infos = append(infos, info)
	}

	return infos, nil
}

// Stat returns metadata for a single file or directory.
func (s *storage) Stat(path string) (storageiface.FileInfo, error) {
	node, err := s.index.GetNodeByPath(path)
	if err != nil {
		return storageiface.FileInfo{}, fmt.Errorf("failed to stat path '%s': %w", path, err)
	}

	info := storageiface.FileInfo{
		Name:    node.Name,
		Path:    path,
		IsDir:   node.NodeType == indexiface.Directory,
		ModTime: node.Mtime,
	}

	if !info.IsDir && node.MetaHash != "" {
		fileMeta, err := s.metaStore.Get(node.MetaHash)
		if err != nil {
			return storageiface.FileInfo{}, fmt.Errorf("failed to get metadata for file '%s': %w", path, err)
		}
		info.Size = fileMeta.Size()
		info.ModTime = fileMeta.ModTime()
	}

	return info, nil
}

// Read opens a file for reading.
func (s *storage) Read(path string) (io.ReadCloser, error) {
	// 1. Find the index node to get the content's metadata ID.
	node, err := s.index.GetNodeByPath(path)
	if err != nil {
		return nil, fmt.Errorf("file not found at path '%s': %w", path, err)
	}

	if node.NodeType != indexiface.File {
		return nil, fmt.Errorf("path '%s' is a directory, not a file", path)
	}

	// 2. The metaStore service provides a readable stream, abstracting away the blob assembly.
	return s.metaStore.Read(node.MetaHash)
}

// Mkdir creates a new directory.
func (s *storage) Mkdir(filepath string) error {
	// Creating a directory is a pure index operation; it has no content.
	dir, name := path.Split(filepath)
	node, err := s.index.GetNodeByPath(dir)
	if err != nil {
		return err
	}
	if err := s.index.CreateNode(&indexiface.Node{
		NodeID:   uuid.Must(uuid.NewRandom()).String(),
		ParentID: node.NodeID,
		Name:     name,
		NodeType: indexiface.Directory,
		Mtime:    time.Now(),
	}); err != nil {
		return err
	}
	return nil
}

// Delete removes a file or an empty directory.
func (s *storage) Delete(path string) error {
	node, err := s.index.GetNodeByPath(path)
	if err != nil {
		return err
	}
	if err := s.index.DeleteNode(node.NodeID); err != nil {
		// This makes the Delete operation idempotent. Deleting a non-existent file is not an error.
		//if errors.Is(err, indexiface.ErrNodeNotFound) { // Assuming index returns a standard error.
		//	return nil
		//}
		return fmt.Errorf("error looking up '%s' for deletion: %w", path, err)
	}

	// For directories, ensure it's empty before deleting to prevent accidental data loss.
	if node.NodeType == indexiface.Directory {
		node, err := s.index.GetNodeByPath(path)
		if err != nil {
			return err
		}
		children, err := s.index.ListChildren(node.NodeID)
		if err != nil {
			return fmt.Errorf("error checking if directory '%s' is empty: %w", path, err)
		}
		if len(children) > 0 {
			return fmt.Errorf("directory '%s' is not empty", path)
		}
	}

	// First, delete the index entry. This makes the file/dir "disappear" from the user's perspective.
	if err := s.index.DeleteNode(node.NodeID); err != nil {
		return fmt.Errorf("failed to delete index entry for '%s': %w", path, err)
	}

	// If it was a file, proceed to delete its content from the meta/blob store.
	if node.NodeType == indexiface.File && node.MetaHash != "" {
		if err := s.metaStore.Delete(node.MetaHash); err != nil {
			// This is a problematic state (orphaned content). Log and return an error.
			log.Errorf("Index entry for '%s' deleted, but failed to delete content (metaID: %s): %v", path, node.MetaHash, err)
			return fmt.Errorf("index entry for '%s' deleted, but failed to delete content (metaID: %s): %w", path, node.MetaHash, err)
		}
	}

	return nil
}

// Rename moves or renames a file or directory.
func (s *storage) Rename(oldPath, newPath string) error {
	oldNode, err := s.index.GetNodeByPath(oldPath)
	if err != nil {
		return err
	}

	newPath = path.Dir(newPath)
	newNode, err := s.index.GetNodeByPath(newPath)
	if err != nil && errors.Is(err, indexiface.ErrNodeNotFound) {
		return err
	}

	// Renaming is a pure index operation. The file content itself (and its metaID) is not touched.
	return s.index.MoveNode(oldNode.NodeID, newNode.NodeID, newNode.Name)
}

// Write creates a new file at the given path with the content from the reader.
func (s *storage) Write(filepath string, data io.Reader, size int64) error {
	dir, name := path.Split(filepath)
	// 1. Write the content stream to the metaStore. It handles chunking and blob storage,
	// returning metadata (including a unique ID) for the stored content.
	metaID, err := s.metaStore.Create(data, size)
	if err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}
	node, err := s.index.GetNodeByPath(dir)
	if err != nil {
		return err
	}
	// 2. Create an entry in the index to link the path to the newly stored content.
	if err := s.index.CreateNode(&indexiface.Node{
		NodeID:   uuid.Must(uuid.NewRandom()).String(),
		ParentID: node.NodeID,
		Name:     name,
		NodeType: indexiface.File,
		Mtime:    time.Now(),
		MetaHash: metaID,
	}); err != nil {
		// IMPORTANT: If creating the index fails, we must clean up the orphaned content.
		log.Warnf("Failed to create index node for path '%s', attempting to clean up orphaned content (metaID: %s)", filepath, metaID)
		if cleanupErr := s.metaStore.Delete(metaID); cleanupErr != nil {
			// This is a critical failure state. The content is now orphaned.
			log.Errorf("Orphaned content cleanup FAILED for metaID '%s': %v", metaID, cleanupErr)
		}
		return fmt.Errorf("failed to create index node for path '%s': %w", filepath, err)
	}

	return nil
}

// NewStorage creates a new Storage service instance based on the provided protobuf configuration.
// This function acts as the entry point for creating the storage system.
func New(cfg *configv1.FileStore) (storageiface.FileStore, error) {

	if cfg.GetDriver() != "local" {
		return nil, fmt.Errorf("this New function only supports 'local' filestore driver, got '%s'", cfg.GetDriver())
	}

	localCfg := cfg.GetLocal()
	if localCfg == nil {
		return nil, fmt.Errorf("local config block is missing for filestore driver 'local'")
	}

	basePath := localCfg.GetRoot()
	if basePath == "" {
		return nil, fmt.Errorf("storage config: filestore.local.root (BasePath) cannot be empty")
	}

	// Use the chunk size from the proto config, with a fallback to the package-level default.
	chunkSize := cfg.GetChunkSize()
	if chunkSize == 0 {
		chunkSize = DefaultChunkSize
	}

	// 1. Create base paths for each component
	blobBasePath := filepath.Join(basePath, "blobs")
	metaBasePath := filepath.Join(basePath, "meta")
	indexPath := filepath.Join(basePath, "index")

	// 2. Instantiate Layout for Blob Store
	blobLayout, err := layoutimpl.NewLocalShardedStorage(blobBasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob layout: %w", err)
	}

	// 3. Instantiate Blob Store
	blobStore := blobimpl.New(blobLayout)

	// 4. Instantiate Content Assembler
	contentAssembler := contentimpl.New(blobStore)

	// 5. Instantiate low-level Meta Store
	lowLevelMetaStore, err := metaimpl.NewStore(metaBasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create meta store: %w", err)
	}

	// 6. Instantiate high-level Meta Service (uses MetaStore, BlobStore, ContentAssembler)
	metaService, err := metaimpl.NewService(lowLevelMetaStore, blobLayout, contentAssembler, chunkSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create meta service: %w", err)
	}

	// 7. Instantiate Index Manager
	indexManager, err := indeximpl.NewManager(indexPath, lowLevelMetaStore)
	if err != nil {
		return nil, fmt.Errorf("failed to create index manager: %w", err)
	}

	return &storage{
		index:     indexManager,
		metaStore: metaService,
	}, nil
}
