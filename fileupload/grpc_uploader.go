package fileupload

import (
	"context"
	"errors"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/origadmin/runtime"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	fileuploadv1 "github.com/origadmin/runtime/gen/go/fileupload/v1"
	"github.com/origadmin/toolkits/fileupload"
)

type grpcUploader struct {
	builder Builder
	client  fileuploadv1.FileUploadServiceClient
	stream  grpc.ClientStreamingClient[fileuploadv1.UploadRequest, fileuploadv1.UploadResponse]
	header  *fileuploadv1.FileHeader
}

func (u *grpcUploader) Resume(ctx context.Context, offset int64) error {
	//TODO implement me
	panic("implement me")
}

func (u *grpcUploader) SetFileHeader(ctx context.Context, header fileupload.FileHeader) error {
	if header == nil {
		return errors.New("invalid file header")
	}
	if v, ok := header.(*fileuploadv1.FileHeader); ok {
		u.header = v
		return nil
	}

	//send fileHeader
	u.header = &fileuploadv1.FileHeader{
		Filename:    header.GetFilename(),
		Size:        header.GetSize(),
		ContentType: header.GetContentType(),
		ModTime:     header.GetModTime(),
		Header:      header.GetHeader(),
	}
	return nil
}

func (u *grpcUploader) UploadFile(ctx context.Context, rd io.Reader) error {
	header, err := proto.Marshal(u.header)
	if err != nil {
		return err
	}
	if err := u.stream.Send(&fileuploadv1.UploadRequest{
		IsHeader: true,
		Data:     header,
	}); err != nil {
		return err
	}
	buf := u.builder.NewBuffer()
	defer u.builder.Free(buf)
	for {
		n, err := rd.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := u.stream.Send(&fileuploadv1.UploadRequest{
			IsHeader: false,
			Data:     buf[:n],
		}); err != nil {
			return err
		}
	}

	return nil
}

func (u *grpcUploader) Finalize(ctx context.Context) (fileupload.UploadResponse, error) {
	resp, err := u.stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewGRPCUploader(ctx context.Context, service *configv1.Service) (fileupload.Uploader, error) {
	clientService, err := runtime.NewGRPCServiceClient(ctx, service)
	if err != nil {
		return nil, err
	}

	client := fileuploadv1.NewFileUploadServiceClient(clientService)
	stream, err := client.Upload(ctx)
	if err != nil {
		return nil, err
	}
	return &grpcUploader{
		client: client,
		stream: stream,
	}, nil

}
