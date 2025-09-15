package interfaces_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/origadmin/runtime/interfaces"
)

func TestOptionSet(t *testing.T) {
	// Test basic value storage and retrieval
	t.Run("basic value storage", func(t *testing.T) {
		key1 := "testKey1"
		value1 := "testValue1"
		key2 := "testKey2"
		value2 := 42

		// Test setting and getting a single value
		opt := interfaces.DefaultOptions().WithOption(key1, value1)
		assert.Equal(t, value1, opt.Option(key1))

		// Test chaining
		opt = opt.WithOption(key2, value2)
		assert.Equal(t, value1, opt.Option(key1))
		assert.Equal(t, value2, opt.Option(key2))

		// Test non-existent key
		assert.Nil(t, opt.Option("nonExistentKey"))
	})

	// Test nil safety
	t.Run("nil safety", func(t *testing.T) {
		var opt *interfaces.OptionSet = nil
		assert.Nil(t, opt.Option("anyKey"))

		// Test WithOption on nil
		opt1 := (*interfaces.OptionSet)(nil).WithOption("key", "value")
		assert.NotNil(t, opt1)
		assert.Equal(t, "value", opt1.Option("key"))
	})

	// Test type safety
	t.Run("type safety", func(t *testing.T) {
		key := "testKey"
		opt := interfaces.DefaultOptions().
			WithOption(key, "string value").
			WithOption(123, 456)

		// Should return the correct type
		assert.IsType(t, "", opt.Option(key))
		assert.IsType(t, 0, opt.Option(123))
	})
}

func TestDefaultOptions(t *testing.T) {
	opt := interfaces.DefaultOptions()
	assert.NotNil(t, opt)
	assert.Nil(t, opt.Option("nonExistent"))
}
