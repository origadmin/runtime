package optionutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/optionutil"
)

type testValue struct {
	Name  string
	Value int
}

var (
	stringKey = optionutil.Key[string]{}
	intKey    = optionutil.Key[int]{}
)

func TestWithOption(t *testing.T) {
	value1 := "test value"

	// Test basic WithOption
	opt := optionutil.Empty()
	opt = optionutil.With(opt, stringKey, value1)
	assert.NotNil(t, opt)

	// Test Value
	got, ok := optionutil.Value(opt, stringKey)
	assert.True(t, ok)
	assert.Equal(t, value1, got)

	// Test with existing emptyContext
	value2 := 42
	opt = optionutil.With(opt, intKey, value2)

	// Verify both values exist
	got1, ok1 := optionutil.Value(opt, stringKey)
	got2, ok2 := optionutil.Value(opt, intKey)

	assert.True(t, ok1)
	assert.Equal(t, value1, got1)
	assert.True(t, ok2)
	assert.Equal(t, value2, got2)
}

var sliceKey = optionutil.Key[[]int]{}

func TestWithSliceOption(t *testing.T) {
	key := sliceKey

	// Test initial slice
	opt := optionutil.Empty()
	opt = optionutil.Append(opt, key, 1, 2, 3)
	slice1 := optionutil.Slice(opt, key)
	assert.Equal(t, []int{1, 2, 3}, slice1)

	// Test appending to slice
	opt = optionutil.Append(opt, key, 4, 5)
	slice2 := optionutil.Slice(opt, key)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, slice2)

	// Test with struct values
	type item struct{ Name string }
	var itemsKey = optionutil.Key[[]item]{}
	opt = optionutil.Append(opt, itemsKey, item{"a"}, item{"b"})
	items := optionutil.Slice(opt, itemsKey)
	assert.Len(t, items, 2)
	assert.Equal(t, "a", items[0].Name)
}

var nonExistentKey = optionutil.Key[string]{}

func TestGetOption_NonExistent(t *testing.T) {
	key := nonExistentKey
	opt := optionutil.Empty()

	// Test with non-existent key
	_, ok := optionutil.Value(opt, key)
	assert.False(t, ok)

	// Test with nil emptyContext
	_, ok = optionutil.Value(nil, key)
	assert.False(t, ok)
}

var sliceNonExistentKey = optionutil.Key[[]int]{}

func TestGetSliceOption_NonExistent(t *testing.T) {
	key := sliceNonExistentKey
	opt := optionutil.Empty()

	// Test with non-existent key
	slice := optionutil.Slice(opt, key)
	nilSlice := []int(nil)
	assert.Equal(t, nilSlice, slice)

	// Test with nil emptyContext
	slice = optionutil.Slice(nil, key)
	assert.Equal(t, nilSlice, slice)
}

func TestUpdateAndApply(t *testing.T) {
	type configA struct {
		Name  string
		Value int
	}
	type configB struct {
		Value int
	}

	var opts []interfaces.Option

	// Create an option to update configA
	opts = append(opts, optionutil.Update(func(cfg *configA) {
		cfg.Name = "updated A"
	}))

	opts = append(opts, optionutil.Update(func(cfg *configA) {
		cfg.Value = 42
	}))

	// Create an option to update configB
	opts = append(opts, optionutil.Update(func(cfg *configB) {
		cfg.Value = 123
	}))

	// --- Test with configA ---
	a := &configA{Name: "initial A"}
	optionutil.Apply(a, opts...)

	// The updater for configA should have run
	assert.Equal(t, "updated A", a.Name)
	assert.Equal(t, 42, a.Value)
	// --- Test with configB ---
	b := &configB{Value: 0}
	optionutil.Apply(b, opts...)

	// The updater for configB should have run
	assert.Equal(t, 123, b.Value)

	// --- Test with an unrelated type ---
	c := &struct{}{} // No updater for this type
	optionutil.Apply(c, opts...)
	// No changes should occur, and it should not panic.
}
