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
	taskId  string
}

func (u *grpcUploader) SetFileHeader(ctx context.Context, header fileupload.FileHeader) error {
	protoHeader := Header2GRPCHeader(header)
	req := &fileuploadv1.CreateUploadTaskRequest{
		FileHeader: protoHeader,
	}
	resp, err := u.client.CreateUploadTask(ctx, req)
	if err != nil {
		return err
	}
	u.taskId = resp.GetTaskId()
	return nil
}

func (u *grpcUploader) UploadFile(ctx context.Context, rd io.Reader) error {
	stream, err := u.client.UploadChunk(ctx)
	if err != nil {
		return err
	}
	defer stream.CloseSend()

	buf := u.builder.NewBuffer()
	chunkNumber := int32(0)
	for {
		n, err := rd.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&fileuploadv1.UploadChunkRequest{
			TaskId:      u.taskId,
			ChunkNumber: chunkNumber,
			Data:        buf[:n],
		}); err != nil {
			return err
		}
		chunkNumber++
	}
	return nil
}

func (u *grpcUploader) Finalize(ctx context.Context) (fileupload.UploadResponse, error) {
	req := &fileuploadv1.FinalizeUploadRequest{
		TaskId: u.taskId,
	}
	resp, err := u.client.FinalizeUpload(ctx, req)
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
	return &grpcUploader{
		builder: Builder{},
		client:  client,
	}, nil
}
