/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta defines the interfaces for the storage meta module.
package meta

import (
	"io/fs"
	"time"
)

// FileMeta 是文件内容元数据的通用接口，所有版本都必须实现。
// 它只包含文件内容固有的元数据，不涉及文件系统路径或权限。
type FileMeta interface {
	// Size() 返回文件内容的字节大小。
	Size() int64

	// ModTime() 返回文件内容本身的最后修改时间。
	ModTime() time.Time

	// BlobRef() 返回指向 Blob 存储中实际文件内容的引用。
	// 对于 V1，这将是单个 Blob ID (内容哈希)。
	// 对于 V2，这将是一个更复杂的结构，例如一个指向分片清单 Blob 的 ID。
	BlobRef() string

	// Version() 返回此元数据记录的版本号。
	Version() int

	// GetExtension() 用于获取文件内容相关的扩展属性。
	GetExtension(key string) (interface{}, bool)
	// SetExtension() 用于设置文件内容相关的扩展属性。
	SetExtension(key string, value interface{})
}

// FileMetaData 是一个泛型包装器，用于统一处理不同版本的文件元数据。
// T 必须是包含文件元数据字段的具体类型（例如 metav1.FileMetaV1, metav2.FileMetaV2）。
type FileMetaData[T any] struct {
	Info       FileInfo               // File system level info (e.g., Name, Mode, IsDir, Hash)
	Data       T                      // Version-specific file meta data
	extensions map[string]interface{} // Extensions for additional metadata
}

// Implement FileMeta interface for FileMetaData[T]
func (f *FileMetaData[T]) Size() int64 {
	// This requires T to have a FileSize() method or a FileSize field
	// We need to cast T to a known type or use reflection, which is complex.
	// For now, let's assume T has a FileSize field that can be accessed.
	// This will be refined when we implement specific versions.
	// For now, we'll use a placeholder.
	return 0 // Placeholder
}

func (f *FileMetaData[T]) ModTime() time.Time {
	// Placeholder
	return time.Time{}
}

func (f *FileMetaData[T]) BlobRef() string {
	// Placeholder
	return ""
}

func (f *FileMetaData[T]) Version() int {
	// Placeholder
	return 0
}

func (f *FileMetaData[T]) GetExtension(key string) (interface{}, bool) {
	if f.extensions == nil {
		return nil, false
	}
	val, ok := f.extensions[key]
	return val, ok
}

func (f *FileMetaData[T]) SetExtension(key string, value interface{}) {
	if f.extensions == nil {
		f.extensions = make(map[string]interface{})
	}
	f.extensions[key] = value
}

// FileInfo 接口定义了文件信息，兼容 io/fs.FileInfo。
// 它将由 Index 模块中的 DirEntry 实现，包含文件系统层面的信息。
type FileInfo interface {
	fs.FileInfo
	// 可以在这里添加其他 Index 模块特有的文件信息方法
}
