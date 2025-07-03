/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package index

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/origadmin/runtime/storage/layout"
	index_interface "github.com/origadmin/runtime/interfaces/storage/index"
)

const (
	pathKeyPrefix   = "path:"
	childrenKeyPrefix = "children:"
)

// FileIndexManager implements the IndexManager interface using the local filesystem
// and BadgerDB for auxiliary indexes.
type FileIndexManager struct {
	layout    layout.ShardedStorage
	db        *badger.DB
	indexPath string // Base path for index data (nodes and badger db)
}

// NewFileIndexManager creates a new FileIndexManager.
func NewFileIndexManager(indexPath string) (*FileIndexManager, error) {
	// Ensure the index path exists
	if err := os.MkdirAll(indexPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create index path: %w", err)
	}

	// Initialize BadgerDB
	dbPath := filepath.Join(indexPath, "badger")
	opts := badger.DefaultOptions(dbPath).WithLogger(nil) // Disable BadgerDB logs for cleaner output
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open BadgerDB: %w", err)
	}

	// Initialize ShardedStorage for nodes
	nodesPath := filepath.Join(indexPath, "nodes")
	ls, err := layout.NewLocalShardedStorage(nodesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create sharded storage for nodes: %w", err)
	}

	manager := &FileIndexManager{
		layout:    ls,
		db:        db,
		indexPath: indexPath,
	}

	// Ensure root node exists
	_, err = manager.GetNodeByPath("/")
	if err != nil && err == badger.ErrKeyNotFound {
		// Create root node if it doesn't exist
		rootNode := &IndexNode{
			NodeID:    uuid.New().String(),
			ParentID:  "", // Root has no parent
			Name:      "/",
			NodeType:  Directory,
			Mode:      os.ModeDir | 0755,
			OwnerID:   "", // Placeholder
			GroupID:   "", // Placeholder
			Atime:     time.Now(),
			Mtime:     time.Now(),
			Ctime:     time.Now(),
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

// Close closes the underlying BadgerDB.
func (m *FileIndexManager) Close() error {
	return m.db.Close()
}

// CreateNode creates a new node in the index.
func (m *FileIndexManager) CreateNode(node *IndexNode) error {
	if node.NodeID == "" {
		node.NodeID = uuid.New().String()
	}

	// Check if path already exists
	_, err := m.GetNodeByPath(filepath.Join(node.ParentID, node.Name))
	if err == nil {
		return fmt.Errorf("node with path %s already exists", filepath.Join(node.ParentID, node.Name))
	}
	if err != badger.ErrKeyNotFound {
		return fmt.Errorf("failed to check path existence: %w", err)
	}

	// Store node data using sharded layout
	nodeBytes, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to marshal node: %w", err)
	}
	if err := m.layout.Write(node.NodeID, nodeBytes); err != nil {
		return fmt.Errorf("failed to write node data: %w", err)
	}

	// Update path index (path -> nodeID)
	fullPath := filepath.Join(node.ParentID, node.Name)
	if err := m.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(pathKeyPrefix+fullPath), []byte(node.NodeID))
	}); err != nil {
		return fmt.Errorf("failed to update path index: %w", err)
	}

	// Update children index (parentID -> []childID)
	if node.ParentID != "" {
		parentPath := node.ParentID
		if parentPath == "/" {
			parentPath = ""
		}
		parentNode, err := m.GetNodeByPath(parentPath)
		if err != nil {
			return fmt.Errorf("failed to get parent node for children update: %w", err)
		}

		childrenKey := []byte(childrenKeyPrefix + parentNode.NodeID)
		return m.db.Update(func(txn *badger.Txn) error {
			item, err := txn.Get(childrenKey)
			var children []string
			if err == nil {
				val, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}
				if err := json.Unmarshal(val, &children); err != nil {
					return err
				}
			} else if err != badger.ErrKeyNotFound {
				return err
			}

			children = append(children, node.NodeID)
			childrenBytes, err := json.Marshal(children)
			if err != nil {
				return err
			}
			return txn.Set(childrenKey, childrenBytes)
		})
	}

	return nil
}

// GetNode retrieves a node by its unique ID.
func (m *FileIndexManager) GetNode(nodeID string) (*IndexNode, error) {
	nodeBytes, err := m.layout.Read(nodeID)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, badger.ErrKeyNotFound // Translate os.ErrNotExist to badger.ErrKeyNotFound for consistency
		}
		return nil, fmt.Errorf("failed to read node data: %w", err)
	}

	var node IndexNode
	if err := json.Unmarshal(nodeBytes, &node); err != nil {
		return nil, fmt.Errorf("failed to unmarshal node: %w", err)
	}
	return &node, nil
}

// GetNodeByPath retrieves a node by its full path.
func (m *FileIndexManager) GetNodeByPath(path string) (*IndexNode, error) {
	var nodeID string
	err := m.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(pathKeyPrefix + path))
		if err != nil {
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		nodeID = string(val)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return m.GetNode(nodeID)
}

// UpdateNode updates an existing node's data.
func (m *FileIndexManager) UpdateNode(node *IndexNode) error {
	// For simplicity, we only allow updating the node data, not its path or parent.
	// Changing path/parent would require updating path and children indexes.
	nodeBytes, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to marshal node for update: %w", err)
	}
	return m.layout.Write(node.NodeID, nodeBytes) // Overwrite the existing node file
}

// DeleteNode removes a node from the index.
func (m *FileIndexManager) DeleteNode(nodeID string) error {
	node, err := m.GetNode(nodeID)
	if err != nil {
		return err
	}

	// Check if it's a directory and not empty
	if node.NodeType == Directory {
		children, err := m.ListChildren(node.NodeID)
		if err != nil {
			return fmt.Errorf("failed to list children for directory deletion check: %w", err)
		}
		if len(children) > 0 {
			return fmt.Errorf("cannot delete non-empty directory: %s", node.Name)
		}
	}

	// Delete from sharded storage
	if err := m.layout.Delete(nodeID); err != nil {
		return fmt.Errorf("failed to delete node data: %w", err)
	}

	// Delete from path index
	fullPath := filepath.Join(node.ParentID, node.Name)
	if err := m.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(pathKeyPrefix + fullPath))
	}); err != nil {
		return fmt.Errorf("failed to delete from path index: %w", err)
	}

	// Remove from parent's children index
	if node.ParentID != "" {
		childrenKey := []byte(childrenKeyPrefix + node.ParentID)
		return m.db.Update(func(txn *badger.Txn) error {
			item, err := txn.Get(childrenKey)
			var children []string
			if err == nil {
				val, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}
				if err := json.Unmarshal(val, &children); err != nil {
					return err
				}
			} else if err != badger.ErrKeyNotFound {
				return err
			}

			// Remove the deleted nodeID from the children list
			newChildren := make([]string, 0, len(children))
			for _, childID := range children {
				if childID != nodeID {
					newChildren = append(newChildren, childID)
				}
			}

			if len(newChildren) == 0 {
				return txn.Delete(childrenKey) // Delete if no children left
			}
			childrenBytes, err := json.Marshal(newChildren)
			if err != nil {
				return err
			}
			return txn.Set(childrenKey, childrenBytes)
		})
	}

	return nil
}

// ListChildren retrieves all immediate children of a directory node.
func (m *FileIndexManager) ListChildren(parentID string) ([]*IndexNode, error) {
	childrenKey := []byte(childrenKeyPrefix + parentID)
	var childIDs []string
	err := m.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(childrenKey)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil // No children
			}
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return json.Unmarshal(val, &childIDs)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get children IDs: %w", err)
	}

	var childrenNodes []*IndexNode
	for _, id := range childIDs {
		node, err := m.GetNode(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get child node %s: %w", id, err)
		}
		childrenNodes = append(childrenNodes, node)
	}
	return childrenNodes, nil
}

// MoveNode moves a node to a new parent directory and/or new name.
func (m *FileIndexManager) MoveNode(nodeID string, newParentID string, newName string) error {
	// Get the node to move
	node, err := m.GetNode(nodeID)
	if err != nil {
		return fmt.Errorf("node not found: %w", err)
	}

	oldFullPath := filepath.Join(node.ParentID, node.Name)
	newFullPath := filepath.Join(newParentID, newName)

	// Check if the new path already exists
	_, err = m.GetNodeByPath(newFullPath)
	if err == nil {
		return fmt.Errorf("target path %s already exists", newFullPath)
	}
	if err != badger.ErrKeyNotFound {
		return fmt.Errorf("failed to check target path existence: %w", err)
	}

	// Update the node itself
	oldParentID := node.ParentID
	node.ParentID = newParentID
	node.Name = newName
	node.Mtime = time.Now() // Update modification time of the node

	// Update node data in sharded storage
	nodeBytes, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to marshal node for move: %w", err)
	}
	if err := m.layout.Write(node.NodeID, nodeBytes); err != nil {
		return fmt.Errorf("failed to write updated node data: %w", err)
	}

	// Update path index
	if err := m.db.Update(func(txn *badger.Txn) error {
		// Delete old path entry
		if err := txn.Delete([]byte(pathKeyPrefix + oldFullPath)); err != nil {
			return err
		}
		// Set new path entry
		return txn.Set([]byte(pathKeyPrefix+newFullPath), []byte(node.NodeID))
	}); err != nil {
		return fmt.Errorf("failed to update path index for move: %w", err)
	}

	// Update children indexes (remove from old parent, add to new parent)
	return m.db.Update(func(txn *badger.Txn) error {
		// Remove from old parent's children list
		if oldParentID != "" {
			oldChildrenKey := []byte(childrenKeyPrefix + oldParentID)
			item, err := txn.Get(oldChildrenKey)
			var oldChildren []string
			if err == nil {
				val, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}
				if err := json.Unmarshal(val, &oldChildren); err != nil {
					return err
				}
			} else if err != badger.ErrKeyNotFound {
				return err
			}

			newOldChildren := make([]string, 0, len(oldChildren))
			for _, childID := range oldChildren {
				if childID != nodeID {
					newOldChildren = append(newOldChildren, childID)
				}
			}
			if len(newOldChildren) == 0 {
				if err := txn.Delete(oldChildrenKey); err != nil {
					return err
				}
			} else {
				oldChildrenBytes, err := json.Marshal(newOldChildren)
				if err != nil {
					return err
				}
				if err := txn.Set(oldChildrenKey, oldChildrenBytes); err != nil {
					return err
				}
			}
		}

		// Add to new parent's children list
		if newParentID != "" {
			newChildrenKey := []byte(childrenKeyPrefix + newParentID)
			item, err := txn.Get(newChildrenKey)
			var newChildren []string
			if err == nil {
				val, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}
				if err := json.Unmarshal(val, &newChildren); err != nil {
					return err
				}
			} else if err != badger.ErrKeyNotFound {
				return err
			}

			newChildren = append(newChildren, nodeID)
			newChildrenBytes, err := json.Marshal(newChildren)
			if err != nil {
				return err
			}
			return txn.Set(newChildrenKey, newChildrenBytes)
		}
		return nil
	})
}
