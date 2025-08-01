package main

import (
	"io"
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

// List returns a slice of FileInfo for the given directory path.
func (fs *LocalStorage) List(path string) ([]storage.FileInfo, error) {
	var files []storage.FileInfo

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

		files = append(files, storage.FileInfo{
			Name:    info.Name(),
			Path:    filepath.Join(path, info.Name()),
			IsDir:   info.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}

	return files, nil
}

// Read returns a reader for the given file path.
func (fs *LocalStorage) Read(path string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(fs.basePath, path))
}

// Write writes data to the given file path.
func (fs *LocalStorage) Write(path string, data io.Reader, size int64) error {
	filePath := filepath.Join(fs.basePath, path)
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, data)
	return err
}

// Stat returns FileInfo for the given path.
func (fs *LocalStorage) Stat(path string) (storage.FileInfo, error) {
	info, err := os.Stat(filepath.Join(fs.basePath, path))
	if err != nil {
		return storage.FileInfo{}, err
	}

	return storage.FileInfo{
		Name:    info.Name(),
		Path:    path,
		IsDir:   info.IsDir(),
		Size:    info.Size(),
		ModTime: info.ModTime(),
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
