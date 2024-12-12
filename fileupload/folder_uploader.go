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

		// 跳过目录本身
		if info.IsDir() {
			return nil
		}

		// 创建新的上传器
		uploader := f.builder.NewUploader(ctx)
		if err != nil {
			return err
		}

		// 打开文件
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// 设置文件头信息
		relPath, _ := filepath.Rel(f.basePath, filePath)
		header := &httpFileHeader{
			Filename:    relPath,
			Size:        uint32(info.Size()),
			ModTime:     uint32(info.ModTime().Unix()),
			ContentType: mime.TypeByExtension(filepath.Ext(filePath)),
			IsDir:       info.IsDir(),
		}

		// 上传文件
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
