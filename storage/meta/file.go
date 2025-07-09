// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"fmt"
	"io/fs"

	"github.com/origadmin/runtime/interfaces/storage/blob"
	"github.com/origadmin/runtime/interfaces/storage/meta"
	metav1 "github.com/origadmin/runtime/storage/meta/v1"
	metav2 "github.com/origadmin/runtime/storage/meta/v2"
)

func OpenMetaFile(storage blob.Store, fileMeta any) (fs.File, error) {
	switch v := fileMeta.(type) {
	case *meta.FileMetaData[metav2.FileMetaV2]: // v2 支持分片
		return metav2.NewMetaFileV2(storage, v)
	case *meta.FileMetaData[metav1.FileMetaV1]: // v1 不分片
		return metav1.NewMetaFileV1(storage, v)
	default:
		return nil, fmt.Errorf("unsupported metaFile meta version")
	}
}
