package fileupload

import (
	"context"
	"io"

	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"

	"github.com/origadmin/runtime"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	fileuploadv1 "github.com/origadmin/runtime/gen/go/fileupload/v1"
	"github.com/origadmin/toolkits/fileupload"
)

type grpcUploader struct {
	builder Builder
	client  fileuploadv1.FileUploadServiceClient
	stream  grpc.ClientStreamingClient[fileuploadv1.UploadRequest, fileuploadv1.UploadResponse]
	header  []byte
	buf     []byte
}

func (u *grpcUploader) Resume(ctx context.Context, offset int64) error {
	//TODO implement me
	panic("implement me")
}

func (u *grpcUploader) SetFileHeader(ctx context.Context, header fileupload.FileHeader) error {
	//u.fileHeader = Header2GRPCHeader(fileHeader)

	//// 初始化stream
	//stream, err := u.client.Upload(ctx)
	//if err != nil {
	//	return err
	//}
	//u.stream = stream

	//send fileHeader
	var err error
	u.header, err = proto.Marshal(&fileuploadv1.FileHeader{
		Filename:    header.GetFilename(),
		Size:        header.GetSize(),
		ContentType: header.GetContentType(),
		ModTime:     header.GetModTime(),
		Header:      header.GetHeader(),
	})

	if err != nil {
		return err
	}
	return nil
}

func (u *grpcUploader) UploadFile(ctx context.Context, rd io.Reader) error {
	if len(u.header) > 0 {
		if err := u.stream.Send(&fileuploadv1.UploadRequest{
			IsHeader: true,
			Data:     u.header,
		}); err != nil {
			return err
		}
	}
	if u.buf == nil {
		u.buf = u.builder.NewBuffer()
	}

	for {
		n, err := rd.Read(u.buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := u.stream.Send(&fileuploadv1.UploadRequest{
			IsHeader: false,
			Data:     u.buf[:n],
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
