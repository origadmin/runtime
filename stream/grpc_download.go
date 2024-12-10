/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package stream

import (
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/fileupload"
	"github.com/origadmin/toolkits/io"
	"google.golang.org/protobuf/proto"

	"github.com/origadmin/runtime/context"
	fileuploadv1 "github.com/origadmin/runtime/gen/go/fileupload/v1"
)

const (
	ErrInvalidFile    = errors.String("invalid file")
	ErrTargetIsNotDir = errors.String("target is not a directory")
	ErrSizeNotMatch   = io.ErrSizeNotMatch
)

type grpcDownloader struct {
	//builder *grpcBuilder
	stream fileuploadv1.FileUploadService_UploadServer
	header *fileuploadv1.FileHeader
	resp   UploadResponse
}

func (d *grpcDownloader) GetFileHeader(ctx context.Context) (fileupload.FileHeader, error) {
	chunk, err := d.stream.Recv()
	if err != nil {
		return nil, err
	}

	if !chunk.IsHeader {
		return nil, ErrInvalidRequest
	}

	var header FileHeader
	if err := proto.Unmarshal(chunk.GetData(), &header); err != nil {
		return nil, err
	}

	d.header = &header
	return d.header, nil
}

func (d *grpcDownloader) DownloadFile(ctx context.Context) (io.Reader, error) {
	if d.header == nil {
		return nil, ErrNoFile
	}

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		for {
			chunk, err := d.stream.Recv()
			if err == io.EOF {
				if resp, recvErr := d.stream.CloseAndRecv(); recvErr == nil {
					d.resp = &grpcFileResponse{
						success:    resp.Success,
						hash:       resp.Hash,
						path:       resp.Path,
						size:       resp.Size,
						failReason: resp.FailReason,
					}
				}
				return
			}
			if err != nil {
				pw.CloseWithError(err)
				return
			}

			if chunk.IsHeader {
				continue // Skip header chunks
			}

			if _, err := pw.Write(chunk.Content); err != nil {
				pw.CloseWithError(err)
				return
			}
		}
	}()

	return pr, nil
}

func (d *grpcDownloader) Finalize(ctx context.Context, resp UploadResponse) error {
	if d.resp == nil {
		return ErrUploadFailed
	}
	if ptr, ok := resp.(*grpcFileResponse); ok {
		*ptr = *d.resp.(*grpcFileResponse)
	}
	return nil
}

type grpcFileResponse struct {
	success    bool
	hash       string
	path       string
	size       uint32
	failReason string
}

func (r *grpcFileResponse) GetSuccess() bool      { return r.success }
func (r *grpcFileResponse) GetHash() string       { return r.hash }
func (r *grpcFileResponse) GetPath() string       { return r.path }
func (r *grpcFileResponse) GetSize() uint32       { return r.size }
func (r *grpcFileResponse) GetFailReason() string { return r.failReason }
