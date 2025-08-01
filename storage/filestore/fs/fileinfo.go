/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package fs

import (
	"io/fs"
	"time"

	indexiface "github.com/origadmin/runtime/interfaces/storage/components/index"
	metaiface "github.com/origadmin/runtime/interfaces/storage/components/meta"
)

// fileInfo implements fs.FileInfo by combining IndexNode and FileMeta.
type fileInfo struct {
	node     *indexiface.Node
	fileMeta metaiface.FileMeta // Optional, only for files
}

// NewFileInfo creates an fs.FileInfo from an IndexNode and optional FileMeta.
func NewFileInfo(node *indexiface.Node, fileMeta metaiface.FileMeta) fs.FileInfo {
	return &fileInfo{node: node, fileMeta: fileMeta}
}

func (fi *fileInfo) Name() string {
	return fi.node.Name
}

func (fi *fileInfo) Size() int64 {
	if fi.node.NodeType == indexiface.File && fi.fileMeta != nil {
		return fi.fileMeta.Size()
	}
	// For directories, size is typically 0 or undefined in fs.FileInfo
	return 0
}

func (fi *fileInfo) Mode() fs.FileMode {
	return fi.node.Mode
}

func (fi *fileInfo) ModTime() time.Time {
	// This should be the modification time of the file's content (from FileMeta)
	// or the node's modification time (from IndexNode) if it's a directory or symlink.
	if fi.node.NodeType == indexiface.File && fi.fileMeta != nil {
		return fi.fileMeta.ModTime()
	}
	return fi.node.Mtime // Node's modification time
}

func (fi *fileInfo) IsDir() bool {
	return fi.node.NodeType == indexiface.Directory
}

func (fi *fileInfo) Sys() any {
	return nil // Or return fi.node for system-specific data
}
