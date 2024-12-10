/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package stream

import (
	"errors"

	"github.com/origadmin/toolkits/fileupload"
	"github.com/origadmin/toolkits/io"
	"google.golang.org/protobuf/proto"

	"github.com/origadmin/runtime/context"
	fileuploadv1 "github.com/origadmin/runtime/gen/go/fileupload/v1"
)

type grpcUploader struct {
	//builder *grpcBuilder
	stream fileuploadv1.FileUploadService_UploadClient
	header *fileuploadv1.FileHeader
}

func (u *grpcUploader) SetFileHeader(ctx context.Context, header fileupload.FileHeader) error {
	if protoHeader, ok := header.(*fileuploadv1.FileHeader); ok {
		u.header = protoHeader
	} else {
		u.header = &fileuploadv1.FileHeader{
			Filename:      header.GetFilename(),
			Size:          header.GetSize(),
			ModTimeString: header.GetModTimeString(),
			ModTime:       header.GetModTime(),
			ContentType:   header.GetContentType(),
			Header:        header.GetHeader(),
			IsDir:         header.GetIsDir(),
		}
	}

	// Send header chunk
	headerData, err := proto.Marshal(u.header)
	if err != nil {
		return err
	}

	return u.stream.Send(&fileuploadv1.UploadRequest{
		IsHeader: true,
		Data:     headerData,
	})
}

func (u *grpcUploader) UploadFile(ctx context.Context, rd io.Reader) error {
	buf := u.builder.bufPool.Get().([]byte)
	defer u.builder.bufPool.Put(buf)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, err := rd.Read(buf)
			if n > 0 {
				chunk := &fileuploadv1.UploadRequest{
					Data: buf[:n],
				}
				if err := u.stream.Send(chunk); err != nil {
					return err
				}
			}
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
}
