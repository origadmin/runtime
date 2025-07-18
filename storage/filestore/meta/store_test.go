/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package meta

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	metaiface "github.com/origadmin/runtime/interfaces/storage/components/meta"
	metav2 "github.com/origadmin/runtime/storage/filestore/meta/v2"
)

func TestFileMetaStore(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "metastore-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new FileMetaStore
	store, err := NewFileMetaStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create FileMetaStore: %v", err)
	}

	// Test FileMetaV2 with embedded data
	embeddedData := []byte("this is some embedded metaFile content")
	metaV2Embedded := &metaiface.StoreMeta[metav2.FileMetaV2]{
		Version: metav2.Version,
		Info: &metaiface.FileMetaInfo{
			Name: "test_embedded.txt",
		},
		Data: &metav2.FileMetaV2{
			Version:      metav2.Version,
			FileSize:     int64(len(embeddedData)),
			ModifyTime:   time.Now().UnixNano(),
			MimeType:     "text/plain",
			RefCount:     1,
			EmbeddedData: embeddedData,
		},
	}

	// Test Create
	idEmbedded, err := store.Create(metaV2Embedded)
	if err != nil {
		t.Fatalf("Create embedded failed: %v", err)
	}
	if idEmbedded == "" {
		t.Fatal("Create embedded returned empty ID")
	}

	// Test Get embedded
	retrievedMetaEmbedded, err := store.Get(idEmbedded, 0)
	if err != nil {
		t.Fatalf("Get embedded failed: %v", err)
	}
	retrievedMetaV2Embedded, ok := retrievedMetaEmbedded.(*metaiface.StoreMeta[metav2.FileMetaV2])
	if !ok {
		t.Fatalf("Retrieved meta is not StoreMeta[FileMetaV2] type")
	}
	if string(retrievedMetaV2Embedded.Data.EmbeddedData) != string(embeddedData) {
		t.Errorf("Embedded data mismatch. got %s, want %s", retrievedMetaV2Embedded.Data.EmbeddedData, embeddedData)
	}
	if retrievedMetaV2Embedded.Data.Size != int64(len(embeddedData)) {
		t.Errorf("Embedded size mismatch. got %d, want %d", retrievedMetaV2Embedded.Data.Size, len(embeddedData))
	}

	// Test FileMetaV2 with blob hashes (no embedded data)
	blobHashes := []string{"hash1", "hash2", "hash3"}
	metaV2Blob := &metaiface.StoreMeta[metav2.FileMetaV2]{
		Version: metav2.Version,
		Info: &metaiface.FileMetaInfo{
			Name: "test_blob.bin",
		},
		Data: &metav2.FileMetaV2{
			Version:    metav2.Version,
			FileSize:   1024 * 1024, // 1MB
			ModifyTime: time.Now().UnixNano(),
			MimeType:   "application/octet-stream",
			RefCount:   1,
			BlobHashes: blobHashes,
		},
	}

	// Test Create blob
	idBlob, err := store.Create(metaV2Blob)
	if err != nil {
		t.Fatalf("Create blob failed: %v", err)
	}
	if idBlob == "" {
		t.Fatal("Create blob returned empty ID")
	}

	// Test Get blob
	retrievedMetaBlob, err := store.Get(idBlob, 0)
	if err != nil {
		t.Fatalf("Get blob failed: %v", err)
	}
	retrievedMetaV2Blob, ok := retrievedMetaBlob.(*metaiface.StoreMeta[metav2.FileMetaV2])
	if !ok {
		t.Fatalf("Retrieved meta is not StoreMeta[FileMetaV2] type")
	}
	if len(retrievedMetaV2Blob.Data.BlobHashes) != len(blobHashes) {
		t.Errorf("Blob hashes count mismatch. got %d, want %d", len(retrievedMetaV2Blob.Data.BlobHashes), len(blobHashes))
	}
	for i, hash := range blobHashes {
		if retrievedMetaV2Blob.Data.BlobHashes[i] != hash {
			t.Errorf("Blob hash mismatch at index %d. got %s, want %s", i, retrievedMetaV2Blob.Data.BlobHashes[i], hash)
		}
	}

	// Test Update
	metaV2Embedded.Data.Size = 200
	metaV2Embedded.Data.MimeType = "image/png"
	err = store.Update(idEmbedded, metaV2Embedded)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify Update
	updatedMetaEmbedded, err := store.Get(idEmbedded, 0)
	if err != nil {
		t.Fatalf("Get after update failed: %v", err)
	}
	updatedMetaV2Embedded, ok := updatedMetaEmbedded.(*metaiface.StoreMeta[metav2.FileMetaV2])
	if !ok {
		t.Fatalf("Updated meta is not StoreMeta[FileMetaV2] type")
	}
	if updatedMetaV2Embedded.Data.Size != 200 {
		t.Errorf("Updated size mismatch. got %d, want %d", updatedMetaV2Embedded.Data.Size, 200)
	}
	if updatedMetaV2Embedded.Data.MimeType != "image/png" {
		t.Errorf("Updated mime type mismatch. got %s, want %s", updatedMetaV2Embedded.Data.MimeType, "image/png")
	}

	// Test Exists
	exists, err := store.Exists(idEmbedded)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Errorf("Exists returned false for a meta that should exist")
	}

	// Test Delete
	err = store.Delete(idEmbedded)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify Delete
	exists, err = store.Exists(idEmbedded)
	if err != nil {
		t.Fatalf("Exists after delete failed: %v", err)
	}
	if exists {
		t.Errorf("Exists returned true for a meta that should have been deleted")
	}

	// Test Get non-existent
	_, err = store.Get(idEmbedded, 0)
	if err == nil {
		t.Errorf("Get non-existent did not return an error")
	}

	// Test unsupported version
	// This requires creating a dummy meta with an unsupported version
	// For now, we'll skip this as it's hard to mock without changing the actual FileMetaV2 struct.
}
