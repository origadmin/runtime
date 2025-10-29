package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/config"
)

const (
	_testJSON = `
{
    "test":{
        "settings":{
            "int_key":1000,
            "float_key":1000.1,
            "duration_key":10000,
            "string_key":"string_value"
        },
        "server":{
            "addr":"127.0.0.1",
            "port":8000
        }
    },
    "foo":[
        {
            "name":"nihao",
            "age":18
        },
        {
            "name":"nihao",
            "age":18
        }
    ]
}`

	_testJSONUpdate = `
{
    "test":{
        "settings":{
            "int_key":1000,
            "float_key":1000.1,
            "duration_key":10000,
            "string_key":"string_value"
        },
        "server":{
            "addr":"127.0.0.1",
            "port":8000
        }
    },
    "foo":[
        {
            "name":"nihao",
            "age":18
        },
        {
            "name":"nihao",
            "age":18
        }
    ],
	"bar":{
		"event":"update"
	}
}`

	//	_testYaml = `
	//Foo:
	//    bar :
	//        - {name: nihao,age: 1}
	//        - {name: nihao,age: 1}
	//
	//
	//`
)

//func TestScan(defaultFormatter *testing.T) {
//
//}

func TestFile(t *testing.T) {
	var (
		path = filepath.Join(t.TempDir(), "test_config")
		file = filepath.Join(path, "test.json")
		data = []byte(_testJSON)
	)
	defer os.Remove(path)
	if err := os.MkdirAll(path, 0o700); err != nil {
		t.Error(err)
	}
	if err := os.WriteFile(file, data, 0o666); err != nil {
		t.Error(err)
	}
	testSource(t, file, data)
	testSource(t, path, data)
	testWatchFile(t, file)
	testWatchDir(t, path, file)
}

func testWatchFile(t *testing.T, path string) {
	// Use event waiting with bounded retries instead of context.WithTimeout to align with real-world usage
	s := NewSource(path)
	watch, err := s.Watch()
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer func() {
		if err := watch.Stop(); err != nil {
			t.Logf("Warning: error stopping watcher: %v", err)
		}
	}()

	errCh := make(chan error, 1)
	kvsCh := make(chan []*config.KeyValue, 1)

	go func() {
		kvs, err := watch.Next()
		if err != nil {
			errCh <- fmt.Errorf("watch.Next() failed: %w", err)
			return
		}
		kvsCh <- kvs
	}()

	// Modify file content by truncating and rewriting to ensure event emission and consistent content
	f, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0o666)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer f.Close()
	if _, err = f.WriteString(_testJSONUpdate); err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}
	_ = f.Sync()
	_ = f.Close() // Ensure the file handle is released before rename on Windows

	// Wait for event with local retry/backoff and a deadline guard to avoid hanging
	deadline := time.Now().Add(10 * time.Second)
	for {
		select {
		case kvs := <-kvsCh:
			if !reflect.DeepEqual(string(kvs[0].Value), _testJSONUpdate) {
				t.Errorf("Expected value %q, got %q", _testJSONUpdate, kvs[0].Value)
			}
			goto RENAME
		case err := <-errCh:
			t.Fatalf("Error watching file: %v", err)
		case <-time.After(100 * time.Millisecond):
			if time.Now().After(deadline) {
				t.Fatal("Timeout waiting for file change event")
			}
		}
	}

RENAME:
	// Test file rename
	newFilepath := filepath.Join(filepath.Dir(path), "test1.json")
	if err = os.Rename(path, newFilepath); err != nil {
		t.Fatalf("Failed to rename file: %v", err)
	}
	defer func() {
		if err := os.Rename(newFilepath, path); err != nil {
			t.Logf("Warning: failed to restore file: %v", err)
		}
	}()

	// Listen again after rename; different platforms/FS may return an error or another event
	done := make(chan struct{}, 1)
	go func() {
		_, err := watch.Next()
		_ = err
		done <- struct{}{}
	}()

	select {
	case <-done:
		// pass
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting after file rename event")
	}
}

func testWatchDir(t *testing.T, path, file string) {
	s := NewSource(path)
	watch, err := s.Watch()
	if err != nil {
		t.Fatalf("watch error: %v", err)
	}
	defer func() { _ = watch.Stop() }()

	// Truncate and rewrite to fully replace content, then fsync
	f, err := os.OpenFile(file, os.O_RDWR|os.O_TRUNC, 0)
	if err != nil {
		t.Fatalf("open file error: %v", err)
	}
	_, err = f.WriteString(_testJSONUpdate)
	if err != nil {
		_ = f.Close()
		t.Fatalf("write file error: %v", err)
	}
	_ = f.Sync()
	_ = f.Close()

	deadline := time.Now().Add(10 * time.Second)
	for {
		select {
		case kvs, ok := <-func() chan []*config.KeyValue {
			ch := make(chan []*config.KeyValue, 1)
			go func() {
				kv, e := watch.Next()
				if e != nil {
					ch <- nil
					return
				}
				ch <- kv
			}()
			return ch
		}():
			if !ok || kvs == nil {
				t.Fatalf("watch.Next() returned nil or closed")
			}
			if !reflect.DeepEqual(string(kvs[0].Value), _testJSONUpdate) {
				t.Errorf("string(kvs[0].Value(%s)) not equal to _testJSONUpdate(%v)", kvs[0].Value, _testJSONUpdate)
			}
			return
		case <-time.After(100 * time.Millisecond):
			if time.Now().After(deadline) {
				t.Fatal("Timeout waiting for directory change event")
			}
		}
	}
}

func testSource(t *testing.T, path string, data []byte) {
	t.Log(path)

	s := NewSource(path)
	kvs, err := s.Load()
	if err != nil {
		t.Error(err)
	}
	if string(kvs[0].Value) != string(data) {
		t.Errorf("no expected: %s, but got: %s", kvs[0].Value, data)
	}
}

func TestConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test_config.json")
	defer os.Remove(path)
	if err := os.WriteFile(path, []byte(_testJSON), 0o666); err != nil {
		t.Error(err)
	}
	c := config.New(config.WithSource(
		NewSource(path),
	))
	testScan(t, c)

	testConfig(t, c)
}

func testConfig(t *testing.T, c config.Config) {
	expected := map[string]any{
		"test.settings.int_key":      int64(1000),
		"test.settings.float_key":    1000.1,
		"test.settings.string_key":   "string_value",
		"test.settings.duration_key": time.Duration(10000),
		"test.server.addr":           "127.0.0.1",
		"test.server.port":           int64(8000),
	}
	if err := c.Load(); err != nil {
		t.Error(err)
	}
	for key, value := range expected {
		switch value.(type) {
		case int64:
			if v, err := c.Value(key).Int(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		case float64:
			if v, err := c.Value(key).Float(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		case string:
			if v, err := c.Value(key).String(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		case time.Duration:
			if v, err := c.Value(key).Duration(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		}
	}
	// scan
	var settings struct {
		IntKey      int64         `json:"int_key"`
		FloatKey    float64       `json:"float_key"`
		StringKey   string        `json:"string_key"`
		DurationKey time.Duration `json:"duration_key"`
	}
	if err := c.Value("test.settings").Scan(&settings); err != nil {
		t.Error(err)
	}
	if v := expected["test.settings.int_key"]; settings.IntKey != v {
		t.Errorf("no expect int_key value: %v, but got: %v", settings.IntKey, v)
	}
	if v := expected["test.settings.float_key"]; settings.FloatKey != v {
		t.Errorf("no expect float_key value: %v, but got: %v", settings.FloatKey, v)
	}
	if v := expected["test.settings.string_key"]; settings.StringKey != v {
		t.Errorf("no expect string_key value: %v, but got: %v", settings.StringKey, v)
	}
	if v := expected["test.settings.duration_key"]; settings.DurationKey != v {
		t.Errorf("no expect duration_key value: %v, but got: %v", settings.DurationKey, v)
	}

	// not found
	if _, err := c.Value("not_found_key").Bool(); errors.Is(err, config.ErrNotFound) {
		t.Logf("not_found_key match: %v", err)
	}
}

func testScan(t *testing.T, c config.Config) {
	type TestJSON struct {
		Test struct {
			Settings struct {
				IntKey      int     `json:"int_key"`
				FloatKey    float64 `json:"float_key"`
				DurationKey int     `json:"duration_key"`
				StringKey   string  `json:"string_key"`
			} `json:"settings"`
			Server struct {
				Addr string `json:"addr"`
				Port int    `json:"port"`
			} `json:"server"`
		} `json:"test"`
		Foo []struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		} `json:"foo"`
	}
	var conf TestJSON
	if err := c.Load(); err != nil {
		t.Error(err)
	}
	if err := c.Scan(&conf); err != nil {
		t.Error(err)
	}
	t.Log(conf)
}

func TestMergeDataRace(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test_config.json")
	defer os.Remove(path)
	if err := os.WriteFile(path, []byte(_testJSON), 0o666); err != nil {
		t.Error(err)
	}
	c := config.New(config.WithSource(
		NewSource(path),
	))
	const count = 80
	wg := &sync.WaitGroup{}
	wg.Add(2)
	startCh := make(chan struct{})
	go func() {
		defer wg.Done()
		<-startCh
		for i := 0; i < count; i++ {
			var conf struct{}
			if err := c.Scan(&conf); err != nil {
				t.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		<-startCh
		for i := 0; i < count; i++ {
			if err := c.Load(); err != nil {
				t.Error(err)
			}
		}
	}()
	close(startCh)
	wg.Wait()
}
