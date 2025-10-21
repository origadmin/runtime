package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	sourcev1 "github.com/origadmin/runtime/api/gen/go/runtime/source/v1"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/interfaces/options"
)

var _ kratosconfig.Source = (*file)(nil)

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

// file represents a file source used to load configuration from the file system
type file struct {
	path      string
	ignores   []string
	formatter Formatter
	optional  bool
}

// NewSource creates a new file source instance
func NewSource(path string, opts ...Option) kratosconfig.Source {
	f := &file{
		path:      path,
		optional:  false,
		ignores:   defaultIgnores,
		formatter: defaultFormatter,
	}
	return applyFileOptions(f, opts...)
}

// loadFile loads a single file from the specified path
func (f *file) loadFile(path string) (*kratosconfig.KeyValue, error) {
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
	return &kratosconfig.KeyValue{
		Key:    info.Name(),
		Format: format(info.Name()),
		Value:  data,
	}, nil
}

// shouldIgnore determines whether a file should be ignored
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

// loadDir loads all non-ignored files from the specified directory
func (f *file) loadDir(path string) (kvs []*kratosconfig.KeyValue, err error) {
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

// Load loads configuration data from the file source
func (f *file) Load() (kvs []*kratosconfig.KeyValue, err error) {
	fi, err := os.Stat(f.path)
	if err != nil {
		if os.IsNotExist(err) && f.optional {
			// File doesn't exist, but it's optional, return empty config
			return []*kratosconfig.KeyValue{}, nil
		}
		abs, _ := filepath.Abs(f.path)
		return nil, fmt.Errorf("failed to stat config path %q", abs)
	}
	if fi.IsDir() {
		return f.loadDir(f.path)
	}

	if f.shouldIgnore(fi.Name()) {
		return nil, nil
	}
	kv, err := f.loadFile(f.path)
	if err != nil {
		if f.optional && (os.IsNotExist(err) || os.IsPermission(err)) {
			return []*kratosconfig.KeyValue{}, nil
		}
		return nil, err
	}

	if kv == nil {
		return []*kratosconfig.KeyValue{}, nil
	}
	return []*kratosconfig.KeyValue{kv}, nil
}

// Watch creates and returns a file watcher instance
func (f *file) Watch() (kratosconfig.Watcher, error) {
	return newWatcher(f)
}

// defaultFormatter is the default formatting function used to process key-value pair data
func defaultFormatter(key string, value []byte) (*kratosconfig.KeyValue, error) {
	return &kratosconfig.KeyValue{
		Key:    key,
		Format: format(key),
		Value:  value,
	}, nil
}

// NewFileSource creates a new file source based on configuration.
// It adapts to the new bootstrap mechanism.
func NewFileSource(cfg *sourcev1.SourceConfig, opts ...options.Option) (kratosconfig.Source, error) {
	fileSrc := cfg.GetFile()
	if fileSrc == nil {
		// This can happen if the source type is "file" but the `file` oneof is not set.
		// Returning nil, nil is a safe default, allowing other sources to proceed.
		return nil, nil
	}
	optional := fileSrc.GetOptional()
	if optional {
		opts = append(opts, WithOptional())
	}

	return NewSource(fileSrc.GetPath(), opts...), nil
}

// init registers the file source during package initialization
func init() {
	runtimeconfig.Register("file", NewFileSource)
}
