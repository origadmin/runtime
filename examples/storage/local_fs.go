package main

import (
	"os"
	"path/filepath"

	"github.com/origadmin/runtime/interfaces/storage"
)

// LocalStorage provides a file system implementation based on the local disk.
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new LocalStorage instance.
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	return &LocalStorage{basePath: basePath}, nil
}

// List returns a slice of ObjectInfo for the given directory path.
func (fs *LocalStorage) List(path string) ([]*storage.ObjectInfo, error) { // Changed return type
	var files []*storage.ObjectInfo // Changed type

	dirPath := filepath.Join(fs.basePath, path)
	fileInfos, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fileInfos {
		info, err := fileInfo.Info()
		if err != nil {
			return nil, err
		}

		// Construct storage.ObjectInfo
		files = append(files, &storage.ObjectInfo{ // Changed type
			Path:    filepath.Join(path, info.Name()),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			Metadata: map[string]interface{}{
				"is_dir": info.IsDir(), // Add IsDir to metadata
			},
		})
	}

	return files, nil
}

// Stat returns ObjectInfo for the given path.
func (fs *LocalStorage) Stat(path string) (*storage.ObjectInfo, error) { // Changed return type
	info, err := os.Stat(filepath.Join(fs.basePath, path))
	if err != nil {
		return nil, err // Changed zero value to nil
	}

	return &storage.ObjectInfo{ // Changed type
		Path:    path,
		Size:    info.Size(),
		ModTime: info.ModTime(),
		Metadata: map[string]interface{}{
			"is_dir": info.IsDir(), // Add IsDir to metadata
		},
	}, nil
}

// Mkdir creates a new directory.
func (fs *LocalStorage) Mkdir(path string) error {
	return os.MkdirAll(filepath.Join(fs.basePath, path), os.ModePerm)
}

// Delete removes a file or directory.
func (fs *LocalStorage) Delete(path string) error {
	return os.RemoveAll(filepath.Join(fs.basePath, path))
}

// Rename renames a file or directory.
func (fs *LocalStorage) Rename(oldPath, newPath string) error {
	oldFullPath := filepath.Join(fs.basePath, oldPath)
	newFullPath := filepath.Join(fs.basePath, newPath)
	return os.Rename(oldFullPath, newFullPath)
}
