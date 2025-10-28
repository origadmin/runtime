package meta

import (
	"os"
	"testing"
	"time"

	metav2 "github.com/origadmin/runtime/storage/filestore/meta/v2"
)

func TestFileMetaStore(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "metastore-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	// Create a new FileMetaStore
	store, err := NewStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create FileMetaStore: %v", err)
	}

	// Test data
	meta := &metav2.FileMetaV2{
		FileSize:   1234,
		ModifyTime: time.Now().UnixNano(),
		MimeType:   "text/plain",
		RefCount:   1,
	}

	// Test Create
	id := "test-id"
	if err := store.Create(id, meta); err != nil {
		t.Fatalf("Failed to create meta: %v", err)
	}

	// Test Get
	retrievedMeta, err := store.Get(id)
	if err != nil {
		t.Fatalf("Failed to get meta: %v", err)
	}

	// Verify type
	switch v := retrievedMeta.(type) {
	case *metav2.FileMetaV2, metav2.FileMetaV2:
		// Type assertion successful
	default:
		t.Fatalf("Expected *metav2.FileMetaV2 or metav2.FileMetaV2, got %T", v)
	}

	// Test Update
	meta.FileSize = 5678
	updatedMeta, err := store.Update(id, meta)
	if err != nil {
		t.Fatalf("Failed to update meta: %v", err)
	}

	// Verify updated meta type and value
	updatedV2Meta, ok := updatedMeta.(*metav2.FileMetaV2)
	if !ok {
		t.Fatalf("Updated meta is not *metav2.FileMetaV2 type, got %T", updatedMeta)
	}

	if updatedV2Meta.FileSize != 5678 {
		t.Errorf("Expected FileSize 5678, got %d", updatedV2Meta.FileSize)
	}
}
