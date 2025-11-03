package meta

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	blobiface "github.com/origadmin/runtime/interfaces/storage/components/blob"
	contentiface "github.com/origadmin/runtime/interfaces/storage/components/content"
	metaiface "github.com/origadmin/runtime/interfaces/storage/components/meta"
	metav2 "github.com/origadmin/runtime/storage/filestore/meta/v2"
)

// Mock implementations for interfaces

type mockMetaStore struct {
	sync.RWMutex
	data map[string]metaiface.FileMeta
}

func newMockMetaStore() *mockMetaStore {
	return &mockMetaStore{
		data: make(map[string]metaiface.FileMeta),
	}
}

func (m *mockMetaStore) Create(id string, fileMeta metaiface.FileMeta) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.data[id]; ok {
		return fmt.Errorf("meta with id %s already exists", id)
	}
	m.data[id] = fileMeta
	return nil
}

func (m *mockMetaStore) Get(id string) (metaiface.FileMeta, error) {
	m.RLock()
	defer m.RUnlock()
	if meta, ok := m.data[id]; ok {
		return meta, nil
	}
	return nil, os.ErrNotExist
}

func (m *mockMetaStore) Exists(id string) (bool, error) {
	m.RLock()
	defer m.RUnlock()
	_, ok := m.data[id]
	return ok, nil
}

func (m *mockMetaStore) Update(id string, fileMeta metaiface.FileMeta) (metaiface.FileMeta, error) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.data[id]; !ok {
		return nil, os.ErrNotExist
	}
	m.data[id] = fileMeta
	return fileMeta, nil
}

func (m *mockMetaStore) Delete(id string) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.data[id]; !ok {
		return os.ErrNotExist
	}
	delete(m.data, id)
	return nil
}

func (m *mockMetaStore) Migrate(id string) (metaiface.FileMeta, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockMetaStore) CurrentVersion() int {
	return metav2.Version
}

type mockBlobStore struct {
	sync.RWMutex
	data map[string][]byte
}

func (m *mockBlobStore) Read(hash string) ([]byte, error) {
	return m.data[hash], nil
}

func (m *mockBlobStore) Exists(hash string) (bool, error) {
	_, ok := m.data[hash]
	return ok, nil
}

func newMockBlobStore() *mockBlobStore {
	return &mockBlobStore{
		data: make(map[string][]byte),
	}
}

func (m *mockBlobStore) Write(content []byte) (string, error) {
	m.Lock()
	defer m.Unlock()
	h := sha256.New()
	h.Write(content)
	hash := hex.EncodeToString(h.Sum(nil))
	m.data[hash] = content
	return hash, nil
}

func (m *mockBlobStore) Delete(hash string) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.data[hash]; !ok {
		return os.ErrNotExist
	}
	delete(m.data, hash)
	return nil
}

func (m *mockBlobStore) Retrieve(hash string) ([]byte, error) {
	m.RLock()
	defer m.RUnlock()
	if content, ok := m.data[hash]; ok {
		return content, nil
	}
	return nil, os.ErrNotExist
}

type mockAssembler struct {
	blobStore blobiface.Store
	chunkSize int64
}

func newMockAssembler(blobStore blobiface.Store, chunkSize int64) *mockAssembler {
	return &mockAssembler{
		blobStore: blobStore,
		chunkSize: chunkSize,
	}
}

func (ma *mockAssembler) NewReader(fileMeta metaiface.FileMeta) (io.Reader, error) {
	if fileMetaV2, ok := fileMeta.(*metav2.FileMetaV2); ok {
		if fileMetaV2.EmbeddedData != nil {
			return bytes.NewReader(fileMetaV2.EmbeddedData), nil
		} else if len(fileMetaV2.BlobHashes) > 0 {
			readers := make([]io.Reader, len(fileMetaV2.BlobHashes))
			for i, hash := range fileMetaV2.BlobHashes {
				data, err := ma.blobStore.Read(hash)
				if err != nil {
					return nil, err
				}
				readers[i] = bytes.NewReader(data)
			}
			return io.MultiReader(readers...), nil
		}
	}
	return nil, fmt.Errorf("unsupported file meta type or empty file")
}

func (ma *mockAssembler) WriteContent(r io.Reader, size int64) (contentID string, fileMeta metaiface.FileMeta, err error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	// Calculate content ID
	h := sha256.New()
	if _, err := io.Copy(h, tee); err != nil {
		return "", nil, fmt.Errorf("failed to hash content: %w", err)
	}
	contentID = hex.EncodeToString(h.Sum(nil))

	fullContent := buf.Bytes()
	actualSize := int64(len(fullContent))

	if size > 0 && size != actualSize {
		return "", nil, fmt.Errorf("size mismatch: expected %d, got %d", size, actualSize)
	}

	if actualSize <= metav2.EmbeddedFileSizeThreshold {
		// Small file, embed data
		fileMeta = &metav2.FileMetaV2{
			FileSize:     actualSize,
			ModifyTime:   time.Now().Unix(),
			MimeType:     "application/octet-stream",
			RefCount:     1,
			EmbeddedData: fullContent,
		}
	} else {
		// Large file, chunk data
		var blobHashes []string
		reader := bytes.NewReader(fullContent)
		for {
			chunk := make([]byte, ma.chunkSize)
			n, err := reader.Read(chunk)
			if err != nil && err != io.EOF {
				return "", nil, fmt.Errorf("failed to read chunk: %w", err)
			}
			if n == 0 {
				break
			}
			hash, storeErr := ma.blobStore.Write(chunk[:n])
			if storeErr != nil {
				return "", nil, fmt.Errorf("failed to store blob chunk: %w", storeErr)
			}
			blobHashes = append(blobHashes, hash)
			if err == io.EOF {
				break
			}
		}

		fileMeta = &metav2.FileMetaV2{
			FileSize:   actualSize,
			ModifyTime: time.Now().Unix(),
			MimeType:   "application/octet-stream",
			RefCount:   1,
			BlobHashes: blobHashes,
			BlobSize:   int32(ma.chunkSize),
		}
	}

	return contentID, fileMeta, nil
}

func (ma *mockAssembler) NewWriter(r io.Reader) (contentiface.Writer, error) {
	return nil, fmt.Errorf("not implemented")
}

// Test functions for the new Service API

func TestService_CreateAndRead(t *testing.T) {
	metaStore := newMockMetaStore()
	blobStore := newMockBlobStore()
	assembler := newMockAssembler(blobStore, 1024) // Use a small chunk size for testing

	service, err := NewService(metaStore, "dummyBasePath", assembler, 1024, func(opts *ServiceOptions) {
		opts.BlobStore = blobStore
	})
	if err != nil {
		t.Fatalf("NewService failed: %v", err)
	}

	// Test small file
	smallContent := "Hello, world!"
	smallReader := bytes.NewReader([]byte(smallContent))
	smallID, err := service.Create(smallReader, int64(len(smallContent)))
	if err != nil {
		t.Fatalf("Create small file failed: %v", err)
	}

	readCloser, err := service.Read(smallID)
	if err != nil {
		t.Fatalf("Read small file failed: %v", err)
	}
	defer readCloser.Close()

	readContent, err := io.ReadAll(readCloser)
	if err != nil {
		t.Fatalf("ReadAll small file failed: %v", err)
	}
	if string(readContent) != smallContent {
		t.Fatalf("Small file content mismatch: expected %q, got %q", smallContent, string(readContent))
	}

	// Test large file (larger than EmbeddedFileSizeThreshold)
	largeContent := make([]byte, metav2.EmbeddedFileSizeThreshold+100)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}
	largeReader := bytes.NewReader(largeContent)
	largeID, err := service.Create(largeReader, int64(len(largeContent)))
	if err != nil {
		t.Fatalf("Create large file failed: %v", err)
	}

	readCloser, err = service.Read(largeID)
	if err != nil {
		t.Fatalf("Read large file failed: %v", err)
	}
	defer readCloser.Close()

	readContent, err = io.ReadAll(readCloser)
	if err != nil {
		t.Fatalf("ReadAll large file failed: %v", err)
	}
	if !bytes.Equal(readContent, largeContent) {
		t.Fatalf("Large file content mismatch")
	}

	// Test non-existent file
	_, err = service.Read("non-existent-id")
	if !os.IsNotExist(err) {
		t.Fatalf("Expected ErrNotExist for non-existent file, got %v", err)
	}
}

func TestService_Get(t *testing.T) {
	metaStore := newMockMetaStore()
	blobStore := newMockBlobStore()
	assembler := newMockAssembler(blobStore, 1024)

	service, err := NewService(metaStore, "dummyBasePath", assembler, 1024, func(opts *ServiceOptions) {
		opts.BlobStore = blobStore
	})
	if err != nil {
		t.Fatalf("NewService failed: %v", err)
	}

	content := "Metadata test content"
	reader := bytes.NewReader([]byte(content))
	id, err := service.Create(reader, int64(len(content)))
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	fileMeta, err := service.Get(id)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if fileMeta.Size() != int64(len(content)) {
		t.Fatalf("Expected file size %d, got %d", len(content), fileMeta.Size())
	}

	// Test non-existent meta
	_, err = service.Get("non-existent-id")
	if !os.IsNotExist(err) {
		t.Fatalf("Expected ErrNotExist for non-existent meta, got %v", err)
	}
}

func TestService_Delete(t *testing.T) {
	metaStore := newMockMetaStore()
	blobStore := newMockBlobStore()
	assembler := newMockAssembler(blobStore, 1024)

	service, err := NewService(metaStore, "dummyBasePath", assembler, 1024, func(opts *ServiceOptions) {
		opts.BlobStore = blobStore
	})
	if err != nil {
		t.Fatalf("NewService failed: %v", err)
	}

	// Create a small file
	smallContent := "Small file for deletion"
	smallReader := bytes.NewReader([]byte(smallContent))
	smallID, err := service.Create(smallReader, int64(len(smallContent)))
	if err != nil {
		t.Fatalf("Create small file failed: %v", err)
	}

	// Verify it exists
	_, err = service.Get(smallID)
	if err != nil {
		t.Fatalf("Small file not found before deletion: %v", err)
	}

	// Delete the small file
	err = service.Delete(smallID)
	if err != nil {
		t.Fatalf("Delete small file failed: %v", err)
	}

	// Verify it's gone
	_, err = service.Get(smallID)
	if !os.IsNotExist(err) {
		t.Fatalf("Expected ErrNotExist after deleting small file, got %v", err)
	}

	// Create a large file
	largeContent := make([]byte, metav2.EmbeddedFileSizeThreshold+100)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}
	largeReader := bytes.NewReader(largeContent)
	largeID, err := service.Create(largeReader, int64(len(largeContent)))
	if err != nil {
		t.Fatalf("Create large file failed: %v", err)
	}

	// Verify its blobs exist
	fileMeta, err := metaStore.Get(largeID)
	if err != nil {
		t.Fatalf("Failed to get large file meta: %v", err)
	}
	if fileMetaV2, ok := fileMeta.(*metav2.FileMetaV2); ok {
		for _, hash := range fileMetaV2.BlobHashes {
			_, err := blobStore.Retrieve(hash)
			if err != nil {
				t.Fatalf("Blob %s not found before deletion: %v", hash, err)
			}
		}
	} else {
		t.Fatalf("Unexpected file meta type for large file")
	}

	// Delete the large file
	err = service.Delete(largeID)
	if err != nil {
		t.Fatalf("Delete large file failed: %v", err)
	}

	// Verify meta is gone
	_, err = service.Get(largeID)
	if !os.IsNotExist(err) {
		t.Fatalf("Expected ErrNotExist after deleting large file meta, got %v", err)
	}

	// Verify blobs are gone
	if fileMetaV2, ok := fileMeta.(*metav2.FileMetaV2); ok {
		for _, hash := range fileMetaV2.BlobHashes {
			_, err := blobStore.Retrieve(hash)
			if !os.IsNotExist(err) {
				t.Fatalf("Expected ErrNotExist for blob %s after deletion, got %v", hash, err)
			}
		}
	}
}

func TestService_Create_UnknownSize(t *testing.T) {
	// Setup test environment
	tempDir, err := os.MkdirTemp("", "filestore-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	metaStore := newMockMetaStore()
	blobStore := newMockBlobStore()
	assembler := newMockAssembler(blobStore, 1024)

	service, err := NewService(metaStore, "dummyBasePath", assembler, 1024, func(opts *ServiceOptions) {
		opts.BlobStore = blobStore
	})
	if err != nil {
		t.Fatalf("NewService failed: %v", err)
	}

	// Test small file with unknown size
	smallContent := "Unknown size small file"
	smallReader := bytes.NewReader([]byte(smallContent))
	smallID, err := service.Create(smallReader, 0) // Pass 0 for unknown size
	if err != nil {
		t.Fatalf("Create unknown size small file failed: %v", err)
	}

	readCloser, err := service.Read(smallID)
	if err != nil {
		t.Fatalf("Read unknown size small file failed: %v", err)
	}
	defer readCloser.Close()

	readContent, err := io.ReadAll(readCloser)
	if err != nil {
		t.Fatalf("ReadAll unknown size small file failed: %v", err)
	}
	if string(readContent) != smallContent {
		t.Fatalf("Unknown size small file content mismatch: expected %q, got %q", smallContent, string(readContent))
	}

	// Test large file with unknown size
	largeSize := int64(metav2.EmbeddedFileSizeThreshold + 500)
	largeContent := make([]byte, largeSize)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}
	largeReader := bytes.NewReader(largeContent)
	largeID, err := service.Create(largeReader, 0) // Pass 0 for unknown size
	if err != nil {
		t.Fatalf("Create unknown size large file failed: %v", err)
	}

	// Verify the metadata
	meta, err := service.Get(largeID)
	if err != nil {
		t.Fatalf("Failed to get metadata for large file: %v", err)
	}
	if meta.Size() != largeSize {
		t.Fatalf("Expected file size %d, got %d", largeSize, meta.Size())
	}

	// Read back the content
	readCloser, err = service.Read(largeID)
	if err != nil {
		t.Fatalf("Read unknown size large file failed: %v", err)
	}
	defer readCloser.Close()

	readContent, err = io.ReadAll(readCloser)
	if err != nil {
		t.Fatalf("ReadAll unknown size large file failed: %v", err)
	}

	// Compare the content
	if len(readContent) != len(largeContent) {
		t.Fatalf("Content length mismatch: expected %d, got %d", len(largeContent), len(readContent))
	}

	// Compare content in chunks to identify where it might differ
	const chunkSize = 1024
	for i := 0; i < len(largeContent); i += chunkSize {
		end := i + chunkSize
		if end > len(largeContent) {
			end = len(largeContent)
		}
		expectedChunk := largeContent[i:end]
		actualChunk := readContent[i:end]
		if !bytes.Equal(expectedChunk, actualChunk) {
			t.Fatalf("Content mismatch at offset %d-%d", i, end)
		}
	}

	// If we get here, the content should match
	if !bytes.Equal(readContent, largeContent) {
		t.Fatalf("Unknown size large file content mismatch (but chunk comparison passed?)")
	}
}
