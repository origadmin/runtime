/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"reflect"
	"testing"

	metav1 "github.com/origadmin/runtime/storage/meta/v1"
)

func TestMarshalFileMeta(t *testing.T) {
	type args struct {
		meta interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				meta: metav1.FileMetaV1{
					BlockHashes: []string{"block1", "block2"},
					BlockSize:   1024,
					Hash:        "hash1",
					MimeType:    "text/plain",
					ModTime:     1638048000,
					Size:        1024,
					Version:     1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalFileMeta(tt.args.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalFileMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			meta, err := UnmarshalFileMeta(got)
			if err != nil {
				t.Errorf("UnmarshalFileMeta() error = %v", err)
			}
			if !reflect.DeepEqual(meta, tt.args.meta) {
				t.Errorf("UnmarshalFileMeta() got = %v, want %v", meta, tt.args.meta)
			}
		})
	}
}
