package meta

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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

func TestNewMeta(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	m, err := New(tmpDir)
	if err != nil {
		t.Fatalf("NewMeta failed: %v", err)
	}

	// Verify root file exists
	rootFilePath := filepath.Join(tmpDir, "index", rootFileName)
	if _, err := os.Stat(rootFilePath); os.IsNotExist(err) {
		t.Fatalf("root file %s not created", rootFilePath)
	}

	// Verify root hash is loaded
	if m.root == nil {
		t.Fatal("root not loaded")
	}

	// Verify root directory can be read
	rootEntries, err := m.ReadDir("/")
	if err != nil {
		t.Fatalf("ReadDir / failed: %v", err)
	}
	if len(rootEntries) != 0 {
		t.Fatalf("expected empty root directory, got %d entries", len(rootEntries))
	}
}

func TestMkdir(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	m, err := New(tmpDir)
	if err != nil {
		t.Fatalf("NewMeta failed: %v", err)
	}

	// Test creating a single directory
	err = m.Mkdir("/testdir", 0755)
	if err != nil {
		t.Fatalf("Mkdir /testdir failed: %v", err)
	}

	entries, err := m.ReadDir("/")
	if err != nil {
		t.Fatalf("ReadDir / failed: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != "testdir" || !entries[0].IsDir() {
		t.Fatalf("expected one directory 'testdir', got %+v", entries)
	}

	// Test creating nested directories
	err = m.Mkdir("/a/b/c", 0755)
	if err == nil {
		t.Fatalf("Mkdir /a/b/c should fail for non-existent parent")
	}

	err = m.Mkdir("/a", 0755)
	if err != nil {
		t.Fatalf("Mkdir /a failed: %v", err)
	}
	err = m.Mkdir("/a/b", 0755)
	if err != nil {
		t.Fatalf("Mkdir /a/b failed: %v", err)
	}
	err = m.Mkdir("/a/b/c", 0755)
	if err != nil {
		t.Fatalf("Mkdir /a/b/c failed: %v", err)
	}

	// Test creating an existing directory
	err = m.Mkdir("/testdir", 0755)
	if !os.IsExist(err) {
		t.Fatalf("expected ErrExist for existing directory, got %v", err)
	}

	// Test creating directory with invalid path
	err = m.Mkdir("invalid/path", 0755)
	if err == nil {
		t.Fatalf("Mkdir invalid/path should fail")
	}
}

func TestReadDir(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	m, err := New(tmpDir)
	if err != nil {
		t.Fatalf("NewMeta failed: %v", err)
	}

	// Create some directories and files
	m.Mkdir("/dir1", 0755)
	m.Mkdir("/dir2", 0755)
	m.WriteFile("/file1.txt", bytes.NewReader([]byte("hello")), 0644)

	// Read root directory
	entries, err := m.ReadDir("/")
	if err != nil {
		t.Fatalf("ReadDir / failed: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries in root, got %d", len(entries))
	}

	foundDir1 := false
	foundDir2 := false
	foundFile1 := false
	for _, e := range entries {
		switch e.Name() {
		case "dir1":
			if !e.IsDir() {
				t.Errorf("dir1 is not a directory")
			}
			foundDir1 = true
		case "dir2":
			if !e.IsDir() {
				t.Errorf("dir2 is not a directory")
			}
			foundDir2 = true
		case "file1.txt":
			if e.IsDir() {
				t.Errorf("file1.txt is a directory")
			}
			foundFile1 = true
		}
	}
	if !foundDir1 || !foundDir2 || !foundFile1 {
		t.Errorf("missing expected entries: dir1=%t, dir2=%t, file1.txt=%t", foundDir1, foundDir2, foundFile1)
	}

	// Read non-existent directory
	_, err = m.ReadDir("/nonexistent")
	if !os.IsNotExist(err) {
		t.Fatalf("expected ErrNotExist for non-existent directory, got %v", err)
	}

	// Read a file as a directory
	_, err = m.ReadDir("/file1.txt")
	if err == nil || !strings.Contains(err.Error(), "not a directory") {
		t.Fatalf("expected 'not a directory' error for file, got %v", err)
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	m, err := New(tmpDir)
	if err != nil {
		t.Fatalf("NewMeta failed: %v", err)
	}

	// Write a small file
	content := "Hello, world!"
	err = m.WriteFile("/test.txt", bytes.NewReader([]byte(content)), 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Verify file exists and content is correct
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

	// Overwrite existing file
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

	// Write a larger file to test chunking (e.g., 5MB)
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

	readLargeContent, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("ReadAll /large.bin failed: %v", err)
	}
	if !bytes.Equal(readLargeContent, largeContent) {
		t.Fatalf("large file content mismatch")
	}

	// Attempt to write to a non-existent directory
	err = m.WriteFile("/nonexistent/file.txt", bytes.NewReader([]byte("data")), 0644)
	if !os.IsNotExist(err) {
		t.Fatalf("expected ErrNotExist for non-existent parent, got %v", err)
	}

	// Attempt to write over a directory
	err = m.WriteFile("/testdir", bytes.NewReader([]byte("data")), 0644)
	if err == nil || !strings.Contains(err.Error(), "is a directory") {
		t.Fatalf("expected 'is a directory' error when writing over dir, got %v", err)
	}
}

func TestOpen(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	m, err := New(tmpDir)
	if err != nil {
		t.Fatalf("NewMeta failed: %v", err)
	}

	// Write a file
	content := "This is the file content."
	err = m.WriteFile("/my_file.txt", bytes.NewReader([]byte(content)), 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Open and read the file
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

	// Attempt to open a non-existent file
	_, err = m.Open("/nonexistent_file.txt")
	if !os.IsNotExist(err) {
		t.Fatalf("expected ErrNotExist for non-existent file, got %v", err)
	}

	// Attempt to open a directory as a file
	m.Mkdir("/my_dir", 0755)
	_, err = m.Open("/my_dir")
	if err == nil || !strings.Contains(err.Error(), "is a directory") {
		t.Fatalf("expected 'is a directory' error when opening dir, got %v", err)
	}
}

func TestStat(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	m, err := New(tmpDir)
	if err != nil {
		t.Fatalf("NewMeta failed: %v", err)
	}

	// Stat root directory
	rootInfo, err := m.Stat("/")
	if err != nil {
		t.Fatalf("Stat / failed: %v", err)
	}
	if rootInfo.Name() != "/" || !rootInfo.IsDir() {
		t.Fatalf("expected root dir info, got %+v", rootInfo)
	}

	// Stat a created directory
	m.Mkdir("/testdir", 0755)
	dirInfo, err := m.Stat("/testdir")
	if err != nil {
		t.Fatalf("Stat /testdir failed: %v", err)
	}
	if dirInfo.Name() != "testdir" || !dirInfo.IsDir() || dirInfo.Mode() != (0755|fs.ModeDir) {
		t.Fatalf("expected testdir info, got %+v", dirInfo)
	}

	// Stat a created file
	content := "file content"
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
