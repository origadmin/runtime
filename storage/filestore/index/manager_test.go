/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package index

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/origadmin/runtime/interfaces/storage/components/index"
)

func TestManager(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "indexmanager-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new Manage
	manager, err := NewManager(tempDir, nil)
	if err != nil {
		t.Fatalf("Failed to create Manage: %v", err)
	}
	defer manager.Close()

	// 1. Test Root Node Creation and GetNodeByPath("/")
	rootNode, err := manager.GetNodeByPath("/")
	if err != nil {
		t.Fatalf("Failed to get root node: %v", err)
	}
	if rootNode.Name != "/" || rootNode.NodeType != index.Directory || rootNode.ParentID != "" {
		t.Errorf("Root node has incorrect properties: %+v", rootNode)
	}

	// 2. Test CreateNode - Directory
	dirNode := &index.Node{
		ParentID: rootNode.NodeID,
		Name:     "my_dir",
		NodeType: index.Directory,
		Mode:     os.ModeDir | 0755,
		OwnerID:  "user1",
		GroupID:  "group1",
		Atime:    time.Now(),
		Mtime:    time.Now(),
		Ctime:    time.Now(),
	}
	if err := manager.CreateNode(dirNode); err != nil {
		t.Fatalf("Failed to create directory node: %v", err)
	}

	// Verify directory node by path
	retrievedDirNode, err := manager.GetNodeByPath(filepath.Join(rootNode.Name, dirNode.Name))
	if err != nil {
		t.Fatalf("Failed to get directory node by path: %v", err)
	}
	if retrievedDirNode.NodeID != dirNode.NodeID || retrievedDirNode.Name != "my_dir" {
		t.Errorf("Retrieved directory node mismatch: %+v", retrievedDirNode)
	}

	// 3. Test CreateNode - File
	fileNode := &index.Node{
		ParentID: dirNode.NodeID,
		Name:     "my_file.txt",
		NodeType: index.File,
		Mode:     0644,
		OwnerID:  "user1",
		GroupID:  "group1",
		Atime:    time.Now(),
		Mtime:    time.Now(),
		Ctime:    time.Now(),
	}
	if err := manager.CreateNode(fileNode); err != nil {
		t.Fatalf("Failed to create file node: %v", err)
	}

	// Verify file node by path
	retrievedFileNode, err := manager.GetNodeByPath(filepath.Join(dirNode.Name, fileNode.Name))
	if err != nil {
		t.Fatalf("Failed to get file node by path: %v", err)
	}
	if retrievedFileNode.NodeID != fileNode.NodeID {
		t.Errorf("Retrieved file node mismatch: %+v", retrievedFileNode)
	}

	// 4. Test ListChildren
	children, err := manager.ListChildren(dirNode.NodeID)
	if err != nil {
		t.Fatalf("Failed to list children: %v", err)
	}
	if len(children) != 1 || children[0].NodeID != fileNode.NodeID {
		t.Errorf("Children list mismatch. Expected 1 child, got %+v", children)
	}

	// 5. Test UpdateNode
	fileNode.Mode = 0777
	fileNode.Mtime = time.Now()
	if err := manager.UpdateNode(fileNode); err != nil {
		t.Fatalf("Failed to update node: %v", err)
	}
	updatedFileNode, err := manager.GetNode(fileNode.NodeID)
	if err != nil {
		t.Fatalf("Failed to get updated node: %v", err)
	}
	if updatedFileNode.Mode != 0777 {
		t.Errorf("Node mode not updated. Got %o, want %o", updatedFileNode.Mode, 0777)
	}

	// 6. Test MoveNode - Rename file
	newFileName := "new_file_name.txt"
	if err := manager.MoveNode(fileNode.NodeID, dirNode.NodeID, newFileName); err != nil {
		t.Fatalf("Failed to rename file: %v", err)
	}

	// Verify old path is gone
	_, err = manager.GetNodeByPath(filepath.Join(dirNode.Name, fileNode.Name))
	if !os.IsNotExist(err) {
		t.Errorf("Old file path still exists or unexpected error: %v", err)
	}

	// Verify new path exists and node is correct
	movedFileNode, err := manager.GetNodeByPath(filepath.Join(dirNode.Name, newFileName))
	if err != nil {
		t.Fatalf("Failed to get moved file by new path: %v", err)
	}
	if movedFileNode.NodeID != fileNode.NodeID || movedFileNode.Name != newFileName {
		t.Errorf("Moved file node mismatch: %+v", movedFileNode)
	}

	// 7. Test DeleteNode - File
	if err := manager.DeleteNode(movedFileNode.NodeID); err != nil {
		t.Fatalf("Failed to delete file node: %v", err)
	}

	// Verify file is gone
	_, err = manager.GetNode(movedFileNode.NodeID)
	if !os.IsNotExist(err) {
		t.Errorf("Deleted file node still exists or unexpected error: %v", err)
	}
	_, err = manager.GetNodeByPath(filepath.Join(dirNode.Name, newFileName))
	if !os.IsNotExist(err) {
		t.Errorf("Deleted file path still exists or unexpected error: %v", err)
	}

	// 8. Test DeleteNode - Non-empty directory (should fail)
	childDir := &index.Node{
		ParentID: dirNode.NodeID,
		Name:     "child_dir",
		NodeType: index.Directory,
		Mode:     os.ModeDir | 0755,
		OwnerID:  "user1",
		GroupID:  "group1",
		Atime:    time.Now(),
		Mtime:    time.Now(),
		Ctime:    time.Now(),
	}
	if err := manager.CreateNode(childDir); err != nil {
		t.Fatalf("Failed to create child directory: %v", err)
	}

	if err := manager.DeleteNode(dirNode.NodeID); err == nil {
		t.Error("Expected error when deleting non-empty directory, got nil")
	}

	// 9. Test DeleteNode - Empty directory
	if err := manager.DeleteNode(childDir.NodeID); err != nil {
		t.Fatalf("Failed to delete empty child directory: %v", err)
	}

	if err := manager.DeleteNode(dirNode.NodeID); err != nil {
		t.Fatalf("Failed to delete empty parent directory: %v", err)
	}

	// Verify directories are gone
	_, err = manager.GetNodeByPath(filepath.Join(rootNode.Name, dirNode.Name))
	if !os.IsNotExist(err) {
		t.Errorf("Deleted directory path still exists or unexpected error: %v", err)
	}

	// 10. Test CreateNode - Duplicate path (should fail)
	duplicateFileNode := &index.Node{
		ParentID: rootNode.NodeID,
		Name:     "duplicate.txt",
		NodeType: index.File,
		Mode:     0644,
	}
	if err := manager.CreateNode(duplicateFileNode); err != nil {
		t.Fatalf("Failed to create first duplicate file: %v", err)
	}

	duplicateFileNode2 := &index.Node{
		ParentID: rootNode.NodeID,
		Name:     "duplicate.txt",
		NodeType: index.File,
		Mode:     0644,
	}
	if err := manager.CreateNode(duplicateFileNode2); err == nil {
		t.Error("Expected error when creating duplicate path, got nil")
	}
}
