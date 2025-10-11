package optionutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/optionutil"
)

// --- Test Fixtures ---

type serverConfig struct {
	Host string
	Port int
}

type dbConfig struct {
	DSN string
}

var (
	stringKey = optionutil.Key[string]{}
	intKey    = optionutil.Key[int]{}
	sliceKey  = optionutil.Key[[]int]{}
)

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

// --- Core API Tests ---

func TestNew(t *testing.T) {
	// New should create a new instance, configure it, and return both context and config.
	ctx, cfg := optionutil.New[serverConfig](
		withHost("new.example.com"),
		withPort(3000),
	)

	assert.NotNil(t, cfg)
	assert.Equal(t, "new.example.com", cfg.Host)
	assert.Equal(t, 3000, cfg.Port)

	// Verify the returned context also contains the config
	retrievedCfg, ok := optionutil.Extract[*serverConfig](ctx)
	assert.True(t, ok)
	assert.Same(t, cfg, retrievedCfg, "Context should contain the same instance that was returned")
}

func TestNewT(t *testing.T) {
	// NewT should create a new instance, configure it, and return only the config.
	cfg := optionutil.NewT[serverConfig](
		withHost("newt.example.com"),
		withPort(3001),
	)

	assert.NotNil(t, cfg)
	assert.Equal(t, "newt.example.com", cfg.Host)
	assert.Equal(t, 3001, cfg.Port)
}

func TestNewContext(t *testing.T) {
	// NewContext should create a new instance, configure it, and return only the context.
	ctx := optionutil.NewContext[serverConfig](
		withHost("newcontext.example.com"),
		withPort(3002),
	)

	assert.NotNil(t, ctx)

	// Verify the returned context contains the configured value
	retrievedCfg, ok := optionutil.Extract[*serverConfig](ctx)
	assert.True(t, ok)
	assert.NotNil(t, retrievedCfg)
	assert.Equal(t, "newcontext.example.com", retrievedCfg.Host)
	assert.Equal(t, 3002, retrievedCfg.Port)
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

	// The original object should be modified
	assert.Equal(t, "example.com", cfg.Host)
	assert.Equal(t, 9090, cfg.Port)

	// Retrieve the configured object from the resulting context
	retrievedCfg, ok := optionutil.Extract[*serverConfig](ctx)
	assert.True(t, ok)
	assert.Equal(t, "example.com", retrievedCfg.Host)
	assert.Equal(t, 9090, retrievedCfg.Port)
	// Ensure it's the same pointer
	assert.Same(t, cfg, retrievedCfg)
}

// --- Update and Implicit Key Tests ---

func TestUpdate(t *testing.T) {
	t.Run("Update existing object", func(t *testing.T) {
		cfg := &serverConfig{Port: 80}
		ctx := optionutil.WithValue(optionutil.Empty(), optionutil.Key[*serverConfig]{}, cfg)

		opt := optionutil.Update(func(c *serverConfig) {
			c.Port = 443
		})

		newCtx := opt(ctx)

		retrievedCfg, _ := optionutil.Extract[*serverConfig](newCtx)
		assert.Equal(t, 443, retrievedCfg.Port)
	})

	t.Run("Update with multiple updaters", func(t *testing.T) {
		cfg := &serverConfig{Host: "initial", Port: 80}
		ctx := optionutil.WithValue(optionutil.Empty(), optionutil.Key[*serverConfig]{}, cfg)

		opt := optionutil.Update(
			func(c *serverConfig) { c.Host = "updated" },
			func(c *serverConfig) { c.Port = 443 },
		)
		newCtx := opt(ctx)
		retrievedCfg, _ := optionutil.Extract[*serverConfig](newCtx)
		assert.Equal(t, "updated", retrievedCfg.Host)
		assert.Equal(t, 443, retrievedCfg.Port)
	})

	t.Run("Update does nothing if type not found", func(t *testing.T) {
		ctx := optionutil.Empty() // Empty context
		opt := withPort(9999)     // An option that targets *serverConfig

		// Apply the option, it should not panic or error
		newCtx := opt(ctx)

		// The context should still not contain the config
		_, ok := optionutil.Extract[*serverConfig](newCtx)
		assert.False(t, ok)
	})
}

func TestImplicitFunctions(t *testing.T) {
	t.Run("Extract and ExtractOr", func(t *testing.T) {
		ctx := optionutil.NewContext[serverConfig](withHost("host"))

		// Value exists
		val, ok := optionutil.Extract[*serverConfig](ctx)
		assert.True(t, ok)
		assert.Equal(t, "host", val.Host)

		// Or returns existing
		valOr := optionutil.ExtractOr(ctx, &serverConfig{Host: "default"})
		assert.Equal(t, "host", valOr.Host)

		// Value does not exist, Or returns default
		dbVal := optionutil.ExtractOr(ctx, &dbConfig{DSN: "default"})
		assert.Equal(t, "default", dbVal.DSN)
	})

	t.Run("AppendValues and ExtractSlice", func(t *testing.T) {
		ctx := optionutil.Empty()
		ctx = optionutil.AppendValues(ctx, 1, 2)
		ctx = optionutil.AppendValues(ctx, 3)

		slice := optionutil.ExtractSlice[int](ctx)
		assert.Equal(t, []int{1, 2, 3}, slice)
	})

	t.Run("ExtractSliceOr", func(t *testing.T) {
		ctx := optionutil.Empty()

		// Slice exists
		ctx = optionutil.AppendValues(ctx, 10)
		slice1 := optionutil.ExtractSliceOr(ctx, []int{99})
		assert.Equal(t, []int{10}, slice1)

		// Slice does not exist
		type dummy struct{}
		slice2 := optionutil.ExtractSliceOr(ctx, []dummy{{}})
		assert.Equal(t, []dummy{{}}, slice2)
	})
}

// --- Explicit Key Tests ---

func TestWithValue(t *testing.T) {
	ctx := optionutil.Empty()

	ctx = optionutil.WithValue(ctx, stringKey, "hello")
	ctx = optionutil.WithValue(ctx, intKey, 123)

	// Test retrieving values
	val1, ok1 := optionutil.Value(ctx, stringKey)
	assert.True(t, ok1)
	assert.Equal(t, "hello", val1)

	val2, ok2 := optionutil.Value(ctx, intKey)
	assert.True(t, ok2)
	assert.Equal(t, 123, val2)

	// Test retrieving non-existent value
	_, ok3 := optionutil.Value(ctx, optionutil.Key[bool]{})
	assert.False(t, ok3)
}

func TestAppendAndSliceValue(t *testing.T) {
	ctx := optionutil.Empty()

	// Append to a nil context
	ctx = optionutil.Append(ctx, sliceKey, 1, 2)
	slice1 := optionutil.SliceValue(ctx, sliceKey)
	assert.Equal(t, []int{1, 2}, slice1)

	// Append to an existing slice
	ctx = optionutil.Append(ctx, sliceKey, 3)
	slice2 := optionutil.SliceValue(ctx, sliceKey)
	assert.Equal(t, []int{1, 2, 3}, slice2)
}

// --- Utility Option Tests ---

func TestWithCond(t *testing.T) {
	opt := func(ctx options.Context) options.Context {
		return optionutil.WithValue(ctx, stringKey, "applied")
	}

	t.Run("Condition is true", func(t *testing.T) {
		conditionalOpt := optionutil.WithCond(true, opt)
		ctx := conditionalOpt(optionutil.Empty())
		val, ok := optionutil.Value(ctx, stringKey)
		assert.True(t, ok)
		assert.Equal(t, "applied", val)
	})

	t.Run("Condition is false", func(t *testing.T) {
		conditionalOpt := optionutil.WithCond(false, opt)
		ctx := conditionalOpt(optionutil.Empty())
		_, ok := optionutil.Value(ctx, stringKey)
		assert.False(t, ok)
	})
}

func TestWithGroup(t *testing.T) {
	opt1 := func(ctx options.Context) options.Context {
		return optionutil.WithValue(ctx, stringKey, "from group 1")
	}
	opt2 := func(ctx options.Context) options.Context {
		return optionutil.WithValue(ctx, intKey, 456)
	}

	groupedOpt := optionutil.WithGroup(opt1, opt2)
	ctx := groupedOpt(optionutil.Empty())

	val1, ok1 := optionutil.Value(ctx, stringKey)
	assert.True(t, ok1)
	assert.Equal(t, "from group 1", val1)

	val2, ok2 := optionutil.Value(ctx, intKey)
	assert.True(t, ok2)
	assert.Equal(t, 456, val2)
}

// --- Other Tests ---

func TestImmutability(t *testing.T) {
	ctx1 := optionutil.Empty()
	ctx2 := optionutil.WithValue(ctx1, stringKey, "hello")

	// ctx1 should not be changed
	_, ok := optionutil.Value(ctx1, stringKey)
	assert.False(t, ok)

	// ctx2 should have the value
	val, ok := optionutil.Value(ctx2, stringKey)
	assert.True(t, ok)
	assert.Equal(t, "hello", val)
}

func TestWithNilKeyPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	var ctx options.Context = optionutil.Empty()
	// This should panic because the key is nil
	ctx.With(nil, "some value")
}
