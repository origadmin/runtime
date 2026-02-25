package file

import (
	"context"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"

	"github.com/go-kratos/kratos/v2/config"
)

var _ config.Watcher = (*watcher)(nil)
var _ config.Watcher = (*emptyWatcher)(nil)

// emptyWatcher is a no-op watcher for optional files that don't exist
type emptyWatcher struct{}

func (w *emptyWatcher) Next() ([]*config.KeyValue, error) {
	return []*config.KeyValue{}, nil
}

func (w *emptyWatcher) Stop() error {
	return nil
}

type watcher struct {
	f  *file
	fw *fsnotify.Watcher

	ctx    context.Context
	cancel context.CancelFunc
}

func (w *watcher) Next() ([]*config.KeyValue, error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case event := <-w.fw.Events:
		if event.Op == fsnotify.Rename {
			if w.f.shouldIgnore(event.Name) {
				return nil, nil
			}
			if _, err := os.Stat(event.Name); err == nil || os.IsExist(err) {
				if err := w.fw.Add(event.Name); err != nil {
					return nil, err
				}
			}
		}

		// Apply filter and ignores to the event file
		fileName := filepath.Base(event.Name)
		if w.f.shouldIgnore(fileName) {
			return nil, nil
		}
		if w.f.filter != "" {
			matched, err := filepath.Match(w.f.filter, fileName)
			if err != nil {
				return nil, err
			}
			if !matched {
				return nil, nil
			}
		}

		// Trigger reload of the entire source
		// Kratos config expects the full list of KVs from the source on each change
		return w.f.Load()
	case err := <-w.fw.Errors:
		return nil, err
	}
}

func (w *watcher) Stop() error {
	w.cancel()
	return w.fw.Close()
}

func newWatcher(f *file) (config.Watcher, error) {
	if f.optional {
		if _, err := os.Stat(f.path); os.IsNotExist(err) {
			return &emptyWatcher{}, nil
		}
	}
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := fw.Add(f.path); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &watcher{f: f, fw: fw, ctx: ctx, cancel: cancel}, nil
}
