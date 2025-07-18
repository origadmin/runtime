/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package blob

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"os"
	"path/filepath"
)

type blobStorage struct {
	Path string
	Hash func() hash.Hash
}

func (m blobStorage) Write(data []byte) (string, error) {
	encodeHash := m.getHash(data)
	path := hashPath(m.Path, encodeHash)

	// Ensure the directory for the blob exists.
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	if err := atomicWrite(path, data); err != nil {
		return "", err
	}
	return encodeHash, nil
}

func (m blobStorage) Read(hash string) ([]byte, error) {
	path := hashPath(m.Path, hash)
	return os.ReadFile(path)
}

func (m blobStorage) Delete(hash string) error {
	path := hashPath(m.Path, hash)
	return os.Remove(path)
}

func (m blobStorage) Exists(hash string) (bool, error) {
	path := hashPath(m.Path, hash)
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return !stat.IsDir(), nil
}

func hashPath(path, hash string) string {
	return path + "/" + hash[:2] + "/" + hash[2:4] + "/" + hash
}

func atomicWrite(path string, content []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, content, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (m blobStorage) getHash(content []byte) string {
	h := m.Hash()
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))
}

func New(path string) *blobStorage {
	return &blobStorage{
		Path: path,
		Hash: sha256.New,
	}
}
