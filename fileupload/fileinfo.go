/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package fileupload implements the functions, types, and interfaces for the module.
package fileupload

import (
	"io/fs"
)

type FileInfo interface {
	fs.FileInfo
	ContentType() string
}
