package fileupload

import (
	"net/http"

	"github.com/origadmin/toolkits/fileupload"
)

// BridgeUploader 实现了HTTP到gRPC的桥接上传
type BridgeUploader struct {
	Builder *Builder
	Service fileupload.Uploader
}

// NewBridgeUploader 创建一个新的桥接上传器
func NewBridgeUploader(builder *Builder, service fileupload.Uploader) fileupload.BridgeUploader {
	return &BridgeUploader{
		Builder: builder,
		Service: service,
	}
}

// ServeHTTP 处理HTTP上传请求并转发到gRPC服务
func (b *BridgeUploader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Create an HTTP receiver to receive uploaded files
	httpReceiver := b.Builder.NewReceiver(r, w)

	// 2. Obtain the file fileHeader information
	header, err := httpReceiver.GetFileHeader(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Create a gRPC uploader
	grpcUploader := b.Service

	// 4. Set the file fileHeader information
	if err := grpcUploader.SetFileHeader(ctx, header); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Get the contents of the uploaded file
	reader, err := httpReceiver.ReceiveFile(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	// 6. Upload the file to the gRPC service
	if err := grpcUploader.UploadFile(ctx, reader); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 7. Complete the upload and get the response
	resp, err := grpcUploader.Finalize(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 8. Forward the gRPC response back to the HTTP client
	httpReceiver.Finalize(ctx, resp)
}
