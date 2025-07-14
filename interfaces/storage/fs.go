package storage

import (
	"io"
	"time"
)

// FileInfo describes a file or directory.
type FileInfo struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime time.Time
}

// FileOperations defines a standard interface for file and directory manipulations,
// abstracting the underlying storage mechanism (e.g., local disk, cloud storage).
type FileOperations interface {
	List(path string) ([]FileInfo, error)
	Read(path string) (io.ReadCloser, error)
	Write(path string, data io.Reader, size int64) error
	Stat(path string) (FileInfo, error)
	Mkdir(path string) error
	Delete(path string) error
	Rename(oldPath, newPath string) error
}
