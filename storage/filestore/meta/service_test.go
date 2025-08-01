package meta

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"testing"
)

// setupTest creates a temporary directory for the test and returns its path.
// It also returns a cleanup function to remove the directory.
func setupTest(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "meta_test_")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return tmpDir, func() {
		os.RemoveAll(tmpDir)
	}
}

type dummyBlobStorage struct {
	data map[string][]byte
}

func (d *dummyBlobStorage) Store(content []byte) (string, error) {
	if d.data == nil {
		d.data = make(map[string][]byte)
	}
	h := sha256.New()
	h.Write(content)
	hash := hex.EncodeToString(h.Sum(nil))
	d.data[hash] = content
	return hash, nil
}

func (d *dummyBlobStorage) Retrieve(hash string) ([]byte, error) {
	if d.data == nil {
		return nil, os.ErrNotExist
	}
	content, ok := d.data[hash]
	if !ok {
		return nil, os.ErrNotExist
	}
	return content, nil
}

func TestNewMeta(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	blob := &dummyBlobStorage{}
	_, err := New(tmpDir, blob)
	if err != nil {
		t.Fatalf("NewMeta failed: %v", err)
	}

	// No specific checks for root metaFile or directory structure as meta now only handles files.
}

func TestWriteFile(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	blob := &dummyBlobStorage{}
	m, err := New(tmpDir, blob)
	if err != nil {
		t.Fatalf("NewMeta failed: %v", err)
	}

	// Write a small metaFile
	content := "Hello, world!"
	err = m.WriteFile("/test.txt", bytes.NewReader([]byte(content)), 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Verify metaFile exists and content is correct
	file, err := m.Open("/test.txt")
	if err != nil {
		t.Fatalf("Open /test.txt failed: %v", err)
	}
	defer file.Close()

	readContent, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("ReadAll /test.txt failed: %v", err)
	}
	if string(readContent) != content {
		t.Fatalf("expected '%s', got '%s'", content, string(readContent))
	}

	// Overwrite existing metaFile
	newContent := "New content here."
	err = m.WriteFile("/test.txt", bytes.NewReader([]byte(newContent)), 0644)
	if err != nil {
		t.Fatalf("Overwrite WriteFile failed: %v", err)
	}

	file, err = m.Open("/test.txt")
	if err != nil {
		t.Fatalf("Open /test.txt after overwrite failed: %v", err)
	}
	defer file.Close()

	readContent, err = io.ReadAll(file)
	if err != nil {
		t.Fatalf("ReadAll /test.txt after overwrite failed: %v", err)
	}
	if string(readContent) != newContent {
		t.Fatalf("expected '%s', got '%s' after overwrite", newContent, string(readContent))
	}

	// Write a larger metaFile to test chunking (e.g., 5MB)
	largeContent := make([]byte, DefaultBlockSize+1024) // Slightly larger than one block
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}
	err = m.WriteFile("/large.bin", bytes.NewReader(largeContent), 0644)
	if err != nil {
		t.Fatalf("WriteFile large.bin failed: %v", err)
	}

	file, err = m.Open("/large.bin")
	if err != nil {
		t.Fatalf("Open /large.bin failed: %v", err)
	}
	defer file.Close()

	readContent, err = io.ReadAll(file)
	if err != nil {
		t.Fatalf("ReadAll /large.bin failed: %v", err)
	}
	if !bytes.Equal(readContent, largeContent) {
		t.Fatalf("large metaFile content mismatch")
	}
}

func TestOpen(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	blob := &dummyBlobStorage{}
	m, err := New(tmpDir, blob)
	if err != nil {
		t.Fatalf("NewMeta failed: %v", err)
	}

	// Write a metaFile
	content := "This is the metaFile content."
	err = m.WriteFile("/my_file.txt", bytes.NewReader([]byte(content)), 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Open and read the metaFile
	file, err := m.Open("/my_file.txt")
	if err != nil {
		t.Fatalf("Open /my_file.txt failed: %v", err)
	}
	defer file.Close()

	readContent, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("ReadAll /my_file.txt failed: %v", err)
	}
	if string(readContent) != content {
		t.Fatalf("expected '%s', got '%s'", content, string(readContent))
	}

	// Attempt to open a non-existent metaFile
	_, err = m.Open("/nonexistent_file.txt")
	if !os.IsNotExist(err) {
		t.Fatalf("expected ErrNotExist for non-existent metaFile, got %v", err)
	}
}

func TestStat(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	blob := &dummyBlobStorage{}
	m, err := New(tmpDir, blob)
	if err != nil {
		t.Fatalf("NewMeta failed: %v", err)
	}

	// Stat a created metaFile
	content := "metaFile content"
	m.WriteFile("/testfile.txt", bytes.NewReader([]byte(content)), 0644)
	fileInfo, err := m.Stat("/testfile.txt")
	if err != nil {
		t.Fatalf("Stat /testfile.txt failed: %v", err)
	}
	if fileInfo.Name() != "testfile.txt" || fileInfo.IsDir() || fileInfo.Size() != int64(len(content)) || fileInfo.Mode() != 0644 {
		t.Fatalf("expected testfile info, got %+v", fileInfo)
	}

	// Stat non-existent path
	_, err = m.Stat("/nonexistent")
	if !os.IsNotExist(err) {
		t.Fatalf("expected ErrNotExist for non-existent path, got %v", err)
	}
}
