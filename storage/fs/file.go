/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package fs

import (
	"io"
	"io/fs"

	index_interface "github.com/origadmin/runtime/interfaces/storage/index"
	meta_interface "github.com/origadmin/runtime/interfaces/storage/meta"
	blob_interface "github.com/origadmin/runtime/interfaces/storage/blob"
)

// file implements fs.File interface.
type file struct {
	reader   io.Reader
	closed   bool
	node     *index_interface.IndexNode
	fileMeta meta_interface.FileMeta
}

// NewFile creates a new fs.File instance.
// It takes an io.Reader for the file content, the IndexNode, and the FileMeta.
func NewFile(node *index_interface.IndexNode, fileMeta meta_interface.FileMeta, blobStore blob_interface.BlobStore) (fs.File, error) {
	reader, err := fileMeta.ContentReader(blobStore)
	if err != nil {
		return nil, err
	}
	return &file{
		reader:   reader,
		node:     node,
		fileMeta: fileMeta,
	}, nil
}

func (f *file) Read(p []byte) (int, error) {
	if f.closed {
		return 0, fs.ErrClosed
	}
	return f.reader.Read(p)
}

func (f *file) Close() error {
	f.closed = true
	// TODO: Potentially close the underlying reader if it's a resource that needs closing
	return nil
}

func (f *file) Stat() (fs.FileInfo, error) {
	if f.closed {
		return nil, fs.ErrClosed
	}
	return NewFileInfo(f.node, f.fileMeta), nil
}