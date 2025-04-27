package fileupload

import (
	"net/http"

	"github.com/origadmin/toolkits/fileupload"
)

// BridgeUploader Implemented HTTP to gRPC bridge upload
type BridgeUploader struct {
	Builder *Builder
	Service fileupload.Uploader
}

// NewBridgeUploader Create a new bridge uploader
func NewBridgeUploader(builder *Builder, service fileupload.Uploader) fileupload.BridgeUploader {
	return &BridgeUploader{
		Builder: builder,
		Service: service,
	}
}

// ServeHTTP Processes HTTP upload requests and forwards them to the gRPC service
func (b *BridgeUploader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	httpReceiver := b.Builder.NewReceiver(r, w)
	header, err := httpReceiver.GetFileHeader(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	grpcUploader := b.Service

	if err := grpcUploader.SetFileHeader(ctx, header); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reader, err := httpReceiver.ReceiveFile(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	if err := grpcUploader.UploadFile(ctx, reader); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := grpcUploader.Finalize(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	errFin := httpReceiver.Finalize(ctx, resp)
	if errFin != nil {
		http.Error(w, errFin.Error(), http.StatusInternalServerError)
	}
}
