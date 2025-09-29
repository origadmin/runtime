package optionutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/origadmin/runtime/interfaces/options"
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

	var opts []options.Option

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

// --- Test Fixtures ---

type serverConfig struct {
	Host string
	Port int
}

type dbConfig struct {
	DSN string
}

func withHost(host string) options.Option {
	return optionutil.Update(func(c *serverConfig) {
		c.Host = host
	})
}

func withPort(port int) options.Option {
	return optionutil.Update(func(c *serverConfig) {
		c.Port = port
	})
}

// --- Test Cases ---

func TestWithValue(t *testing.T) {
	ctx := optionutil.Empty()
	key1 := optionutil.Key[string]{}
	key2 := optionutil.Key[int]{}

	ctx = optionutil.With(ctx, key1, "hello")
	ctx = optionutil.With(ctx, key2, 123)

	// Test retrieving values
	val1, ok1 := optionutil.Value(ctx, key1)
	assert.True(t, ok1)
	assert.Equal(t, "hello", val1)

	val2, ok2 := optionutil.Value(ctx, key2)
	assert.True(t, ok2)
	assert.Equal(t, 123, val2)

	// Test retrieving non-existent value
	_, ok3 := optionutil.Value(ctx, optionutil.Key[bool]{})
	assert.False(t, ok3)
}

func TestImmutability(t *testing.T) {
	ctx1 := optionutil.Empty()
	key1 := optionutil.Key[string]{}
	ctx2 := optionutil.With(ctx1, key1, "hello")

	// ctx1 should not be changed
	_, ok := optionutil.Value(ctx1, key1)
	assert.False(t, ok)

	// ctx2 should have the value
	val, ok := optionutil.Value(ctx2, key1)
	assert.True(t, ok)
	assert.Equal(t, "hello", val)
}

func TestApply(t *testing.T) {
	// Start with a default config
	cfg := &serverConfig{
		Host: "localhost",
		Port: 8080,
	}

	// Apply options to it
	ctx := optionutil.Apply(cfg,
		withHost("example.com"),
		withPort(9090),
	)

	// Retrieve the configured object from the resulting context
	retrievedCfg, ok := optionutil.FromContext[*serverConfig](ctx)
	assert.True(t, ok)
	assert.Equal(t, "example.com", retrievedCfg.Host)
	assert.Equal(t, 9090, retrievedCfg.Port)
	// Ensure it's the same pointer
	assert.Same(t, cfg, retrievedCfg)
}

func TestApplyNew(t *testing.T) {
	// ApplyNew should create a new instance and configure it
	cfg := optionutil.ApplyNew[serverConfig](
		withHost("new.example.com"),
		withPort(3000),
	)

	assert.NotNil(t, cfg)
	assert.Equal(t, "new.example.com", cfg.Host)
	assert.Equal(t, 3000, cfg.Port)
}

func TestUpdate(t *testing.T) {
	t.Run("Update existing object", func(t *testing.T) {
		cfg := &serverConfig{Port: 80}
		ctx := optionutil.With(optionutil.Empty(), optionutil.Key[*serverConfig]{}, cfg)

		opt := optionutil.Update(func(c *serverConfig) {
			c.Port = 443
		})

		newCtx := opt(ctx)

		retrievedCfg, _ := optionutil.FromContext[*serverConfig](newCtx)
		assert.Equal(t, 443, retrievedCfg.Port)
	})

	t.Run("Update does nothing if type not found", func(t *testing.T) {
		ctx := optionutil.Empty() // Empty context
		opt := withPort(9999)     // An option that targets *serverConfig

		// Apply the option, it should not panic or error
		newCtx := opt(ctx)

		// The context should still not contain the config
		_, ok := optionutil.FromContext[*serverConfig](newCtx)
		assert.False(t, ok)
	})
}

func TestWithContext(t *testing.T) {
	key1 := optionutil.Key[string]{}
	key2 := optionutil.Key[int]{}

	ctx1 := optionutil.With(optionutil.Empty(), key1, "from ctx1")
	ctx2 := optionutil.With(optionutil.Empty(), key2, 999)

	// Create an option that will replace ctx1 with ctx2
	opt := optionutil.WithContext(ctx2)

	// Apply the option
	ctx3 := opt(ctx1)

	// The resulting context should be identical to ctx2
	val2, ok2 := optionutil.Value(ctx3, key2)
	assert.True(t, ok2)
	assert.Equal(t, 999, val2)

	// The value from ctx1 should not exist
	_, ok1 := optionutil.Value(ctx3, key1)
	assert.False(t, ok1)
}

func TestAppendAndSlice(t *testing.T) {
	key := optionutil.Key[[]string]{}
	ctx := optionutil.Empty()

	// Append to a nil context
	ctx = optionutil.Append(ctx, key, "a", "b")
	slice1 := optionutil.Slice(ctx, key)
	assert.Equal(t, []string{"a", "b"}, slice1)

	// Append to an existing slice
	ctx = optionutil.Append(ctx, key, "c")
	slice2 := optionutil.Slice(ctx, key)
	assert.Equal(t, []string{"a", "b", "c"}, slice2)

	// Test that Slice returns a copy
	slice2[0] = "z"
	slice3 := optionutil.Slice(ctx, key)
	assert.Equal(t, "a", slice3[0], "Slice should return a copy, original should not be modified")
}

func TestChainingAndDependencies(t *testing.T) {
	// An option that sets dbConfig
	withDB := func(dsn string) options.Option {
		return func(ctx options.Context) options.Context {
			return optionutil.With(ctx, optionutil.Key[*dbConfig]{}, &dbConfig{DSN: dsn})
		}
	}

	// An option that depends on dbConfig to set the server host
	withHostFromDB := func() options.Option {
		return func(ctx options.Context) options.Context {
			if db, ok := optionutil.Value(ctx, optionutil.Key[*dbConfig]{}); ok {
				// Now, update the serverConfig based on the dbConfig
				updateOpt := optionutil.Update(func(sc *serverConfig) {
					sc.Host = db.DSN // e.g., use DSN as host
				})
				return updateOpt(ctx)
			}
			return ctx
		}
	}

	// Apply options in order. withHostFromDB depends on withDB.
	finalCtx := optionutil.Apply(&serverConfig{},
		withDB("my-database-dsn"),
		withHostFromDB(),
	)

	finalCfg, ok := optionutil.FromContext[*serverConfig](finalCtx)
	assert.True(t, ok)
	assert.Equal(t, "my-database-dsn", finalCfg.Host)
}
