/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package upload is the http multipart upload package
package upload

import (
	"io/fs"
	"mime/multipart"
	"time"
)

type FileInfo interface {
	fs.FileInfo
	ContentType() string
}

type multipartInfo struct {
	modTime     time.Time
	isDir       bool
	sys         any
	name        string
	size        int64
	contentType string
}

func (f *multipartInfo) ContentType() string {
	return f.contentType
}

func (f *multipartInfo) Name() string {
	return f.name
}

func (f *multipartInfo) Size() int64 {
	return f.size
}

func (f *multipartInfo) Sys() any {
	return f.sys
}

func (f *multipartInfo) Mode() fs.FileMode {
	return 0o644
}

func (f *multipartInfo) ModTime() time.Time {
	return f.modTime
}

func (f *multipartInfo) IsDir() bool {
	return f.isDir
}

func ParseMultipart(header *multipart.FileHeader) FileInfo {
	mod, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", header.Header.Get("Last-Modified"))
	return &multipartInfo{
		name:        header.Filename,
		size:        header.Size,
		sys:         header.Header,
		contentType: header.Header.Get("Content-Type"),
		modTime:     mod,
		isDir:       false,
	}
}

var _ fs.FileInfo = (*multipartInfo)(nil)
var _ FileInfo = (*multipartInfo)(nil)
