/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package fileupload implements the functions, types, and interfaces for the module.
package fileupload

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/goexts/generic/settings"
	"github.com/google/uuid"

	"github.com/origadmin/toolkits/fileupload"
)

const (
	bufSize = 32 * 1024
)

type Builder struct {
	uri         string
	hash        func(string) string
	bufPool     *sync.Pool
	bufSize     int
	serviceType fileupload.ServiceType
	timeout     time.Duration
}

func (b *Builder) Init(ss ...BuilderOption) *Builder {
	oldBufSize := b.bufSize
	settings.Apply(b, ss)
	// Initialize the buffer pool
	if b.bufPool == nil || b.bufSize != oldBufSize {
		b.bufPool = &sync.Pool{
			New: func() interface{} {
				return make([]byte, b.bufSize)
			},
		}
	}
	return b
}

func (b *Builder) NewUploader(ctx context.Context) fileupload.Uploader {
	//	switch b.serviceType {
	//	case ServiceTypeGRPC:
	//		return NewGRPCUploader(ctx, b.uri)
	//	default:
	//		return NewHTTPUploader(ctx, b.uri)
	//	}
	//	if b.serviceType >= serviceTypeMax {
	//		return nil, ErrInvalidServiceType
	//	}
	//	return b.services[b.serviceType], nil
	//}
	//
	//func (b *Builder) NewReceiver(ctx context.Context) Receiver {
	//	switch b.serviceType {
	//	case ServiceTypeGRPC:
	//		return NewGRPCReceiver(ctx)
	//	default:
	//		return NewHTTPReceiver(ctx)
	//	}
	return newUploader(ctx, b)
}

func (b *Builder) Free(buf []byte) {
	buf = buf[:0]
	b.bufPool.Put(buf)
}

func (b *Builder) NewBuffer() []byte {
	buf := b.bufPool.Get().([]byte)
	return buf
}

func (b *Builder) Timeout() time.Duration {
	return b.timeout
}

func (b *Builder) NewReceiver(r *http.Request, w http.ResponseWriter) fileupload.Receiver {
	return newReceiver(b, r, w)
}

type BuilderOption = func(builder *Builder)

func WithURI(uri string) BuilderOption {
	return func(builder *Builder) {
		builder.uri = uri
	}
}

func WithHash(hash func(name string) string) BuilderOption {
	return func(builder *Builder) {
		builder.hash = hash
	}
}

func GenerateHash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func GenerateFileHash(filename string) string {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	return GenerateHash(name) + ext
}

func GenerateRandomHash() string {
	id := uuid.Must(uuid.NewRandom())
	return hex.EncodeToString(id[:])
}

// NewBuilder creates a new httpBuilder with the given options
func NewBuilder(ss ...BuilderOption) *Builder {
	b := &Builder{
		hash:    GenerateHash,
		bufSize: bufSize, // default 32kb
		timeout: 5 * time.Minute,
	}
	return b.Init(ss...)
}

func WithBufferSize(size int) BuilderOption {
	return func(builder *Builder) {
		builder.bufSize = size
	}
}

func WithServiceType(st fileupload.ServiceType) BuilderOption {
	return func(builder *Builder) {
		builder.serviceType = st
	}
}

func WithTimeout(timeout time.Duration) BuilderOption {
	return func(builder *Builder) {
		builder.timeout = timeout
	}
}
