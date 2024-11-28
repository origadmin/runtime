/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package upload

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/io"
)

const (
	ErrInvalidFile    = errors.String("invalid file")
	ErrTargetIsNotDir = errors.String("target is not a directory")
	ErrSizeNotMatch   = io.ErrSizeNotMatch
)

type Uploader interface {
	UploadFromMultipart(req *http.Request, name string) (fs.FileInfo, error)
}

type uploader struct {
	keys []string
}

// FileUpload upload file
func FileUpload(req *http.Request, path, name string) (fs.FileInfo, error) {
	// replace path to os path
	path = filepath.FromSlash(path)
	abspath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	src, header, err := req.FormFile("file")
	if err != nil {
		return nil, err
	}

	if header.Size == 0 || header.Filename == "" {
		return nil, ErrInvalidFile
	}

	switch name == "" {
	case true:
		name = header.Filename
	case false:
		name += filepath.Ext(header.Filename)
	}

	dst, err := os.OpenFile(filepath.Join(abspath, name), os.O_RDONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, err
	}
	defer dst.Close()
	copied, err := io.CopyContext(req.Context(), dst, src)
	if err != nil {
		return nil, err
	}
	if copied != header.Size {
		return nil, ErrSizeNotMatch
	}
	return ParseMultipart(header), nil
}
