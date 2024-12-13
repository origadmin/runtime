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
			want: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08.jpg",
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
