/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package stream is the http multipart upload package
package stream

import (
	"io/fs"
	"mime/multipart"
	"time"

	fileuploadv1 "github.com/origadmin/runtime/gen/go/fileupload/v1"
)

type FileInfo interface {
	fs.FileInfo
	ContentType() string
}

type multipartFileInfo struct {
	header *fileuploadv1.FileHeader
}

func (f *multipartFileInfo) ContentType() string {
	return f.header.ContentType
}

func (f *multipartFileInfo) Name() string {
	return f.header.Filename
}

func (f *multipartFileInfo) Size() int64 {
	return int64(f.header.Size)
}

func (f *multipartFileInfo) Sys() any {
	return nil
}

func (f *multipartFileInfo) Mode() fs.FileMode {
	return 0o644
}

func (f *multipartFileInfo) ModTime() time.Time {
	return time.Unix(int64(f.header.ModTime), 0)
}

func (f *multipartFileInfo) IsDir() bool {
	return f.header.IsDir
}

func ParseMultipart(header *multipart.FileHeader) FileInfo {
	headerStr := header.Header.Get("Last-Modified")
	mod, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", headerStr)
	requestHeader := make(map[string]string, len(header.Header))
	for k, v := range header.Header {
		requestHeader[k] = v[0]
	}
	fileHeader := &fileuploadv1.FileHeader{
		ContentType:   header.Header.Get("Content-Type"),
		Filename:      header.Filename,
		Header:        requestHeader,
		ModTime:       uint32(mod.Unix()),
		ModTimeString: headerStr,
		Size:          uint32(header.Size),
		IsDir:         false,
	}
	return &multipartFileInfo{
		header: fileHeader,
	}
}

var _ fs.FileInfo = (*multipartFileInfo)(nil)
var _ FileInfo = (*multipartFileInfo)(nil)
