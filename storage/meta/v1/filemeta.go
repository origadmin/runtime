/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package metav1 implements the functions, types, and interfaces for the module.
package metav1

import (
	"time"
)

const Version = 1

type FileMetaV1 struct {
	Version    int32  `msgpack:"v"` // File meta version
	FileSize   int64  `msgpack:"s"` // File size
	MimeType   string `msgpack:"m"` // MIME type
	ModifyTime int64  `msgpack:"t"` // Modify time
	RefCount   int32  `msgpack:"r"` // Reference count
}

func (f FileMetaV1) CurrentVersion() int32 {
	return Version
}

func (f FileMetaV1) Size() int64 {
	return f.FileSize
}

func (f FileMetaV1) ModTime() time.Time {
	return time.Unix(f.ModifyTime, 0)
}
