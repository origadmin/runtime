/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

/*
import (
	"testing"

	"github.com/vmihailenco/msgpack/v5"

	metainterfaces "github.com/origadmin/runtime/interfaces/store/meta"
	metav1 "github.com/origadmin/runtime/store/meta/v1"
	metav2 "github.com/origadmin/runtime/store/meta/v2"

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
			input: &metainterfaces.FileMetaData[metav1.FileMetaV1]{
				Info: FileIndexEntry{EntryName: "test1.txt"},
				Data: &metav1.FileMetaV1{
					FileSize:     1024,
					MimeType: "text/plain",
					ModifyTime:  1717182000,
				},
			},
			wantErr: false,
		},
		{
			name: "Value of FileMetaV1",
			input: metainterfaces.FileMetaData[metav1.FileMetaV1]{
				Info: FileIndexEntry{EntryName: "test2.txt"},
				Data: &metav1.FileMetaV1{
					FileSize:     2048,
					MimeType: "image/png",
					ModifyTime:  1717182001,
				},
			},
			wantErr: false,
		},
		{
			name: "Pointer to FileMetaV2",
			input: &metainterfaces.FileMetaData[metav2.FileMetaV2]{
				Info: FileIndexEntry{EntryName: "test3.json"},
				Data: &metav2.FileMetaV2{
					FileSize:        4096,
					MimeType:    "application/json",
					ModifyTime:     1717182002,
					BlockSize:   1024,
					BlockHashes: []string{"h5", "h6"},
				},
			},
			wantErr: false,
		},
		{
			name: "Value of FileMetaV2",
			input: metainterfaces.FileMetaData[metav2.FileMetaV2]{
				Info: FileIndexEntry{EntryName: "test4.mp4"},
				Data: &metav2.FileMetaV2{
					FileSize:        8192,
					MimeType:    "video/mp4",
					ModifyTime:     1717182003,
					BlockSize:   2048,
					BlockHashes: []string{"h7", "h8"},
				},
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
*/
