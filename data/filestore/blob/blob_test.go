/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package blob

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBlobStore(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "blobstore-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new BlobStore
	store, err := New(tempDir)
	if err != nil {
		t.Fatalf("Failed to create BlobStore: %v", err)
	}

	// Test data
	data := []byte("hello world")
	// SHA256 hash of "hello world" is b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9
	expectedID := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

	// Test Write
	id, err := store.Write(data)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if id != expectedID {
		t.Errorf("Write returned incorrect ID. got %v, want %v", id, expectedID)
	}

	// Check that the file was created in the correct sharded path
	expectedPath := filepath.Join(tempDir, expectedID[:2], expectedID[2:4], expectedID)
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Blob not found at expected path: %v", expectedPath)
	}

	// Test Exists
	exists, err := store.Exists(id)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Errorf("Exists returned false for a file that should exist")
	}

	// Test Read
	readData, err := store.Read(id)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if string(readData) != string(data) {
		t.Errorf("Read returned incorrect data. got %v, want %v", string(readData), string(data))
	}

	// Test Delete
	err = store.Delete(id)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Test Exists after delete
	exists, err = store.Exists(id)
	if err != nil {
		t.Fatalf("Exists after delete failed: %v", err)
	}
	if exists {
		t.Errorf("Exists returned true for a file that should have been deleted")
	}
}
