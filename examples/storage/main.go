package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	// 1. Setup a temporary directory for the local storage base path.
	baseDir, err := os.MkdirTemp("", "storage_server_test")
	if err != nil {
		fmt.Printf("Failed to create temp directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(baseDir) // Clean up the temporary directory.
	fmt.Printf("Using temporary directory for storage: %s\n", baseDir)

	// 2. Initialize the LocalStorage service.
	fs, err := NewLocalStorage(baseDir)
	if err != nil {
		fmt.Printf("Failed to initialize storage service: %v\n", err)
		os.Exit(1)
	}

	// 3. Demonstrate storage operations.
	fmt.Println("--- Testing Storage Operations ---")

	// Create a directory.
	dirPath := "my-test-dir"
	fmt.Printf("Creating directory: %s\n", dirPath)
	if err := fs.Mkdir(dirPath); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		os.Exit(1)
	}

	// Put a file inside the new directory.
	filePath := "my-test-dir/hello.txt"
	fileContent := []byte("Hello from OrigAdmin Storage!")
	fmt.Printf("Putting file: %s\n", filePath)
	if err := fs.Put(filePath, fileContent); err != nil {
		fmt.Printf("Failed to put file: %v\n", err)
		os.Exit(1)
	}

	// List files in the directory.
	fmt.Printf("Listing files in: %s\n", dirPath)
	files, err := fs.List(dirPath)
	if err != nil {
		fmt.Printf("Failed to list files: %v\n", err)
		os.Exit(1)
	}
	for _, file := range files {
		fmt.Printf("  - Found: %s (IsDir: %v, Size: %d)\n", file.Path, file.Metadata["is_dir"], file.Size)
	}

	// Stat the file.
	fmt.Printf("Stating file: %s\n", filePath)
	info, err := fs.Stat(filePath)
	if err != nil {
		fmt.Printf("Failed to stat file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  - Stat result: Path=%s, Size=%d, ModTime=%s\n", info.Path, info.Size, info.ModTime)

	// Get the file content.
	fmt.Printf("Getting file content: %s\n", filePath)
	readContent, err := fs.Get(filePath)
	if err != nil {
		fmt.Printf("Failed to get file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  - Content: \"%s\"\n", string(readContent))

	// Rename the file.
	newFilePath := "my-test-dir/hello_renamed.txt"
	fmt.Printf("Renaming file from %s to %s\n", filePath, newFilePath)
	if err := fs.Rename(filePath, newFilePath); err != nil {
		fmt.Printf("Failed to rename file: %v\n", err)
		os.Exit(1)
	}

	// Verify rename by listing again.
	fmt.Printf("Listing files after rename in: %s\n", dirPath)
	files, err = fs.List(dirPath)
	if err != nil {
		fmt.Printf("Failed to list files after rename: %v\n", err)
		os.Exit(1)
	}
	for _, file := range files {
		fmt.Printf("  - Found: %s\n", file.Path)
	}

	// Delete the file.
	fmt.Printf("Deleting file: %s\n", newFilePath)
	if err := fs.Delete(newFilePath); err != nil {
		fmt.Printf("Failed to delete file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("--- Storage Operations Test Complete ---")

	// Keep the application running for a moment to see logs if needed.
	<-context.Background().Done()
}
