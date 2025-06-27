/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"os"

	"github.com/origadmin/runtime/interfaces/storage/meta"
)

type blobStorage struct {
	Path string
	Hash func() hash.Hash
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

func (m blobStorage) Store(content []byte) (string, error) {
	encodeHash := m.getHash(content)
	path := hashPath(m.Path, encodeHash)

	if err := atomicWrite(path, content); err != nil {
		return "", err
	}
	return encodeHash, nil
}

func (m blobStorage) Retrieve(hash string) ([]byte, error) {
	path := hashPath(m.Path, hash)
	return os.ReadFile(path)
}

func NewBlobStorage(path string) meta.BlobStorage {
	return &blobStorage{
		Path: path,
		Hash: sha256.New,
	}
}
