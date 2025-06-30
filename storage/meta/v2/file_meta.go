/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package metav2 implements the functions, types, and interfaces for the module.
package metav2

const Version = 2

type FileMetaV2 struct {
	Version  int32  `msgpack:"v"` // Schema version, e.g., 1
	Name     string `msgpack:"n"` // File name
	Hash     string `msgpack:"h"` // Content hash
	Size     int64  `msgpack:"s"` // File size
	MimeType string `msgpack:"m"` // MIME type
	ModTime  int64  `msgpack:"t"` // Modify time

	// if version > 1, then we have the following fields:
	BlockSize   int32             `msgpack:"bs"` // New field
	BlockHashes []string          `msgpack:"bh"` // New field
	Extra       map[string]string `msgpack:"e"`  // extra data
}

func (f FileMetaV2) CurrentVersion() int32 {
	return Version
}
