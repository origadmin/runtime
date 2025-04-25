package fileupload

import (
	"context"
	"mime"
	"os"
	"path/filepath"
)

type FolderUploader struct {
	builder  Builder
	basePath string
}

func NewFolderUploader(builder Builder, basePath string) *FolderUploader {
	return &FolderUploader{
		builder:  builder,
		basePath: basePath,
	}
}

func (f *FolderUploader) UploadFolder(ctx context.Context, path string) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		uploader := f.builder.NewUploader(ctx)
		if err != nil {
			return err
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		relPath, _ := filepath.Rel(f.basePath, filePath)
		header := &httpFileHeader{
			Filename:    relPath,
			Size:        uint32(info.Size()),
			ModTime:     uint32(info.ModTime().Unix()),
			ContentType: mime.TypeByExtension(filepath.Ext(filePath)),
			IsDir:       info.IsDir(),
		}

		if err := uploader.SetFileHeader(ctx, header); err != nil {
			return err
		}

		if err := uploader.UploadFile(ctx, file); err != nil {
			return err
		}

		_, err = uploader.Finalize(ctx)
		return err
	})
}
