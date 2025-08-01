/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package index

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/origadmin/runtime/interfaces/storage/components/index"
	layoutiface "github.com/origadmin/runtime/interfaces/storage/components/layout"
	metaiface "github.com/origadmin/runtime/interfaces/storage/components/meta"
	"github.com/origadmin/runtime/storage/filestore/layout"
)

const (
	pathKeyPrefix     = "_paths"
	childrenKeyPrefix = "_children"
	nodesKeyPrefix    = "nodes"
)

// Manage implements the Manager interface using the local filesystem.
type Manage struct {
	layout    layoutiface.ShardedStorage
	indexPath string // Base path for index data
}

// NewManager creates a new Manage.
func NewManager(indexPath string, metaStore metaiface.Store) (index.Manager, error) {
	// Ensure the index path exists
	if err := os.MkdirAll(indexPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create index path: %w", err)
	}

	// Initialize ShardedStorage
	ls, err := layout.NewLocalShardedStorage(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create sharded storage: %w", err)
	}

	manager := &Manage{
		layout:    ls,
		indexPath: indexPath,
	}

	// Ensure root node exists
	_, err = manager.GetNodeByPath("/")
	if err != nil && errors.Is(err, os.ErrNotExist) {
		rootNode := &index.Node{
			NodeID:   uuid.New().String(),
			ParentID: "", // Root has no parent
			Name:     "/",
			NodeType: index.Directory,
			Mode:     os.ModeDir | 0755,
			OwnerID:  "", // Placeholder
			GroupID:  "", // Placeholder
			Atime:    time.Now(),
			Mtime:    time.Now(),
			Ctime:    time.Now(),
		}

		err = manager.CreateNode(rootNode)
		if err != nil {
			return nil, fmt.Errorf("failed to create root node: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to check for root node: %w", err)
	}

	return manager, nil
}

// Close closes the file manager.
func (m *Manage) Close() error {
	// No-op since layout doesn't require explicit closing in this implementation.
	return nil
}

func (m *Manage) pathKey(path string) string {
	return filepath.Join(pathKeyPrefix, path)
}

func (m *Manage) childrenKey(nodeID string) string {
	return filepath.Join(childrenKeyPrefix, nodeID)
}

func (m *Manage) nodeKey(nodeID string) string {
	return filepath.Join(nodesKeyPrefix, nodeID)
}

// CreateNode creates a new node in the index.
func (m *Manage) CreateNode(node *index.Node) error {
	if node.NodeID == "" {
		node.NodeID = uuid.New().String()
	}

	fullPath := filepath.Join(node.ParentID, node.Name)
	if node.Name == "/" {
		fullPath = "/"
	}

	// Check if path already exists
	_, err := m.GetNodeByPath(fullPath)
	if err == nil {
		return fmt.Errorf("node with path %s already exists", fullPath)
	}
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to check path existence: %w", err)
	}

	// Store node data
	nodeBytes, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to marshal node: %w", err)
	}
	if err := m.layout.Write(m.nodeKey(node.NodeID), nodeBytes); err != nil {
		return fmt.Errorf("failed to write node data: %w", err)
	}

	// Update path index (path -> nodeID)
	if err := m.layout.Write(m.pathKey(fullPath), []byte(node.NodeID)); err != nil {
		err := m.layout.Delete(m.nodeKey(node.NodeID))
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to update path index: %w", err)
	}

	// Update children index (parentID -> []childID)
	if node.ParentID != "" {
		children, err := m.listChildrenIDs(node.ParentID)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			m.layout.Delete(m.nodeKey(node.NodeID))
			m.layout.Delete(m.pathKey(fullPath))
			return fmt.Errorf("failed to list children IDs for update: %w", err)
		}

		children = append(children, node.NodeID)
		childrenBytes, err := json.Marshal(children)
		if err != nil {
			return fmt.Errorf("failed to marshal children IDs: %w", err)
		}

		if err := m.layout.Write(m.childrenKey(node.ParentID), childrenBytes); err != nil {
			m.layout.Delete(m.nodeKey(node.NodeID))
			m.layout.Delete(m.pathKey(fullPath))
			return fmt.Errorf("failed to update children index: %w", err)
		}
	}

	return nil
}

// GetNode retrieves a node by its unique ID.
func (m *Manage) GetNode(nodeID string) (*index.Node, error) {
	nodeBytes, err := m.layout.Read(m.nodeKey(nodeID))
	if err != nil {
		return nil, err
	}

	var node index.Node
	if err := json.Unmarshal(nodeBytes, &node); err != nil {
		return nil, fmt.Errorf("failed to unmarshal node: %w", err)
	}
	return &node, nil
}

// GetNodeByPath retrieves a node by its full path.
func (m *Manage) GetNodeByPath(path string) (*index.Node, error) {
	nodeIDBytes, err := m.layout.Read(m.pathKey(path))
	if err != nil {
		return nil, err
	}
	return m.GetNode(string(nodeIDBytes))
}

// UpdateNode updates an existing node's data.
func (m *Manage) UpdateNode(node *index.Node) error {
	nodeBytes, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to marshal node for update: %w", err)
	}
	return m.layout.Write(m.nodeKey(node.NodeID), nodeBytes)
}

// DeleteNode removes a node from the index.
func (m *Manage) DeleteNode(nodeID string) error {
	node, err := m.GetNode(nodeID)
	if err != nil {
		return err
	}

	if node.NodeType == index.Directory {
		children, err := m.ListChildren(node.NodeID)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to list children for directory deletion check: %w", err)
		}
		if len(children) > 0 {
			return fmt.Errorf("cannot delete non-empty directory: %s", node.Name)
		}
	}

	if err := m.layout.Delete(m.nodeKey(nodeID)); err != nil {
		return fmt.Errorf("failed to delete node data: %w", err)
	}

	fullPath := filepath.Join(node.ParentID, node.Name)
	if node.Name == "/" {
		fullPath = "/"
	}
	m.layout.Delete(m.pathKey(fullPath))

	if node.ParentID != "" {
		children, err := m.listChildrenIDs(node.ParentID)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			return fmt.Errorf("failed to list children for update: %w", err)
		}

		newChildren := make([]string, 0, len(children))
		for _, childID := range children {
			if childID != nodeID {
				newChildren = append(newChildren, childID)
			}
		}

		if len(newChildren) == 0 {
			return m.layout.Delete(m.childrenKey(node.ParentID))
		}

		childrenBytes, err := json.Marshal(newChildren)
		if err != nil {
			return fmt.Errorf("failed to marshal updated children list: %w", err)
		}
		return m.layout.Write(m.childrenKey(node.ParentID), childrenBytes)
	}

	return nil
}

// ListChildren retrieves all immediate children of a directory node.
func (m *Manage) ListChildren(parentID string) ([]*index.Node, error) {
	childIDs, err := m.listChildrenIDs(parentID)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	childrenNodes := make([]*index.Node, 0, len(childIDs))
	for _, id := range childIDs {
		node, err := m.GetNode(id)
		if err != nil {
			continue
		}
		childrenNodes = append(childrenNodes, node)
	}
	return childrenNodes, nil
}

func (m *Manage) listChildrenIDs(parentID string) ([]string, error) {
	childrenBytes, err := m.layout.Read(m.childrenKey(parentID))
	if err != nil {
		return nil, err
	}

	var childIDs []string
	if err := json.Unmarshal(childrenBytes, &childIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal children IDs: %w", err)
	}
	return childIDs, nil
}

// MoveNode moves a node to a new parent directory and/or new name.
func (m *Manage) MoveNode(nodeID string, newParentID string, newName string) error {
	node, err := m.GetNode(nodeID)
	if err != nil {
		return fmt.Errorf("node not found: %w", err)
	}

	oldParentID := node.ParentID
	oldName := node.Name
	oldFullPath := filepath.Join(oldParentID, oldName)
	if oldName == "/" {
		oldFullPath = "/"
	}
	newFullPath := filepath.Join(newParentID, newName)

	_, err = m.GetNodeByPath(newFullPath)
	if err == nil {
		return fmt.Errorf("target path %s already exists", newFullPath)
	}
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to check target path existence: %w", err)
	}

	// Simplified move operation without full transaction support.
	// A failure during this process can lead to an inconsistent state.

	// 1. Remove from old parent's children list
	if oldParentID != "" {
		if err := m.removeChildFromParent(nodeID, oldParentID); err != nil {
			return fmt.Errorf("failed to remove from old parent's children list: %w", err)
		}
	}

	// 2. Add to new parent's children list
	if newParentID != "" {
		if err := m.addChildToParent(nodeID, newParentID); err != nil {
			// Try to add back to old parent to rollback
			m.addChildToParent(nodeID, oldParentID)
			return fmt.Errorf("failed to add to new parent's children list: %w", err)
		}
	}

	// 3. Update path index
	m.layout.Delete(m.pathKey(oldFullPath))
	if err := m.layout.Write(m.pathKey(newFullPath), []byte(nodeID)); err != nil {
		// Rollback children change and path deletion
		m.layout.Delete(m.pathKey(newFullPath))
		m.layout.Write(m.pathKey(oldFullPath), []byte(nodeID))
		if newParentID != "" {
			m.removeChildFromParent(nodeID, newParentID)
		}
		if oldParentID != "" {
			m.addChildToParent(nodeID, oldParentID)
		}
		return fmt.Errorf("failed to update path index for move: %w", err)
	}

	// 4. Update node's own properties
	node.ParentID = newParentID
	node.Name = newName
	node.Mtime = time.Now()
	if err := m.UpdateNode(node); err != nil {
		// Rollback all previous steps
		m.layout.Delete(m.pathKey(newFullPath))
		m.layout.Write(m.pathKey(oldFullPath), []byte(nodeID))
		if newParentID != "" {
			m.removeChildFromParent(nodeID, newParentID)
		}
		if oldParentID != "" {
			m.addChildToParent(nodeID, oldParentID)
		}
		return fmt.Errorf("failed to write updated node data for move: %w", err)
	}

	return nil
}

func (m *Manage) addChildToParent(childID, parentID string) error {
	children, err := m.listChildrenIDs(parentID)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	children = append(children, childID)
	childrenBytes, err := json.Marshal(children)
	if err != nil {
		return err
	}
	return m.layout.Write(m.childrenKey(parentID), childrenBytes)
}

func (m *Manage) removeChildFromParent(childID, parentID string) error {
	children, err := m.listChildrenIDs(parentID)
	if err != nil {
		return err
	}
	newChildren := make([]string, 0, len(children))
	for _, id := range children {
		if id != childID {
			newChildren = append(newChildren, id)
		}
	}
	if len(newChildren) == 0 {
		return m.layout.Delete(m.childrenKey(parentID))
	}
	childrenBytes, err := json.Marshal(newChildren)
	if err != nil {
		return err
	}
	return m.layout.Write(m.childrenKey(parentID), childrenBytes)
}
