/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package fileupload implements the functions, types, and interfaces for the module.
package fileupload

import (
	"testing"
)

func TestGenerateRandomHash(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		// TODO: Add test cases.
		{
			name: "test",
			want: 32,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateRandomHash(); len(got) != tt.want {
				t.Errorf("GenerateRandomHash() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestGenerateFileNameHash(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				name: "test.jpg",
			},
			want: "0d407ee6406a1216f2366674a1a9ff71361d5bef47021f8eb8b51f95e319dd56.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateFileHash(tt.args.name); got != tt.want {
				t.Errorf("GenerateFileNameHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
