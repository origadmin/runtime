package fileupload

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/origadmin/runtime"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/toolkits/fileupload"
)

type httpUploader struct {
	builder *Builder
	client  *http.Client
	request *http.Request
	body    io.ReadCloser
	header  fileupload.FileHeader
	uri     string
	buf     []byte
	offset  int64
}

func (u *httpUploader) SetFileHeader(ctx context.Context, header fileupload.FileHeader) error {
	log.Printf("Setting file fileHeader: %+v", header)
	u.header = header

	// Create new request
	req, err := http.NewRequestWithContext(ctx, "POST", u.uri, nil)
	if err != nil {
		log.Printf("Error creating new request: %v", err)
		return err
	}
	// Set headers
	req.Header.Set("Content-Type", header.GetContentType())
	req.Header.Set("Content-Length", fmt.Sprintf("%d", header.GetSize()))
	for k, v := range header.GetHeader() {
		req.Header.Set(k, v)
	}

	u.request = req
	log.Printf("File fileHeader set successfully")
	return nil
}

func (u *httpUploader) UploadFile(ctx context.Context, rd io.Reader) error {
	log.Printf("Uploading file...")
	if u.request == nil {
		log.Printf("Invalid request: request is nil")
		return ErrInvalidRequest
	}

	// Set the resumable upload Range fileHeader
	if u.offset > 0 {
		u.request.Header.Set("Range", fmt.Sprintf("bytes=%d-", u.offset))
		log.Printf("Setting Range fileHeader: bytes=%d-", u.offset)
	}

	u.request.Body = io.NopCloser(rd)
	if u.buf == nil {
		u.buf = u.builder.NewBuffer()
		log.Printf("Allocated new buffer: %v", len(u.buf))
	}

	if u.client == nil {
		u.client = &http.Client{
			//Timeout: u.builder.Timeout(),
		}
		log.Printf("Created new HTTP client: %+v", u.client)
	}
	//log.Printf("Uploading file with request: %+v", u.request)
	resp, err := u.client.Do(u.request)
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return err
	}
	log.Printf("Received response: %+v", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Invalid response status code: %d", resp.StatusCode)
		return ErrInvalidReceiverResponse
	}
	u.body = resp.Body
	log.Printf("File uploaded successfully: %+v", resp)
	return nil
}

func (u *httpUploader) Finalize(ctx context.Context) (fileupload.UploadResponse, error) {
	if u.buf != nil {
		log.Printf("Releasing buffer: %v", len(u.buf))
		u.builder.Free(u.buf)
		u.buf = nil
	}
	log.Printf("Finalizing upload...")
	var resp httpFileResponse
	if u.body == nil {
		log.Printf("Invalid response: response is nil")
		return &resp, ErrInvalidReceiverResponse
	}
	decoder := json.NewDecoder(u.body)
	if err := decoder.Decode(&resp); err != nil {
		log.Printf("Error decoding response: %v", err)
		return &resp, err
	}
	defer u.body.Close()
	log.Printf("Upload finalized successfully: %+v", resp)
	return &resp, nil
}

func (u *httpUploader) Resume(ctx context.Context, offset int64) error {
	log.Printf("Resuming upload at offset: %d", offset)
	u.offset = offset
	return nil
}

func NewHTTPUploader(ctx context.Context, service *configv1.Service) (fileupload.Uploader, error) {
	log.Printf("Creating new HTTP uploader...")
	clientService, err := runtime.NewHTTPServiceClient(ctx, service)
	if err != nil {
		log.Printf("Error creating HTTP service client: %v", err)
		return nil, err
	}
	_ = clientService
	return nil, errors.New("not implemented")
	//beyondcorp.ClientConnectorService{}
	//client := fileuploadv1.NewFileUploadServiceClient(beyondcorp.ClientConnectorService{})
	//return &httpUploader{
	//	builder: builder,
	//	uri:     builder.uri,
	//}
}

func newUploader(ctx context.Context, builder *Builder) fileupload.Uploader {
	log.Printf("Creating new uploader...")
	//clientConfig, err := runtime.NewHTTPServiceClient(ctx, service)
	//if err != nil {
	//	return nil, err
	//}
	//clientConfig
	//client := fileuploadv1.NewFileUploadServiceClient(clientConfig)
	return &httpUploader{
		builder: builder,
		uri:     builder.uri,
		client: &http.Client{
			Timeout: builder.Timeout(),
		},
	}
}
