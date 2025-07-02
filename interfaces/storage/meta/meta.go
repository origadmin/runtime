/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"io/fs"
)

type BlobStorage interface {
	Store(content []byte) (string, error)
	Retrieve(hash string) ([]byte, error)
}

type FileMeta struct {
	Info *FileMetaInfo `json:"info" msgpack:"i"`
	Data any           `json:"data" msgpack:"d"`
}

type FileMetaVersion struct {
	Version int32 `msgpack:"v"`
}

func (f FileMetaVersion) CurrentVersion() int32 {
	return f.Version
}

type FileMetaInfo struct {
	Name     string      `json:"name" msgpack:"n"`      // File name
	Hash     string      `json:"hash" msgpack:"h"`      // File hash
	FileMode fs.FileMode `json:"file_mode" msgpack:"m"` // File mode
	Sys      any         `json:"sys" msgpack:"s"`       // System-specific data
}

type FileMetaData[T any] struct {
	Version int32         `json:"version" msgpack:"v"`
	Info    *FileMetaInfo `json:"info" msgpack:"i"`
	Data    *T            `json:"data" msgpack:"d"`
}
