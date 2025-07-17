/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	blobiface "github.com/origadmin/runtime/interfaces/storage/blob"
	contentiface "github.com/origadmin/runtime/interfaces/storage/content"
	metaiface "github.com/origadmin/runtime/interfaces/storage/meta"
	metav2 "github.com/origadmin/runtime/storage/meta/v2"
)

// chunkData reads from the reader and splits the content into blocks of a fixed size.
// For each block, it calls the provided store function.
func (m *Meta) chunkData(r io.Reader) ([]string, int64, error) {
	var hashes []string
	var totalSize int64
	buf := make([]byte, m.chunkSize)

	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return nil, 0, err
		}
		if n == 0 {
			break
		}

		data := buf[:n]
		hash, storeErr := m.blobStorage.Write(data)
		if storeErr != nil {
			return nil, 0, storeErr
		}
		hashes = append(hashes, hash)
		totalSize += int64(n)

		if err == io.EOF {
			break
		}
	}

	return hashes, totalSize, nil
}

// Meta is a high-level service for managing file content and its metadata.
// It orchestrates interactions between the metadata store (metaStore) and the blob store (blobStorage).
// It is stateless and operates on content IDs (metaID), not paths.
type Meta struct {
	metaStore   metaiface.Store
	blobStorage blobiface.Store
	assembler   contentiface.Assembler
	chunkSize   int64 // Configurable chunk size for writing large files
}

// New creates a new Meta instance.
func New(metaStore metaiface.Store, blobStorage blobiface.Store, assembler contentiface.Assembler, chunkSize int64) (*Meta, error) {
	if chunkSize <= 0 {
		chunkSize = metav2.EmbeddedFileSizeThreshold // Use a sensible default if not provided
	}
	m := &Meta{
		metaStore:   metaStore,
		blobStorage: blobStorage,
		assembler:   assembler,
		chunkSize:   chunkSize,
	}
	return m, nil
}

// Create reads content from a reader, stores it (either embedded or as sharded blobs),
// persists the metadata, and returns the unique metadata ID.
// It optimizes based on whether the input `size` is known.
func (m *Meta) Create(r io.Reader, size int64) (string, error) {
	var fileMeta metaiface.FileMeta
	var isLargeFile bool

	if size > 0 { // Case 1: Size is known
		if size <= metav2.EmbeddedFileSizeThreshold {
			// Known small file: read exact size and embed.
			contentBytes := make([]byte, size)
			if _, err := io.ReadFull(r, contentBytes); err != nil {
				return "", fmt.Errorf("failed to read content for known small file: %w", err)
			}
			fileMeta = &metav2.FileMetaV2{
				FileSize:     size,
				ModifyTime:   time.Now().Unix(),
				MimeType:     "application/octet-stream",
				RefCount:     1,
				EmbeddedData: contentBytes,
			}
		} else {
			// Known large file: chunk directly.
			isLargeFile = true
			hashes, totalSize, chunkErr := m.chunkData(r)
			if chunkErr != nil {
				return "", fmt.Errorf("failed to chunk and store data for known large file: %w", chunkErr)
			}
			// Sanity check
			if totalSize != size {
				// Cleanup blobs if size doesn't match
				for _, h := range hashes {
					_ = m.blobStorage.Delete(h)
				}
				return "", fmt.Errorf("stream size mismatch: provided size %d, but read %d", size, totalSize)
			}
			fileMeta = &metav2.FileMetaV2{
				FileSize:   totalSize,
				ModifyTime: time.Now().Unix(),
				MimeType:   "application/octet-stream",
				RefCount:   1,
				BlobHashes: hashes,
			}
		}
	} else { // Case 2: Size is unknown, fall back to peeking.
		buf := new(bytes.Buffer)
		tee := io.TeeReader(r, buf)
		_, err := io.CopyN(io.Discard, tee, metav2.EmbeddedFileSizeThreshold+1)
		if err != nil && err != io.EOF {
			return "", err
		}

		contentBytes := buf.Bytes()
		if err == io.EOF { // It's a small file
			fileMeta = &metav2.FileMetaV2{
				FileSize:     int64(len(contentBytes)),
				ModifyTime:   time.Now().Unix(),
				MimeType:     "application/octet-stream",
				RefCount:     1,
				EmbeddedData: contentBytes,
			}
		} else { // It's a large file
			isLargeFile = true
			fullStream := io.MultiReader(bytes.NewReader(contentBytes), r)
			hashes, totalSize, chunkErr := m.chunkData(fullStream)
			if chunkErr != nil {
				return "", fmt.Errorf("failed to chunk and store data for unknown size file: %w", chunkErr)
			}
			fileMeta = &metav2.FileMetaV2{
				FileSize:   totalSize,
				ModifyTime: time.Now().Unix(),
				MimeType:   "application/octet-stream",
				RefCount:   1,
				BlobHashes: hashes,
			}
		}
	}

	// Persist the metadata using the underlying metaStore
	id, err := m.metaStore.Create(fileMeta)
	if err != nil {
		// If persisting metadata fails, we must clean up any blobs we just created.
		if isLargeFile {
			if v2, ok := fileMeta.(*metav2.FileMetaV2); ok {
				for _, h := range v2.BlobHashes {
					_ = m.blobStorage.Delete(h) // Best-effort cleanup
				}
			}
		}
		return "", fmt.Errorf("failed to create meta in store: %w", err)
	}

	return id, nil
}

// Get retrieves file metadata by its ID.
func (m *Meta) Get(id string) (metaiface.FileMeta, error) {
	return m.metaStore.Get(id)
}

// Read creates a reader for a file's content by its metadata ID.
func (m *Meta) Read(id string) (io.ReadCloser, error) {
	fileMeta, err := m.metaStore.Get(id)
	if err != nil {
		return nil, err
	}

	reader, err := m.assembler.NewReader(fileMeta)
	if err != nil {
		return nil, err
	}
	// The assembler's reader is not necessarily a ReadCloser, so we wrap it.
	return io.NopCloser(reader), nil
}

// Delete orchestrates the deletion of file content (blobs) and the metadata record.
func (m *Meta) Delete(id string) error {
	// 1. Get metadata to find blob hashes for large files.
	fileMeta, err := m.metaStore.Get(id)
	if err != nil {
		if os.IsNotExist(err) { // Or a more specific error from the layout store
			return nil // Idempotent: if it's already gone, we're done.
		}
		return fmt.Errorf("failed to get metadata for deletion (id: %s): %w", id, err)
	}

	// 2. If it's a large file with sharded blobs, delete them.
	if v2, ok := fileMeta.(*metav2.FileMetaV2); ok {
		if len(v2.BlobHashes) > 0 {
			// TODO: Consider collecting errors and continuing instead of stopping on the first one.
			for _, blobHash := range v2.BlobHashes {
				if err := m.blobStorage.Delete(blobHash); err != nil && !os.IsNotExist(err) {
					return fmt.Errorf("failed to delete blob %s for meta %s: %w", blobHash, id, err)
				}
			}
		}
	}

	// 3. Finally, delete the metadata record itself.
	return m.metaStore.Delete(id)
}
