/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package fileupload implements the functions, types, and interfaces for the module.
package fileupload

import (
	"context"
	"io"
	"net/http"
)

const (
	ModTimeFormat = "Mon, 02 Jan 2006 15:04:05 MST"
)

// FileHeader represents a file header with metadata.
type FileHeader interface {
	// GetFilename returns the name of the file.
	GetFilename() string

	// GetSize returns the size of the file in bytes.
	GetSize() uint32

	// GetModTime returns the last modified time of the file as a string.
	GetModTime() uint32

	// GetModTimeString returns the last modified time of the file as a Unix timestamp.
	GetModTimeString() string

	// GetContentType returns the MIME type of the file.
	GetContentType() string

	// GetHeader returns a map of additional file headers.
	GetHeader() map[string]string

	// GetIsDir returns true if the file is a directory.
	GetIsDir() bool
}

// UploadResponse represents a response to a file upload request.
type UploadResponse interface {
	// GetSuccess indicates whether the file upload was successful.
	GetSuccess() bool

	// GetHash returns the hash of the uploaded file.
	GetHash() string

	// GetPath returns the path where the file was uploaded.
	GetPath() string

	// GetSize returns the size of the uploaded file in bytes.
	GetSize() uint32

	// GetFailReason returns the failure reason of the file upload.
	GetFailReason() string
}

type Uploader interface {
	// SetFileHeader sets the file header after the file has been uploaded.
	SetFileHeader(ctx context.Context, header FileHeader) error
	// UploadFile uploads the file first and then sets the header.
	UploadFile(ctx context.Context, rd io.Reader) error
	// Resume resumes upload from specified offset
	Resume(ctx context.Context, offset int64) error
	// Finalize finalizes the upload process.
	Finalize(ctx context.Context) (UploadResponse, error)
}

type Receiver interface {
	// GetFileHeader retrieves the file header after the file has been received.
	GetFileHeader(ctx context.Context) (FileHeader, error)
	// ReceiveFile receives the file first and then retrieves the header.
	ReceiveFile(ctx context.Context) (io.ReadCloser, error)
	// GetOffset returns current offset for resume support
	GetOffset(ctx context.Context) (int64, error)
	// Finalize finalizes the receipt process.
	Finalize(ctx context.Context, resp UploadResponse) error
}

type BridgeUploader interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
