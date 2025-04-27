package fileupload

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/origadmin/runtime"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/toolkits/fileupload"
	fileuploadv1 "github.com/origadmin/toolkits/fileupload/gen/go/fileupload/v1"
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
	taskId  string
}

func (u *httpUploader) SetFileHeader(ctx context.Context, header fileupload.FileHeader) error {
	protoHeader := Header2GRPCHeader(header)
	reqBody, err := json.Marshal(map[string]interface{}{
		"file_header": protoHeader,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", u.uri+"/upload/create", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var respData fileuploadv1.CreateUploadTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return err
	}
	u.taskId = respData.TaskId
	return nil
}

func (u *httpUploader) UploadFile(ctx context.Context, rd io.Reader) error {
	chunkSize := u.builder.ChunkSize()
	var chunkNumber int32 = 0
	var wg sync.WaitGroup
	for {
		chunk := make([]byte, chunkSize)
		n, err := rd.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		wg.Add(1)
		go func(data []byte, chunkNum int32) {
			defer wg.Done()
			reqBody := bytes.NewReader(data[:n])
			req, _ := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/upload/%s/chunk", u.uri, u.taskId), reqBody)
			req.Header.Set("Content-Type", "application/octet-stream")
			req.Header.Set("Chunk-Number", fmt.Sprintf("%d", chunkNum))

			_, err := u.client.Do(req)
			if err != nil {
				log.Printf("Chunk %d upload failed: %v", chunkNum, err)
			}
		}(chunk, chunkNumber)
		chunkNumber++
	}
	wg.Wait()
	return nil
}

func (u *httpUploader) Finalize(ctx context.Context) (fileupload.UploadResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/upload/%s/finalize", u.uri, u.taskId), nil)
	if err != nil {
		return nil, err
	}
	resp, err := u.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respData fileuploadv1.UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}
	return &respData, nil
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
