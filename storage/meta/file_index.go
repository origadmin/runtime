// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"io/fs"
	"time"
)

type FileIndexEntry struct {
	EntryName   string      `json:"name" msgpack:"n"` // Renamed from Name
	Hash        string      `json:"hash" msgpack:"h"` // 指向 FileMeta 或 DirectoryIndex 的 hash
	FileMode    fs.FileMode `json:"mode" msgpack:"m"`
	IsDirectory bool        `json:"is_dir" msgpack:"d"`
}

type DirectoryIndex struct {
	Version int32            `json:"version" msgpack:"v"`
	Path    string           `json:"path" msgpack:"p"`
	Entries []FileIndexEntry `json:"entries" msgpack:"e"`
}

// --- fs.DirEntry interface implementation ---

// Name returns the name of the file (or subdirectory) described by the entry.
func (f FileIndexEntry) Name() string {
	return f.EntryName // Use the renamed field
}

// IsDir reports whether the entry describes a directory.
func (f FileIndexEntry) IsDir() bool {
	return f.IsDirectory
}

// Type returns the type of file system object described by the entry.
func (f FileIndexEntry) Type() fs.FileMode {
	return f.FileMode.Type()
}

// Info returns the FileInfo for the file or subdirectory described by the entry.
func (f FileIndexEntry) Info() (fs.FileInfo, error) {
	// This returns the entry itself, which partially implements FileInfo.
	// For complete info like size, a full Stat() is needed.
	return f, nil
}

// --- fs.FileInfo interface implementation ---

// Size returns the size of the file. For directories, it is 0.
func (f FileIndexEntry) Size() int64 {
	// Size is not stored in the directory entry. A full Stat() is required.
	return 0
}

// Mode returns the file mode bits.
func (f FileIndexEntry) Mode() fs.FileMode {
	return f.FileMode
}

// ModTime returns the modification time.
func (f FileIndexEntry) ModTime() time.Time {
	// Modification time is not stored in this version.
	return time.Time{}
}

// Sys returns underlying data source (can be nil).
func (f FileIndexEntry) Sys() interface{} {
	return nil
}
