package fileupload

import (
	"context"
	"errors"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/origadmin/toolkits/fileupload"

	fileuploadv1 "github.com/origadmin/runtime/gen/go/fileupload/v1"
)

type grpcReceiver struct {
	header *fileuploadv1.FileHeader
	stream grpc.ServerStream
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
				continue
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
