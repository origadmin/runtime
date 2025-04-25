/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package fileupload implements the functions, types, and interfaces for the module.
package fileupload

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploaderDownloaderSequence(t *testing.T) {
	b := NewBuilder()
	// Create a test server for upload
	uploadTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Create a receiver
		receiver := b.NewReceiver(r, w)

		_, err := receiver.GetFileHeader(r.Context())
		if err != nil {
			return
		}

		// Receive the file
		reader, err := receiver.ReceiveFile(context.Background())
		if err != nil {
			t.Errorf("ReceiveFile failed: %v", err)
			return
		}
		defer reader.Close()

		// Read the downloaded content
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(reader)
		size := buf.Len()
		t.Logf("ReceiveFile successful: %+v", size)
		// Check the content
		content := buf.String()
		if content != "Hello, World!" {
			t.Errorf("Expected 'Hello, World!', got %s", content)
		}
		hash := GenerateFileHash("test.txt")
		resp := &httpFileResponse{
			Success: true,
			Path:    "test.txt",
			Size:    uint32(size),
			Hash:    hash,
		}

		// Finalize the receiver
		err = receiver.Finalize(context.Background(), resp)
		if err != nil {
			t.Errorf("Finalize failed: %v", err)
		}
		t.Logf("ReceiveFile successful: %+v", content)
	}))
	defer uploadTS.Close()

	//// Create a test server for download
	//downloadTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	if r.Method != http.MethodGet {
	//		t.Errorf("Expected GET request, got %s", r.Method)
	//	}
	//	w.Write([]byte("Hello, World!"))
	//}))
	//defer downloadTS.Close()
	b.Init(WithURI(uploadTS.URL))
	// Create an uploader
	uploader := b.NewUploader(context.Background())

	// Create a buffer to hold the form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a form file part
	part, err := writer.CreateFormFile("file", "testfile.txt")
	if err != nil {
		t.Fatalf("Failed to create form file part: %v", err)
	}

	// Create a test file
	testFile := []byte("Hello, World!")
	_, _ = part.Write(testFile)
	writer.Close()
	//reader := bytes.NewReader(testFile)
	header := httpFileHeader{
		Filename:    "testfile.txt",
		Size:        uint32(body.Len()),
		ContentType: writer.FormDataContentType(),
		Header: map[string]string{
			"Content-Length": fmt.Sprintf("%d", body.Len()),
		},
		IsDir: false,
	}
	err = uploader.SetFileHeader(context.Background(), &header)
	if err != nil {
		t.Errorf("UploadFile failed: %v", err)
	}
	// Upload the file
	t.Logf("Starting UploadFile...")
	err = uploader.UploadFile(context.Background(), body)
	if err != nil {
		t.Errorf("UploadFile failed: %v", err)
	}
	t.Logf("UploadFile completed.")

	// Finalize the upload
	resp, err := uploader.Finalize(context.Background())
	if err != nil {
		t.Errorf("Finalize failed: %v", err)
	}

	// Check the response
	if !resp.GetSuccess() {
		t.Errorf("Expected success, got failure")
	}
	t.Logf("Uploader successful: %+v", resp)
}

func TestMyGRPCService(t *testing.T) {

}
