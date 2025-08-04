package meta

import (
	"crypto/sha256"
	"encoding/hex"
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
	defer os.RemoveAll(tempDir)

	// Create a new FileMetaStore
	store, err := NewStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create FileMetaStore: %v", err)
	}

	// Test FileMetaV2 with embedded data
	embeddedData := []byte("this is some embedded metaFile content")
	hashBytesEmbedded := sha256.Sum256(embeddedData)
	hashStringEmbedded := hex.EncodeToString(hashBytesEmbedded[:])

	metaV2Embedded := &metav2.FileMetaV2{
		FileSize:     int64(len(embeddedData)),
		ModifyTime:   time.Now().UnixNano(),
		MimeType:     "text/plain",
		RefCount:     1,
		EmbeddedData: embeddedData,
	}

	// Test Create
	err = store.Create(hashStringEmbedded, metaV2Embedded) // Corrected Create call
	if err != nil {
		t.Fatalf("Create embedded failed: %v", err)
	}

	// Test Get embedded
	retrievedMetaEmbedded, err := store.Get(hashStringEmbedded) // Corrected Get call
	if err != nil {
		t.Fatalf("Get embedded failed: %v", err)
	}
	retrievedMetaV2Embedded, ok := retrievedMetaEmbedded.(*metav2.FileMetaV2) // Changed type assertion
	if !ok {
		t.Fatalf("Retrieved meta is not *metav2.FileMetaV2 type")
	}
	if string(retrievedMetaV2Embedded.GetEmbeddedData()) != string(embeddedData) { // Used GetEmbeddedData()
		t.Errorf("Embedded data mismatch. got %s, want %s", retrievedMetaV2Embedded.GetEmbeddedData(), embeddedData)
	}
	if retrievedMetaV2Embedded.Size() != int64(len(embeddedData)) { // Used Size()
		t.Errorf("Embedded size mismatch. got %d, want %d", retrievedMetaV2Embedded.Size(), len(embeddedData))
	}

	// Test FileMetaV2 with blob hashes (no embedded data)
	blobHashes := []string{"hash1", "hash2", "hash3"}
	hashBytesBlob := sha256.Sum256([]byte("blob_content_placeholder")) // Placeholder for blob content hash
	hashStringBlob := hex.EncodeToString(hashBytesBlob[:])

	metaV2Blob := &metav2.FileMetaV2{
		FileSize:   1024 * 1024, // 1MB
		ModifyTime: time.Now().UnixNano(),
		MimeType:   "application/octet-stream",
		RefCount:   1,
		BlobHashes: blobHashes,
	}

	// Test Create blob
	err = store.Create(hashStringBlob, metaV2Blob) // Corrected Create call
	if err != nil {
		t.Fatalf("Create blob failed: %v", err)
	}

	// Test Get blob
	retrievedMetaBlob, err := store.Get(hashStringBlob) // Corrected Get call
	if err != nil {
		t.Fatalf("Get blob failed: %v", err)
	}
	retrievedMetaV2Blob, ok := retrievedMetaBlob.(*metav2.FileMetaV2) // Changed type assertion
	if !ok {
		t.Fatalf("Retrieved meta is not *metav2.FileMetaV2 type")
	}
	if len(retrievedMetaV2Blob.BlobHashes) != len(blobHashes) {
		t.Errorf("Blob hashes count mismatch. got %d, want %d", len(retrievedMetaV2Blob.BlobHashes), len(blobHashes))
	}
	for i, hash := range blobHashes {
		if retrievedMetaV2Blob.BlobHashes[i] != hash {
			t.Errorf("Blob hash mismatch at index %d. got %s, want %s", i, retrievedMetaV2Blob.BlobHashes[i], hash)
		}
	}

	// Test Update
	metaV2Embedded.FileSize = 200 // Corrected field name
	metaV2Embedded.MimeType = "image/png"
	err = store.Update(hashStringEmbedded, metaV2Embedded) // Corrected ID
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify Update
	updatedMetaEmbedded, err := store.Get(hashStringEmbedded) // Corrected ID
	if err != nil {
		t.Fatalf("Get after update failed: %v", err)
	}
	updatedMetaV2Embedded, ok := updatedMetaEmbedded.(*metav2.FileMetaV2) // Changed type assertion
	if !ok {
		t.Fatalf("Updated meta is not *metav2.FileMetaV2 type")
	}
	if updatedMetaV2Embedded.FileSize != 200 { // Corrected field name
		t.Errorf("Updated size mismatch. got %d, want %d", updatedMetaV2Embedded.FileSize, 200)
	}
	if updatedMetaV2Embedded.MimeType != "image/png" {
		t.Errorf("Updated mime type mismatch. got %s, want %s", updatedMetaV2Embedded.MimeType, "image/png")
	}

	// Test Exists
	exists, err := store.Exists(hashStringEmbedded) // Corrected ID
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Errorf("Exists returned false for a meta that should exist")
	}

	// Test Delete
	err = store.Delete(hashStringEmbedded) // Corrected ID
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify Delete
	exists, err = store.Exists(hashStringEmbedded) // Corrected ID
	if err != nil {
		t.Fatalf("Exists after delete failed: %v", err)
	}
	if exists {
		t.Errorf("Exists returned true for a meta that should have been deleted")
	}

	// Test Get non-existent
	_, err = store.Get(hashStringEmbedded) // Corrected ID
	if err == nil {
		t.Errorf("Get non-existent did not return an error")
	}

	// Test unsupported version
	// This requires creating a dummy meta with an unsupported version
	// For now, we'll skip this as it's hard to mock without changing the actual FileMetaV2 struct.
}
