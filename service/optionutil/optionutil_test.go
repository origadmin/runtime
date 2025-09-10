package optionutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/service/optionutil"
)

type testValue struct {
	Name  string
	Value int
}

var (
	stringKey = optionutil.OptionKey[string]{}
	intKey    = optionutil.OptionKey[int]{}
)

func TestWithOption(t *testing.T) {
	value1 := "test value"

	// Test basic WithOption
	opt := optionutil.WithOption(nil, stringKey, value1)
	assert.NotNil(t, opt)

	// Test GetOption
	got, ok := optionutil.GetOption(opt, stringKey)
	assert.True(t, ok)
	assert.Equal(t, value1, got)

	// Test with existing options
	value2 := 42
	opt = optionutil.WithOption(opt, intKey, value2)

	// Verify both values exist
	got1, ok1 := optionutil.GetOption(opt, stringKey)
	got2, ok2 := optionutil.GetOption(opt, intKey)

	assert.True(t, ok1)
	assert.Equal(t, value1, got1)
	assert.True(t, ok2)
	assert.Equal(t, value2, got2)
}

var sliceKey = optionutil.OptionKey[[]int]{}

func TestWithSliceOption(t *testing.T) {
	key := sliceKey

	// Test initial slice
	opt := optionutil.WithSliceOption(nil, key, 1, 2, 3)
	slice1 := optionutil.GetSliceOption[int](opt, key)
	assert.Equal(t, []int{1, 2, 3}, slice1)

	// Test appending to slice
	opt = optionutil.WithSliceOption(opt, key, 4, 5)
	slice2 := optionutil.GetSliceOption[int](opt, key)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, slice2)

	// Test with struct values
	type item struct{ Name string }
	var itemsKey = optionutil.OptionKey[[]item]{}
	opt = optionutil.WithSliceOption(opt, itemsKey, item{"a"}, item{"b"})
	items := optionutil.GetSliceOption[item](opt, itemsKey)
	assert.Len(t, items, 2)
	assert.Equal(t, "a", items[0].Name)
}

var nonExistentKey = optionutil.OptionKey[string]{}

func TestGetOption_NonExistent(t *testing.T) {
	key := nonExistentKey
	opt := interfaces.DefaultOptions()

	// Test with non-existent key
	_, ok := optionutil.GetOption(opt, key)
	assert.False(t, ok)

	// Test with nil options
	_, ok = optionutil.GetOption(nil, key)
	assert.False(t, ok)
}

var sliceNonExistentKey = optionutil.OptionKey[[]int]{}

func TestGetSliceOption_NonExistent(t *testing.T) {
	key := sliceNonExistentKey
	opt := interfaces.DefaultOptions()

	// Test with non-existent key
	slice := optionutil.GetSliceOption[int](opt, key)
	nilSlice := []int(nil)
	assert.Equal(t, nilSlice, slice)

	// Test with nil options
	slice = optionutil.GetSliceOption[int](nil, key)
	assert.Equal(t, nilSlice, slice)
}
