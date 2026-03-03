package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// 1. Setup a temporary directory for the local storage base path.
	baseDir, err := os.MkdirTemp("", "storage_server_test")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(baseDir)
	}() // Clean up the temporary directory.
	fmt.Printf("Using temporary directory for storage: %s\n", baseDir)

	// 2. Initialize the LocalStorage service.
	fs, err := NewLocalStorage(baseDir)
	if err != nil {
		return fmt.Errorf("failed to initialize storage service: %v", err)
	}

	// 3. Demonstrate storage operations.
	fmt.Println("--- Testing Storage Operations ---")

	// Create a directory.
	dirPath := "my-test-dir"
	fmt.Printf("Creating directory: %s\n", dirPath)
	if err := fs.Mkdir(dirPath); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Put a file inside the new directory.
	filePath := "my-test-dir/hello.txt"
	fileContent := []byte("Hello from OrigAdmin Storage!")
	fmt.Printf("Putting file: %s\n", filePath)
	if err := fs.Put(filePath, fileContent); err != nil {
		return fmt.Errorf("failed to put file: %v", err)
	}

	// List files in the directory.
	fmt.Printf("Listing files in: %s\n", dirPath)
	files, err := fs.List(dirPath)
	if err != nil {
		return fmt.Errorf("failed to list files: %v", err)
	}
	for _, file := range files {
		fmt.Printf("  - Found: %s (IsDir: %v, Size: %d)\n", file.Path, file.Metadata["is_dir"], file.Size)
	}

	// Stat the file.
	fmt.Printf("Stating file: %s\n", filePath)
	info, err := fs.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to stat file: %v", err)
	}
	fmt.Printf("  - Stat result: Path=%s, Size=%d, ModTime=%s\n", info.Path, info.Size, info.ModTime)

	// Get the file content.
	fmt.Printf("Getting file content: %s\n", filePath)
	readContent, err := fs.Get(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file: %v", err)
	}
	fmt.Printf("  - Content: \"%s\"\n", string(readContent))

	// Rename the file.
	newFilePath := "my-test-dir/hello_renamed.txt"
	fmt.Printf("Renaming file from %s to %s\n", filePath, newFilePath)
	if err := fs.Rename(filePath, newFilePath); err != nil {
		return fmt.Errorf("failed to rename file: %v", err)
	}

	// Verify rename by listing again.
	fmt.Printf("Listing files after rename in: %s\n", dirPath)
	files, err = fs.List(dirPath)
	if err != nil {
		return fmt.Errorf("failed to list files after rename: %v", err)
	}
	for _, file := range files {
		fmt.Printf("  - Found: %s\n", file.Path)
	}

	// Delete the file.
	fmt.Printf("Deleting file: %s\n", newFilePath)
	if err := fs.Delete(newFilePath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	fmt.Println("--- Storage Operations Test Complete ---")

	// Keep the application running for a moment to see logs if needed.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Simply finish the demo.
	_ = ctx
	return nil
}
