/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"

	blobiface "github.com/origadmin/runtime/interfaces/storage/components/blob"
	contentiface "github.com/origadmin/runtime/interfaces/storage/components/content"
	metaiface "github.com/origadmin/runtime/interfaces/storage/components/meta"
	blobimpl "github.com/origadmin/runtime/storage/filestore/blob"
	metav2 "github.com/origadmin/runtime/storage/filestore/meta/v2"
)

// Service is a high-level service for managing file content and its metadata.
// It orchestrates interactions between the metadata store (metaStore) and the blob store (blobStore).
// It is stateless and operates on content IDs (metaID), not paths.
type Service struct {
	metaStore metaiface.Store
	blobStore blobiface.Store
	assembler contentiface.Assembler
	chunkSize int64 // Configurable chunk size for writing large files
}

// calculateContentHash is a helper to generate a hash from byte data.
func calculateContentHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// chunkData reads from the reader and splits the content into blocks of a fixed size.
// It processes a large file stream by chunking it, storing the chunks as blobs,
// and returns a fully populated FileMetaV2 object along with a hash of the entire content stream.
func (s *Service) chunkData(r io.Reader) (string, *metav2.FileMetaV2, error) {
	var hashes []string
	var totalSize int64
	buf := make([]byte, s.chunkSize)

	// Create a hasher that will calculate the hash of the entire stream.
	streamHasher := sha256.New()
	// TeeReader copies data from r to streamHasher as it's being read.
	teeReader := io.TeeReader(r, streamHasher)

	for {
		// Read from the teeReader to ensure the hasher gets the data.
		n, err := teeReader.Read(buf)
		if err != nil && err != io.EOF {
			return "", nil, err
		}
		if n == 0 {
			break
		}

		data := buf[:n]
		// Write the chunk to the blob storage. The blob store will return the hash of this chunk.
		blobHash, storeErr := s.blobStore.Write(data)
		if storeErr != nil {
			return "", nil, storeErr
		}
		hashes = append(hashes, blobHash)
		totalSize += int64(n)

		if err == io.EOF {
			break
		}
	}

	// Finalize the hash of the entire stream.
	overallHash := hex.EncodeToString(streamHasher.Sum(nil))

	meta := &metav2.FileMetaV2{
		FileSize:   totalSize,
		ModifyTime: time.Now().Unix(),
		MimeType:   "application/octet-stream",
		RefCount:   1,
		BlobHashes: hashes,
		BlobSize:   int32(s.chunkSize),
	}

	return overallHash, meta, nil
}

// NewService creates a new Service instance.
func NewService(metaStore metaiface.Store, basePath string, assembler contentiface.Assembler, chunkSize int64) (*Service, error) {
	if chunkSize <= 0 {
		chunkSize = metav2.EmbeddedFileSizeThreshold // Use a sensible default if not provided
	}

	// Instantiate layout internally
	bstore, err := blobimpl.New(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create local sharded storage for meta service: %w", err)
	}

	s := &Service{
		metaStore: metaStore,
		blobStore: bstore, // Pass the layout to blob.New
		assembler: assembler,
		chunkSize: chunkSize,
	}
	return s, nil
}

// Create reads content from a reader, stores it (either embedded or as sharded blobs),
// persists the metadata, and returns the unique metadata ID.
// It optimizes based on whether the input `size` is known.
func (s *Service) Create(r io.Reader, size int64) (string, error) {
	var fileMeta metaiface.FileMeta
	var isLargeFile bool
	var id string

	if size > 0 { // Case 1: Size is known
		if size <= metav2.EmbeddedFileSizeThreshold {
			// Known small file: read exact size and embed.
			contentBytes := make([]byte, size)
			if _, err := io.ReadFull(r, contentBytes); err != nil {
				return "", fmt.Errorf("failed to read content for known small file: %w", err)
			}
			id = calculateContentHash(contentBytes)
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
			chunkID, meta, chunkErr := s.chunkData(r)
			if chunkErr != nil {
				// Note: chunkData does not return partial blob hashes, so no cleanup needed here.
				return "", fmt.Errorf("failed to chunk and store data for known large file: %w", chunkErr)
			}
			// Sanity check to ensure the stream provided the expected amount of data.
			if meta.FileSize != size {
				// Cleanup blobs if size doesn't match
				for _, h := range meta.BlobHashes {
					_ = s.blobStore.Delete(h)
				}
				return "", fmt.Errorf("stream size mismatch: provided size %d, but read %d", size, meta.FileSize)
			}
			id = chunkID
			fileMeta = meta
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
			id = calculateContentHash(contentBytes)
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
			chunkID, meta, chunkErr := s.chunkData(fullStream)
			if chunkErr != nil {
				return "", fmt.Errorf("failed to chunk and store data for unknown size file: %w", chunkErr)
			}
			id = chunkID
			fileMeta = meta
		}
	}

	// This case handles a known-size-0 or unknown-size empty stream.
	if id == "" && (fileMeta == nil || fileMeta.Size() == 0) {
		id = calculateContentHash([]byte{})
		if fileMeta == nil {
			fileMeta = &metav2.FileMetaV2{
				FileSize:   0,
				ModifyTime: time.Now().Unix(),
				MimeType:   "application/octet-stream",
				RefCount:   1,
			}
		}
	}

	// Persist the metadata using the underlying metaStore
	if err := s.metaStore.Create(id, fileMeta); err != nil {
		// If persisting metadata fails, we must clean up any blobs we just created.
		if isLargeFile {
			if v2, ok := fileMeta.(*metav2.FileMetaV2); ok {
				for _, h := range v2.BlobHashes {
					_ = s.blobStore.Delete(h) // Best-effort cleanup
				}
			}
		}
		return "", fmt.Errorf("failed to create meta in store: %w", err)
	}

	return id, nil
}

// Get retrieves file metadata by its ID.
func (s *Service) Get(id string) (metaiface.FileMeta, error) {
	return s.metaStore.Get(id)
}

// Read creates a reader for a file's content by its metadata ID.
func (s *Service) Read(id string) (io.ReadCloser, error) {
	fileMeta, err := s.metaStore.Get(id)
	if err != nil {
		return nil, err
	}

	reader, err := s.assembler.NewReader(fileMeta)
	if err != nil {
		return nil, err
	}
	// The assembler's reader is not necessarily a ReadCloser, so we wrap it.
	return io.NopCloser(reader), nil
}

// Delete orchestrates the deletion of file content (blobs) and the metadata record.
func (s *Service) Delete(id string) error {
	// 1. Get metadata to find blob hashes for large files.
	fileMeta, err := s.metaStore.Get(id)
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
				if err := s.blobStore.Delete(blobHash); err != nil && !os.IsNotExist(err) {
					return fmt.Errorf("failed to delete blob %s for meta %s: %w", blobHash, id, err)
				}
			}
		}
	}

	// 3. Finally, delete the metadata record itself.
	return s.metaStore.Delete(id)
}
