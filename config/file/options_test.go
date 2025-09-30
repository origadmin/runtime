package file

import (
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/stretchr/testify/assert"

	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/optionutil"
)

func TestFromOptions_NilOptions(t *testing.T) {
	// Test with nil options
	opts := applyFileOptions(&file{})
	assert.Empty(t, opts, "Expected empty options slice when options is nil")
}

func TestFromOptions_UninitializedOptions(t *testing.T) {
	// Test with uninitialized options
	option := optionutil.WithContext(nil)
	opts := applyFileOptions(&file{}, option)
	assert.Empty(t, opts, "Expected empty options slice when options.Option is nil")
}

func TestWithIgnores(t *testing.T) {
	tests := []struct {
		name     string
		ignores  []string
		expected int
	}{
		{
			name:     "single ignore",
			ignores:  []string{"test"},
			expected: 1,
		},
		{
			name:     "multiple ignores",
			ignores:  []string{"test1", "test2"},
			expected: 2,
		},
		{
			name:     "empty ignores",
			ignores:  []string{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithIgnores(tt.ignores...)

			// Apply the options to a file
			f := &file{
				ignores: defaultIgnores, // Initialize with default ignores
			}
			f = applyFileOptions(f, option)
			// Only check the newly added ignores, not the default ones
			if tt.expected > 0 {
				// The last 'tt.expected' elements should be our test ignores
				totalIgnores := len(f.ignores)
				startIdx := totalIgnores - tt.expected
				assert.Equal(t, tt.ignores, f.ignores[startIdx:])
			} else {
				// Should only have default ignores
				assert.Equal(t, defaultIgnores, f.ignores)
			}
		})
	}
}

func TestWithFormatter(t *testing.T) {
	t.Log("1. Getting default options")
	var opts []options.Option

	// Define a test formatter
	testFormatter := func(key string, value []byte) (*config.KeyValue, error) {
		t.Logf("Formatter called with key: %s, value: %v", key, value)
		return &config.KeyValue{
			Key:   key,
			Value: value,
		}, nil
	}

	t.Log("2. Setting formatter in options")
	opts = append(opts, WithFormatter(testFormatter))

	// Apply the options to a file
	t.Log("3. Creating file instance")
	f := &file{
		ignores: defaultIgnores, // Initialize with default ignores
	}

	t.Log("4. Applying options to file")
	t.Logf("Number of options found: %d", len(opts))
	for i, opt := range opts {
		t.Logf("Option %d: %T", i, opt)
	}

	f = applyFileOptions(f, opts...)

	t.Log("5. Verifying formatter was set")
	if f.formatter == nil {
		t.Fatal("Formatter is still nil after applyFileOptions")
	}
	assert.NotNil(t, f.formatter)

	t.Log("6. Testing formatter function")
	kv, err := f.formatter("test", []byte("value"))
	if err != nil {
		t.Fatalf("Formatter returned error: %v", err)
	}
	if kv == nil {
		t.Fatal("Formatter returned nil KeyValue")
	}
	assert.Equal(t, "test", kv.Key)
	assert.Equal(t, []byte("value"), kv.Value)
	t.Log("7. Test completed successfully")
}

func TestEmptyOptions(t *testing.T) {
	// Test that FromOptions works with empty options
	opts := optionutil.WithContext(nil)
	f := applyFileOptions(&file{}, opts)
	assert.Empty(t, f, "Expected empty options slice when no options are set")
}

func TestMultipleOptions(t *testing.T) {
	// Test that multiple options can be set and retrieved correctly
	var opts []options.Option

	// Set multiple options
	opts = append(opts, WithIgnores("test1", "test2"))

	// Set a formatter
	formatterCalled := false
	testFormatter := func(key string, value []byte) (*config.KeyValue, error) {
		formatterCalled = true
		return &config.KeyValue{
			Key:   key,
			Value: value,
		}, nil
	}
	opts = append(opts, WithFormatter(testFormatter))

	// Apply the options
	f := &file{
		ignores: defaultIgnores, // Initialize with default ignores
	}
	f = applyFileOptions(f, opts...)

	// Verify both options were applied
	// The ignores should include both default ignores and our test ignores
	assert.GreaterOrEqual(t, len(f.ignores), 2, "Should have at least 2 ignores")
	// Check that our test ignores are at the end of the slice
	testIgnores := f.ignores[len(f.ignores)-2:]
	assert.Equal(t, []string{"test1", "test2"}, testIgnores)
	assert.NotNil(t, f.formatter)

	// Test the formatter
	formatterReturned, err := f.formatter("test", []byte("value"))
	if err != nil {
		t.Fatalf("Formatter returned error: %v", err)
	}
	assert.Equal(t, "test", formatterReturned.Key)
	assert.Equal(t, []byte("value"), formatterReturned.Value)
	assert.True(t, formatterCalled, "Expected formatter to be called")
}
