/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package fileupload implements the functions, types, and interfaces for the module.
package fileupload

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/fileupload"
)

type HTTPReceiver struct {
	builder    *Builder
	fileHeader fileupload.FileHeader
	file       multipart.File
	header     *multipart.FileHeader
	response   http.ResponseWriter
	err        error
	pw         *io.PipeWriter
}

type httpFileResponse struct {
	Success    bool   `json:"success"`
	Hash       string `json:"hash,omitempty"`
	Path       string `json:"path,omitempty"`
	Size       uint32 `json:"size,omitempty"`
	FailReason string `json:"fail_reason,omitempty"`
}

func (r *httpFileResponse) String() string {
	bytes, _ := json.MarshalIndent(r, "", "  ")
	return string(bytes)
}

func (r *httpFileResponse) GetSuccess() bool {
	if r == nil {
		return false
	}
	return r.Success
}

func (r *httpFileResponse) GetHash() string {
	return r.Hash
}

func (r *httpFileResponse) GetPath() string {
	return r.Path
}

func (r *httpFileResponse) GetSize() uint32 {
	return r.Size
}

func (r *httpFileResponse) GetFailReason() string {
	return r.FailReason
}

// GetFileHeader read the file fileHeader from the request.
func (r *HTTPReceiver) GetFileHeader(ctx context.Context) (fileupload.FileHeader, error) {
	if r.err != nil {
		log.Printf("Error getting file fileHeader: %v", r.err)
		return nil, r.err
	}
	log.Printf("File fileHeader: %+v", r.fileHeader)
	return r.fileHeader, nil
}

// ReceiveFile read the file data to the server with path.
func (r *HTTPReceiver) ReceiveFile(ctx context.Context) (io.ReadCloser, error) {
	if r.err != nil {
		log.Printf("Error receiving file: %v", r.err)
		return nil, r.err
	}
	if r.file == nil {
		r.err = ErrNoFile
		log.Printf("No file provided")
		return nil, r.err
	}

	if rangeHeader := r.header.Header.Get("Range"); rangeHeader != "" {
		log.Printf("parsing range fileHeader: %v", rangeHeader)
		var offset int64
		_, err := fmt.Sscanf(rangeHeader, "bytes=%r-", &offset)
		if err != nil {
			log.Printf("Error parsing range fileHeader: %v", err)
			return nil, err
		}
		if _, err := r.file.Seek(offset, io.SeekStart); err != nil {
			log.Printf("Error seeking file: %v", err)
			return nil, err
		}
	}

	var pr *io.PipeReader
	pr, r.pw = io.Pipe()

	go func() {
		buf := r.builder.NewBuffer()
		defer func() {
			r.file.Close()
			r.builder.Free(buf)
		}()

		for {
			select {
			case <-ctx.Done():
				log.Printf("Context cancelled")
				r.pw.CloseWithError(ctx.Err())
				return
			default:

				n, err := r.file.Read(buf)
				log.Print("start reading file", n)
				if n > 0 {
					if _, werr := r.pw.Write(buf[:n]); werr != nil {
						log.Printf("Error writing to pipe: %v", werr)
						r.pw.CloseWithError(werr)
						return
					}
				}
				if err == io.EOF {
					r.pw.Close()
					return
				}
				if err != nil {
					log.Printf("Error reading file: %v", err)
					r.pw.CloseWithError(err)
					return
				}
			}
		}
	}()

	log.Printf("File received successfully")
	return pr, nil
}

// Finalize write the finalize status to the client and close the upload process.
func (r *HTTPReceiver) Finalize(ctx context.Context, resp fileupload.UploadResponse) error {
	if r.err != nil {
		log.Printf("Error finalizing upload: %v", r.err)
		return r.err
	}
	if resp.GetSuccess() {
		r.response.WriteHeader(http.StatusOK)
		// Write response headers
		r.response.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(r.response).Encode(resp)
	}
	r.response.WriteHeader(http.StatusInternalServerError)
	log.Printf("Failed to upload file: %v", resp.GetFailReason())
	return json.NewEncoder(r.response).Encode(map[string]string{
		"error": resp.GetFailReason(),
	})
}
func (r *HTTPReceiver) GetOffset(ctx context.Context) (int64, error) {
	if rangeHeader := r.header.Header.Get("Range"); rangeHeader != "" {
		var offset int64
		_, err := fmt.Sscanf(rangeHeader, "bytes=%r-", &offset)
		if err != nil {
			return 0, err
		}
		return offset, nil
	}
	return 0, nil
}

func newReceiver(builder *Builder, req *http.Request, resp http.ResponseWriter) *HTTPReceiver {
	// Read the file fileHeader from the request
	var receiver HTTPReceiver
	file, header, err := req.FormFile("file")
	if err != nil {
		receiver.err = errors.Wrap(err, "failed to get file from request")
		return &receiver
	}

	fileheader := make(map[string]string, len(header.Header))
	for k, v := range header.Header {
		fileheader[k] = v[0]
	}
	modTime := uint32(time.Now().Unix())
	if mod, err := time.Parse(fileupload.ModTimeFormat, header.Header.Get("Last-Modified")); err == nil {
		modTime = uint32(mod.Unix())
	}

	// Create a FileHeader struct and populate it with the file information
	fileHeader := &httpFileHeader{
		Filename:    header.Filename,
		Size:        uint32(header.Size),
		ContentType: header.Header.Get("Content-Type"),
		Header:      fileheader,
		ModTime:     modTime,
	}

	return &HTTPReceiver{
		builder:    builder,
		file:       file,
		response:   resp,
		fileHeader: fileHeader,
		header:     header,
		//request:  req,
	}
}

var _ fileupload.Receiver = &HTTPReceiver{}
