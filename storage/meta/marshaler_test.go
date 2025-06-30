/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"testing"

	"github.com/vmihailenco/msgpack/v5"

	metav1 "github.com/origadmin/runtime/storage/meta/v1"
	metav2 "github.com/origadmin/runtime/storage/meta/v2"

	"github.com/stretchr/testify/assert"
)

func TestMarshalFileMeta(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name: "Pointer to FileMetaV1",
			input: &metav1.FileMetaV1{
				Hash:     "hash1",
				Size:     1024,
				MimeType: "text/plain",
				ModTime:  1717182000,
			},
			wantErr: false,
		},
		{
			name: "Value of FileMetaV1",
			input: metav1.FileMetaV1{
				Hash:     "hash2",
				Size:     2048,
				MimeType: "image/png",
				ModTime:  1717182001,
			},
			wantErr: false,
		},
		{
			name: "Pointer to FileMetaV2",
			input: &metav2.FileMetaV2{
				Hash:        "hash3",
				Size:        4096,
				MimeType:    "application/json",
				ModTime:     1717182002,
				BlockSize:   1024,
				BlockHashes: []string{"h5", "h6"},
				Extra:       map[string]string{"k1": "v1", "k2": "v2"},
			},
			wantErr: false,
		},
		{
			name: "Value of FileMetaV2",
			input: metav2.FileMetaV2{
				Hash:        "hash4",
				Size:        8192,
				MimeType:    "video/mp4",
				ModTime:     1717182003,
				BlockSize:   2048,
				BlockHashes: []string{"h7", "h8"},
				Extra:       map[string]string{"k1": "v1", "k2": "v2"},
			},
			wantErr: false,
		},
		{
			name:    "Unsupported type",
			input:   "invalid-type",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := MarshalFileMeta(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, data)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, data)

				// Optional: Deserialize back to verify content consistency
				var version struct {
					Version int32 `msgpack:"v"`
				}
				err = msgpack.Unmarshal(data, &version)
				assert.NoError(t, err)
				assert.NotZero(t, version.Version)
			}
		})
	}
}
