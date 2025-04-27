package fileupload

import (
	"context"
	"errors"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/toolkits/fileupload"

	fileuploadv1 "github.com/origadmin/runtime/gen/go/fileupload/v1"
)

type grpcReceiver struct {
	//header *fileuploadv1.FileHeader
	//stream grpc.ServerStream
	server fileuploadv1.FileUploadServiceServer
	pr     *io.PipeReader
	pw     *io.PipeWriter
}

func (r *grpcReceiver) GetOffset(ctx context.Context) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func NewGRPCReceiver(stream grpc.ServerStream) fileupload.Receiver {
	pr, pw := io.Pipe()
	return &grpcReceiver{
		stream: stream,
		pr:     pr,
		pw:     pw,
	}
}

func (r *grpcReceiver) GetFileHeader(ctx context.Context) (fileupload.FileHeader, error) {
	r.server.CreateUploadTask(ctx, &fileuploadv1.CreateUploadTaskRequest{})
	if r.header != nil {
		return r.header, nil
	}
	var req fileuploadv1.UploadRequest
	err := r.stream.RecvMsg(&req)
	if err != nil {
		return nil, err
	}

	if !req.IsHeader {
		return nil, errors.New("expected fileHeader")
	}

	var header fileuploadv1.FileHeader
	if err := proto.Unmarshal(req.Data, &header); err != nil {
		return nil, err
	}

	r.header = &header
	return r.header, nil
}

func (r *grpcReceiver) ReceiveFile(ctx context.Context) (io.ReadCloser, error) {
	_, err := r.GetFileHeader(ctx)
	if err != nil {
		return nil, err
	}
	go func() {
		defer r.pw.Close()

		for {
			var req fileuploadv1.UploadRequest
			err := r.stream.RecvMsg(&req)
			if err == io.EOF {
				return
			}
			if err != nil {
				r.pw.CloseWithError(err)
				return
			}

			if req.IsHeader {
				r.pr.CloseWithError(errors.New("unexpected file header"))
			}

			if _, err := r.pw.Write(req.Data); err != nil {
				r.pw.CloseWithError(err)
				return
			}
		}
	}()

	return r.pr, nil
}

func (r *grpcReceiver) Finalize(ctx context.Context, resp fileupload.UploadResponse) error {
	return r.stream.SendMsg(resp)
}

func NewGRPCServerWithStorage(ctx context.Context, service *configv1.S) (fileuploadv1.FileUploadServiceServer,
	error) {

}
