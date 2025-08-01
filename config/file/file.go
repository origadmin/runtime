package file

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/goexts/generic/settings"
)

var _ config.Source = (*file)(nil)

// Temporary file suffixes that are ignored by default
var defaultIgnores = []string{
	// Linux
	"~",
	// macOS
	".DS_Store",
	".AppleDouble",
	".LSOverride",
	// Windows
	".tmp",
	".temp",
	".bak",
}

type file struct {
	path      string
	ignores   []string
	formatter Formatter
}

// NewSource new a file source.
func NewSource(path string, opts ...Option) config.Source {
	f := &file{
		path:      path,
		ignores:   defaultIgnores,
		formatter: defaultFormatter,
	}
	return settings.Apply(f, opts)
}

func (f *file) loadFile(path string) (*config.KeyValue, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if f.formatter != nil {
		return f.formatter(info.Name(), data)
	}
	return &config.KeyValue{
		Key:    info.Name(),
		Format: format(info.Name()),
		Value:  data,
	}, nil
}

func (f *file) shouldIgnore(filename string) bool {
	if len(f.ignores) == 0 {
		return false
	}
	ext := strings.ToLower(filepath.Ext(filename))
	for _, ignoreExt := range f.ignores {
		if strings.HasSuffix(ext, ignoreExt) {
			return true
		}
	}
	return false
}

func (f *file) loadDir(path string) (kvs []*config.KeyValue, err error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		// ignore hidden files
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") || f.shouldIgnore(file.Name()) {
			continue
		}
		kv, err := f.loadFile(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}
		kvs = append(kvs, kv)
	}
	return
}

func (f *file) Load() (kvs []*config.KeyValue, err error) {
	fi, err := os.Stat(f.path)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return f.loadDir(f.path)
	}

	if f.shouldIgnore(fi.Name()) {
		return nil, nil
	}
	kv, err := f.loadFile(f.path)
	if err != nil {
		return nil, err
	}
	return []*config.KeyValue{kv}, nil
}

func (f *file) Watch() (config.Watcher, error) {
	return newWatcher(f)
}

func defaultFormatter(key string, value []byte) (*config.KeyValue, error) {
	return &config.KeyValue{
		Key:    key,
		Format: format(key),
		Value:  value,
	}, nil
}
