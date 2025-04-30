/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package envars implements the functions, types, and interfaces for the module.
package envars

import (
	"os"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/stretchr/testify/assert"
)

func TestNewSource(t *testing.T) {
	tests := []struct {
		name     string
		envVars  []string
		prefixes []string
		expected []*config.KeyValue
	}{
		{
			name:     "no conditions",
			envVars:  []string{"FOO=bar", "BAZ=qux"},
			prefixes: []string{},
			expected: []*config.KeyValue{
				{Key: "FOO", Value: []byte("bar"), Format: ""},
				{Key: "BAZ", Value: []byte("qux"), Format: ""},
			},
		},
		{
			name:     "with conditions",
			envVars:  []string{"FOO=bar", "BAZ=qux", "TEST=123"},
			prefixes: []string{"FOO", "TEST"},
			expected: []*config.KeyValue{},
		},
		{
			name:     "condition not found",
			envVars:  []string{"FOO=bar", "BAZ=qux"},
			prefixes: []string{"NOT_EXIST"},
			expected: []*config.KeyValue{},
		},
		{
			name:     "empty environment",
			envVars:  []string{},
			prefixes: []string{},
			expected: []*config.KeyValue{},
		},
		{
			name:     "malformed env var",
			envVars:  []string{"MALFORMED"},
			prefixes: []string{},
			expected: []*config.KeyValue{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			originalEnv := os.Environ()
			os.Clearenv()
			defer func() {
				os.Clearenv()
				for _, e := range originalEnv {
					parts := strings.SplitN(e, "=", 2)
					if len(parts) == 2 {
						os.Setenv(parts[0], parts[1])
					}
				}
			}()

			// Correctly set environment variables
			for _, env := range tt.envVars {
				parts := strings.SplitN(env, "=", 2)
				if len(parts) == 2 {
					os.Setenv(parts[0], parts[1])
				}
			}

			// Test
			source := NewSource(tt.prefixes...)
			envarsSource, ok := source.(*envars)
			assert.True(t, ok)

			// Call the Load method to load data
			_, err := envarsSource.Load()
			assert.NoError(t, err)

			// Use ElementsMatch to compare slices, ignoring order
			assert.ElementsMatch(t, tt.expected, envarsSource.data, "Mismatch in KeyValue elements")
		})
	}
}

func TestEnvars_Load(t *testing.T) {
	// 保存原始环境变量
	originalEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, e := range originalEnv {
			parts := strings.SplitN(e, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
	}()

	tests := []struct {
		name     string
		envVars  map[string]string
		prefixes []string
		expected []*config.KeyValue
	}{
		{
			name: "single prefix match",
			envVars: map[string]string{
				"APP_DB_HOST": "localhost",
				"APP_DB_PORT": "5432",
				"APP_DB_USER": "user",
			},
			prefixes: []string{"APP_"},
			expected: []*config.KeyValue{
				{Key: "DB_HOST", Value: []byte("localhost")},
				{Key: "DB_PORT", Value: []byte("5432")},
				{Key: "DB_USER", Value: []byte("user")},
			},
		},
		{
			name: "multiple prefixes",
			envVars: map[string]string{
				"APP_DB_HOST": "localhost",
				"APP_DB_PORT": "5432",
				"APP_DB_USER": "user",
				"OTHER_VAR":   "value",
			},
			prefixes: []string{"APP_", "OTHER_"},
			expected: []*config.KeyValue{
				{Key: "DB_HOST", Value: []byte("localhost")},
				{Key: "DB_PORT", Value: []byte("5432")},
				{Key: "DB_USER", Value: []byte("user")},
				{Key: "VAR", Value: []byte("value")},
			},
		},
		{
			name:     "no matching prefix",
			envVars:  map[string]string{},
			prefixes: []string{"NONE_"},
			expected: []*config.KeyValue{},
		},
		{
			name: "empty prefix",
			envVars: map[string]string{
				"APP_DB_HOST":   "localhost",
				"APP_DB_PORT":   "5432",
				"APP_DB_USER":   "user",
				"OTHER_VAR":     "value",
				"NO_PREFIX_VAR": "no_prefix",
			},
			prefixes: []string{""},
			expected: []*config.KeyValue{
				{Key: "APP_DB_HOST", Value: []byte("localhost")},
				{Key: "APP_DB_PORT", Value: []byte("5432")},
				{Key: "APP_DB_USER", Value: []byte("user")},
				{Key: "OTHER_VAR", Value: []byte("value")},
				{Key: "NO_PREFIX_VAR", Value: []byte("no_prefix")},
			},
		},
		{
			name: "multiple prefixes with one matching",
			envVars: map[string]string{
				"APP_DB_HOST": "localhost",
				"APP_DB_PORT": "5432",
				"APP_DB_USER": "user",
			},
			prefixes: []string{"APP_", "NONE_"},
			expected: []*config.KeyValue{
				{Key: "DB_HOST", Value: []byte("localhost")},
				{Key: "DB_PORT", Value: []byte("5432")},
				{Key: "DB_USER", Value: []byte("user")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清除所有环境变量
			os.Clearenv()

			// 设置测试用例所需的环境变量
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			e := NewSource(tt.prefixes...)
			actual, err := e.Load()
			assert.NoError(t, err)

			// 使用 assert.ElementsMatch 忽略顺序进行比较
			assert.ElementsMatch(t, tt.expected, actual)
		})
	}
}

func TestLoadEnviron(t *testing.T) {
	tests := []struct {
		name     string
		data     []string
		prefixes []string
		expected []*config.KeyValue
	}{
		{
			name:     "empty inputs",
			data:     []string{},
			prefixes: []string{},
			expected: []*config.KeyValue{},
		},
		{
			name:     "no matching conditions",
			data:     []string{"FOO=bar", "BAZ=qux"},
			prefixes: []string{"NONE"},
			expected: []*config.KeyValue{},
		},
		{
			name:     "exact match, key empty after trim",
			data:     []string{"FOO=bar", "BAZ=qux"},
			prefixes: []string{"FOO"},
			expected: []*config.KeyValue{},
		},
		{
			name:     "case insensitive match",
			data:     []string{"Foo=bar", "BAZ=qux"},
			prefixes: []string{"foo"},
			expected: []*config.KeyValue{},
		},
		{
			name:     "multiple matches",
			data:     []string{"FOO=bar", "BAZ=qux", "FOO2=bar2"},
			prefixes: []string{"FOO", "FOO2"},
			expected: []*config.KeyValue{
				{
					Key:   "2",
					Value: []byte("bar2"),
				},
			},
		},
		{
			name: "multiple prefixes",
			data: []string{
				"APP_DB_HOST=localhost",
				"APP_DB_PORT=5432",
				"APP_DB_USER=user",
				"OTHER_VAR=value",
				"NO_PREFIX_VAR=no_prefix",
			},
			prefixes: []string{"APP_", "OTHER_"},
			expected: []*config.KeyValue{
				{
					Key:   "DB_HOST",
					Value: []byte("localhost"),
				},
				{
					Key:   "DB_PORT",
					Value: []byte("5432"),
				},
				{
					Key:   "DB_USER",
					Value: []byte("user"),
				},
				{
					Key:   "VAR",
					Value: []byte("value"),
				},
			},
		},
		{
			name:     "empty prefix returns all",
			data:     []string{"FOO=bar", "BAZ=qux"},
			prefixes: []string{""},
			expected: []*config.KeyValue{
				{Key: "FOO", Value: []byte("bar")},
				{Key: "BAZ", Value: []byte("qux")},
			},
		},
		{
			name:     "case-insensitive prefix match",
			data:     []string{"FooBar=baz"},
			prefixes: []string{"foobar"},
			expected: []*config.KeyValue{},
		},
		{
			name:     "nested prefix match",
			data:     []string{"APP_CONFIG_KEY=value"},
			prefixes: []string{"APP_"},
			expected: []*config.KeyValue{
				{Key: "CONFIG_KEY", Value: []byte("value")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := loadEnviron(tt.data, tt.prefixes)
			assert.ElementsMatch(t, tt.expected, actual)
		})
	}
}

func TestMatchFold(t *testing.T) {
	tests := []struct {
		name     string
		envs     []string
		data     string
		expected string
		found    bool
	}{
		{
			name:     "empty envs and data",
			envs:     []string{},
			data:     "",
			expected: "",
			found:    false,
		},
		{
			name:     "empty envs",
			envs:     []string{},
			data:     "TEST",
			expected: "",
			found:    false,
		},
		{
			name:     "empty data",
			envs:     []string{"TEST"},
			data:     "",
			expected: "",
			found:    false,
		},
		{
			name:     "exact match",
			envs:     []string{"TEST"},
			data:     "TEST",
			expected: "TEST",
			found:    true,
		},
		{
			name:     "prefix match",
			envs:     []string{"TEST"},
			data:     "TEST_VALUE",
			expected: "TEST",
			found:    true,
		},
		{
			name:     "no match",
			envs:     []string{"TEST"},
			data:     "OTHER",
			expected: "",
			found:    false,
		},
		{
			name:     "multiple envs first match",
			envs:     []string{"TEST1", "TEST2"},
			data:     "TEST1_VALUE",
			expected: "TEST1",
			found:    true,
		},
		{
			name:     "multiple envs second match",
			envs:     []string{"TEST1", "TEST2"},
			data:     "TEST2_VALUE",
			expected: "TEST2",
			found:    true,
		},
		{
			name:     "case sensitive no match",
			envs:     []string{"test"},
			data:     "TEST",
			expected: "",
			found:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, found := matchPrefix(tt.envs, tt.data)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.found, found)
		})
	}
}
