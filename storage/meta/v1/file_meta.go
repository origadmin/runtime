/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package metav1 implements the functions, types, and interfaces for the module.
package metav1

type FileMetaV1 struct {
	Version  int32  `msgpack:"v"` // Schema version, e.g., 1
	Hash     string `msgpack:"h"` // Content hash
	Size     int64  `msgpack:"s"` // File size
	MimeType string `msgpack:"m"` // MIME type
	ModTime  int64  `msgpack:"t"` // Modify time

	// ... other fields
	BlockSize   int32    `msgpack:"bs"` // New field
	BlockHashes []string `msgpack:"bh"` // New field
}

func (f FileMetaV1) CurrentVersion() int32 {
	return 1
}
