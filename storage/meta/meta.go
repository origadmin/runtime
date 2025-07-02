/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	p "path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/origadmin/runtime/interfaces/storage/meta"
	metav1 "github.com/origadmin/runtime/storage/meta/v1"
	metav2 "github.com/origadmin/runtime/storage/meta/v2"
)

const (
	// rootFileName defines the name of the file that stores the hash of the root directory index.
	rootFileName = ".root"
)

// Meta is the main entry point for managing file system metadata.
// It orchestrates interactions between directory indexes, file metadata, and blob storage.
// It is safe for concurrent use.
type Meta struct {
	path     string
	rootFile string
	storage  meta.BlobStorage

	// root holds the current root DirectoryIndex.
	root *DirectoryIndex
	// rootLock ensures that operations that modify the tree are serialized.
	rootLock sync.Mutex

	// dirCache caches recently accessed DirectoryIndex objects.
	// The key is the hash of the directory's index blob.
	dirCache  map[string]*DirectoryIndex
	cacheLock sync.RWMutex
}

// New creates and initializes a new Meta manager.
// It sets up the blob storage and ensures the root directory exists.
func New(path string) (*Meta, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}
	blobPath := filepath.Join(path, "blobs")
	if err := os.MkdirAll(blobPath, 0755); err != nil {
		return nil, err
	}
	rootPath := filepath.Join(path, "index")
	if err := os.MkdirAll(rootPath, 0755); err != nil {
		return nil, err
	}

	storage := NewBlobStorage(blobPath)
	m := &Meta{
		path:     path,
		rootFile: filepath.Join(rootPath, rootFileName),
		storage:  storage,
		dirCache: make(map[string]*DirectoryIndex),
	}

	if err := m.initRoot(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Meta) initRoot() error {
	content, err := os.ReadFile(m.rootFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		rootIndex := &DirectoryIndex{
			Version: 1,
			Path:    "/",
			Entries: []FileIndexEntry{},
		}
		marshal, err := msgpack.Marshal(rootIndex)
		if err != nil {
			return err
		}

		if err := atomicWrite(m.rootFile, marshal); err != nil {
			return err
		}
		m.root = rootIndex
		return nil
	}
	if err := msgpack.Unmarshal(content, &m.root); err != nil {
		return err
	}
	return nil
}

func (m *Meta) getDirectoryIndexByHash(hash string) (*DirectoryIndex, error) {
	m.cacheLock.RLock()
	if idx, found := m.dirCache[hash]; found {
		m.cacheLock.RUnlock()
		return idx, nil
	}
	m.cacheLock.RUnlock()

	data, err := m.storage.Retrieve(hash)
	if err != nil {
		return nil, err
	}

	var idx DirectoryIndex
	if err := msgpack.Unmarshal(data, &idx); err != nil {
		return nil, err
	}

	cachedIdx := idx

	m.cacheLock.Lock()
	m.dirCache[hash] = &cachedIdx
	m.cacheLock.Unlock()

	return &cachedIdx, nil
}

func (m *Meta) storeIndex(idx *DirectoryIndex) (string, error) {
	data, err := msgpack.Marshal(idx)
	if err != nil {
		return "", err
	}
	hash, err := m.storage.Store(data)
	if err != nil {
		return "", err
	}

	m.cacheLock.Lock()
	m.dirCache[hash] = idx
	m.cacheLock.Unlock()
	return hash, nil
}

func (d *DirectoryIndex) findEntry(name string) (*FileIndexEntry, bool) {
	for i, entry := range d.Entries {
		if entry.EntryName == name {
			return &d.Entries[i], true
		}
	}
	return nil, false
}

// findEntryByPath traverses the directory tree to find the entry for a given path.
func (m *Meta) findEntryByPath(path string) (*FileIndexEntry, error) {
	cleanPath := p.Clean(path)

	slog.Info("find entry by path", "path", path, "cleanPath", cleanPath)
	if !strings.HasPrefix(cleanPath, "/") {
		return nil, &fs.PathError{Op: "find", Path: path, Err: errors.New("path must be absolute and start with /")}
	}

	if cleanPath == "/" {
		// Return a virtual entry for the root directory itself.
		return &FileIndexEntry{
			EntryName:   "/",
			Hash:        "",
			FileMode:    fs.ModeDir,
			IsDirectory: true,
		}, nil
	}

	parts := strings.Split(strings.TrimPrefix(cleanPath, "/"), "/")
	current := m.root

	for i, part := range parts {
		entry, found := current.findEntry(part)
		if !found {
			return nil, &fs.PathError{Op: "find", Path: path, Err: fs.ErrNotExist}
		}

		// If this is the last part of the path, we found our entry.
		if i == len(parts)-1 {
			return entry, nil
		}

		// If it's not the last part, it must be a directory to continue.
		if !entry.IsDirectory {
			return nil, &fs.PathError{Op: "find", Path: path, Err: errors.New("not a directory")}
		}
		var err error
		current, err = m.getDirectoryIndexByHash(entry.Hash)
		if err != nil {
			return nil, err
		}
	}

	// This part should not be reached.
	return nil, &fs.PathError{Op: "find", Path: path, Err: fs.ErrNotExist}
}

func (m *Meta) Mkdir(path string, perm fs.FileMode) error {
	m.rootLock.Lock()
	defer m.rootLock.Unlock()

	cleanPath := p.Clean(path)
	if !strings.HasPrefix(cleanPath, "/") {
		return &fs.PathError{Op: "mkdir", Path: path, Err: errors.New("path must be absolute and start with /")}
	}

	if cleanPath == "/" {
		return &fs.PathError{Op: "mkdir", Path: path, Err: errors.New("cannot create root directory")}
	}

	parentPath := p.Dir(cleanPath)
	name := p.Base(cleanPath)

	newDirIndex := &DirectoryIndex{
		Version: 1,
		Path:    cleanPath,
		Entries: []FileIndexEntry{},
	}
	newDirHash, err := m.storeIndex(newDirIndex)
	if err != nil {
		return err
	}

	newEntry := FileIndexEntry{
		EntryName:   name,
		Hash:        newDirHash,
		FileMode:    perm | fs.ModeDir,
		IsDirectory: true,
	}

	_, err = m.updateTree(parentPath, func(dir *DirectoryIndex) (*FileIndexEntry, error) {
		if _, found := dir.findEntry(name); found {
			return nil, fs.ErrExist
		}
		return &newEntry, nil
	})

	if err != nil {
		return err
	}

	return m.writeRootFile()
}

// ReadDir reads the directory named by path and returns a list of directory entries.
func (m *Meta) ReadDir(path string) ([]fs.DirEntry, error) {
	entry, err := m.findEntryByPath(path)
	if err != nil {
		return nil, err
	}

	if !entry.IsDirectory {
		return nil, &fs.PathError{Op: "readdir", Path: path, Err: errors.New("not a directory")}
	}

	// Special case for the root directory
	if path == "/" {
		dirEntries := make([]fs.DirEntry, len(m.root.Entries))
		for i, e := range m.root.Entries {
			dirEntries[i] = e
		}
		return dirEntries, nil
	}

	dirIndex, err := m.getDirectoryIndexByHash(entry.Hash)
	if err != nil {
		return nil, err
	}

	// Convert []FileIndexEntry to []fs.DirEntry
	dirEntries := make([]fs.DirEntry, len(dirIndex.Entries))
	for i, e := range dirIndex.Entries {
		dirEntries[i] = e
	}

	return dirEntries, nil
}

// WriteFile writes the content from r to the file at the given path.
// It creates a new file or overwrites an existing one.
func (m *Meta) WriteFile(path string, r io.Reader, perm fs.FileMode) error {
	m.rootLock.Lock()
	defer m.rootLock.Unlock()

	cleanPath := p.Clean(path)
	if !strings.HasPrefix(cleanPath, "/") {
		return &fs.PathError{Op: "write", Path: path, Err: errors.New("path must be absolute and start with /")}
	}

	if cleanPath == "/" {
		return &fs.PathError{Op: "write", Path: path, Err: errors.New("cannot write to root directory")}
	}

	parentPath := p.Dir(cleanPath)
	name := p.Base(cleanPath)

	// 1. Chunk and store data in blob storage
	blockHashes, totalSize, err := chunkData(r, m.storage.Store)
	if err != nil {
		return err
	}

	// 2. Create FileMetaV2
	fileMetaV2 := &metav2.FileMetaV2{
		Version:     metav2.Version,
		Size:        totalSize,
		ModTime:     time.Now().Unix(),
		MimeType:    "application/octet-stream", // Placeholder, can be improved with content sniffing
		RefCount:    1,                          // Initial ref count
		BlockSize:   DefaultBlockSize,
		BlockHashes: blockHashes,
	}

	// 3. Marshal and store FileMetaV2 in blob storage
	fileMetaBytes, err := msgpack.Marshal(fileMetaV2)
	if err != nil {
		return err
	}
	fileMetaHash, err := m.storage.Store(fileMetaBytes)
	if err != nil {
		return err
	}

	// 4. Create FileIndexEntry
	newEntry := FileIndexEntry{
		EntryName:   name,
		Hash:        fileMetaHash,
		FileMode:    perm,
		IsDirectory: false,
	}

	// 5. Update the directory tree
	_, err = m.updateTree(parentPath, func(dir *DirectoryIndex) (*FileIndexEntry, error) {
		// Check if file already exists, if so, replace it.
		for i, entry := range dir.Entries {
			if entry.EntryName == name {
				if entry.IsDirectory {
					return nil, &fs.PathError{Op: "write", Path: path, Err: errors.New("is a directory")}
				}
				// Replace existing file entry
				dir.Entries[i] = newEntry
				return &newEntry, nil
			}
		}
		// Add new file entry
		dir.Entries = append(dir.Entries, newEntry)
		return &newEntry, nil
	})

	if err != nil {
		return err
	}

	return m.writeRootFile()
}

// A temporary struct to unmarshal only the version to determine the actual FileMeta type.
type versionOnly struct {
	Version int32 `msgpack:"v"`
}

// Open opens the named file for reading.
func (m *Meta) Open(path string) (fs.File, error) {
	entry, err := m.findEntryByPath(path)
	if err != nil {
		return nil, err
	}

	if entry.IsDirectory {
		return nil, &fs.PathError{Op: "open", Path: path, Err: errors.New("is a directory")}
	}

	fileMetaBytes, err := m.storage.Retrieve(entry.Hash)
	if err != nil {
		return nil, err
	}

	var vOnly versionOnly
	if err := msgpack.Unmarshal(fileMetaBytes, &vOnly); err != nil {
		return nil, err
	}

	var fileMeta any
	switch vOnly.Version {
	case metav1.Version:
		var fm metav1.FileMeta
		if err := msgpack.Unmarshal(fileMetaBytes, &fm); err != nil {
			return nil, err
		}
		fileMeta = &fm
	case metav2.Version:
		var fm metav2.FileMeta
		if err := msgpack.Unmarshal(fileMetaBytes, &fm); err != nil {
			return nil, err
		}
		fileMeta = &fm
	default:
		return nil, fmt.Errorf("unsupported file meta version: %d", vOnly.Version)
	}

	return OpenMetaFile(m.storage, fileMeta)
}

// Stat returns a FileInfo describing the file or directory named by path.
func (m *Meta) Stat(path string) (fs.FileInfo, error) {
	entry, err := m.findEntryByPath(path)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (m *Meta) updateTree(path string, modifier func(*DirectoryIndex) (*FileIndexEntry, error)) (string, error) {
	m.rootLock.Lock()
	defer m.rootLock.Unlock()

	cleanPath := p.Clean(path)
	if !strings.HasPrefix(cleanPath, "/") {
		return "", errors.New("path must be absolute")
	}

	parentDir, err := m.findDirectoryByPath(cleanPath)
	if err != nil {
		return "", err
	}

	_, err = modifier(parentDir)
	if err != nil {
		return "", err
	}

	return m.storeDirectoryIndex(parentDir)
}

func (m *Meta) findDirectoryByPath(path string) (*DirectoryIndex, error) {
	cleanPath := p.Clean(path)
	if !strings.HasPrefix(cleanPath, "/") {
		return nil, errors.New("path must be absolute")
	}

	if cleanPath == "/" {
		return m.root, nil
	}

	parts := strings.Split(strings.TrimPrefix(cleanPath, "/"), "/")
	dir := m.root

	for _, part := range parts {
		entry, found := dir.findEntry(part)
		if !found || !entry.IsDirectory {
			return nil, fs.ErrNotExist
		}

		idx, err := m.getDirectoryIndexByHash(entry.Hash)
		if err != nil {
			return nil, err
		}
		dir = idx
	}

	return dir, nil
}

func (m *Meta) storeDirectoryIndex(idx *DirectoryIndex) (string, error) {
	data, err := msgpack.Marshal(idx)
	if err != nil {
		return "", err
	}

	hash, err := m.storage.Store(data)
	if err != nil {
		return "", err
	}

	m.cacheLock.Lock()
	m.dirCache[hash] = idx
	m.cacheLock.Unlock()

	return hash, nil
}

func (m *Meta) writeRootFile() error {
	data, err := msgpack.Marshal(m.root)
	if err != nil {
		return err
	}

	return atomicWrite(m.rootFile, data)
}
