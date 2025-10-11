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
	ctx, cfg := optionutil.New[serverConfig](
		withHost("new.example.com"),
		withPort(3000),
	)

	assert.NotNil(t, cfg)
	assert.Equal(t, "new.example.com", cfg.Host)
	assert.Equal(t, 3000, cfg.Port)

	retrievedCfg, ok := optionutil.Value[*serverConfig](ctx)
	assert.True(t, ok)
	assert.Same(t, cfg, retrievedCfg)
}

func TestNewT(t *testing.T) {
	cfg := optionutil.NewT[serverConfig](
		withHost("newt.example.com"),
		withPort(3001),
	)

	assert.NotNil(t, cfg)
	assert.Equal(t, "newt.example.com", cfg.Host)
	assert.Equal(t, 3001, cfg.Port)
}

func TestNewContext(t *testing.T) {
	ctx := optionutil.NewContext[serverConfig](
		withHost("newcontext.example.com"),
		withPort(3002),
	)

	assert.NotNil(t, ctx)

	retrievedCfg, ok := optionutil.Value[*serverConfig](ctx)
	assert.True(t, ok)
	assert.NotNil(t, retrievedCfg)
	assert.Equal(t, "newcontext.example.com", retrievedCfg.Host)
	assert.Equal(t, 3002, retrievedCfg.Port)
}

func TestApply(t *testing.T) {
	cfg := &serverConfig{Host: "localhost", Port: 8080}

	ctx := optionutil.Apply(cfg,
		withHost("example.com"),
		withPort(9090),
	)

	assert.Equal(t, "example.com", cfg.Host)
	assert.Equal(t, 9090, cfg.Port)

	retrievedCfg, ok := optionutil.Value[*serverConfig](ctx)
	assert.True(t, ok)
	assert.Same(t, cfg, retrievedCfg)
}

// --- Implicit Key Data Access ---

func TestValue(t *testing.T) {
	ctx := optionutil.Empty()
	ctx = optionutil.WithValue(ctx, "hello")
	ctx = optionutil.WithValue(ctx, 123)

	valStr, okStr := optionutil.Value[string](ctx)
	assert.True(t, okStr)
	assert.Equal(t, "hello", valStr)

	valInt, okInt := optionutil.Value[int](ctx)
	assert.True(t, okInt)
	assert.Equal(t, 123, valInt)
}

func TestValueOr(t *testing.T) {
	ctx := optionutil.WithValue(optionutil.Empty(), "actual")

	val1 := optionutil.ValueOr(ctx, "default")
	assert.Equal(t, "actual", val1)

	val2 := optionutil.ValueOr(ctx, 999)
	assert.Equal(t, 999, val2)
}

func TestSliceAndAppend(t *testing.T) {
	ctx := optionutil.Empty()
	ctx = optionutil.Append(ctx, "a", "b")
	ctx = optionutil.Append(ctx, "c")

	slice := optionutil.Slice[string](ctx)
	assert.Equal(t, []string{"a", "b", "c"}, slice)
}

func TestSliceOr(t *testing.T) {
	ctx := optionutil.Append(optionutil.Empty(), "a")
	defaultValue := []string{"default"}

	val1 := optionutil.SliceOr(ctx, defaultValue)
	assert.Equal(t, []string{"a"}, val1)

	val2 := optionutil.SliceOr(ctx, []int{99})
	assert.Equal(t, []int{99}, val2)
}

// --- Explicit Key (ByKey) Data Access ---

func TestByKeyFunctions(t *testing.T) {
	ctx := optionutil.Empty()

	// Test WithValue
	ctx = optionutil.WithValue(ctx, "value A")
	valA, okA := optionutil.Value[string](ctx)
	assert.True(t, okA)
	assert.Equal(t, "value A", valA)

	ctx = optionutil.WithValue(ctx, "value B")
	valB, okB := optionutil.Value[string](ctx)
	assert.True(t, okB)
	assert.Equal(t, "value B", valB)

	// Test appendValues
	ctx = optionutil.Append(ctx, 1, 2)
	ctx = optionutil.Append(ctx, 3)

	// Test slice
	slice := optionutil.Slice[int](ctx)
	assert.Equal(t, []int{1, 2, 3}, slice)
}

// --- Update Function Test ---

func TestUpdate(t *testing.T) {
	// This test now implicitly tests that Update calls UpdateByKey correctly
	opt := optionutil.Update(
		func(c *serverConfig) { c.Host = "updated" },
		func(c *serverConfig) { c.Port = 443 },
	)

	ctx, cfg := optionutil.New[serverConfig](opt)

	assert.Equal(t, "updated", cfg.Host)
	assert.Equal(t, 443, cfg.Port)

	retrieved, _ := optionutil.Value[*serverConfig](ctx)
	assert.Same(t, cfg, retrieved)
}

// --- Utility Option Tests ---

func TestWithCond(t *testing.T) {
	opt := func(ctx options.Context) options.Context {
		return optionutil.WithValue(ctx, "applied")
	}

	t.Run("Condition is true", func(t *testing.T) {
		conditionalOpt := optionutil.WithCond(true, opt)
		ctx := conditionalOpt(optionutil.Empty())
		val, ok := optionutil.Value[string](ctx)
		assert.True(t, ok)
		assert.Equal(t, "applied", val)
	})

	t.Run("Condition is false", func(t *testing.T) {
		conditionalOpt := optionutil.WithCond(false, opt)
		ctx := conditionalOpt(optionutil.Empty())
		_, ok := optionutil.Value[string](ctx)
		assert.False(t, ok)
	})
}

func TestWithGroup(t *testing.T) {
	opt1 := func(ctx options.Context) options.Context {
		return optionutil.WithValue(ctx, "from group 1")
	}
	opt2 := func(ctx options.Context) options.Context {
		return optionutil.WithValue(ctx, 456)
	}

	groupedOpt := optionutil.WithGroup(opt1, opt2)
	ctx := groupedOpt(optionutil.Empty())

	val1, ok1 := optionutil.Value[string](ctx)
	assert.True(t, ok1)
	assert.Equal(t, "from group 1", val1)

	val2, ok2 := optionutil.Value[int](ctx)
	assert.True(t, ok2)
	assert.Equal(t, 456, val2)
}

// --- Other Tests ---

func TestImmutability(t *testing.T) {
	ctx1 := optionutil.Empty()
	ctx2 := optionutil.WithValue(ctx1, "hello")

	_, ok := optionutil.Value[string](ctx1)
	assert.False(t, ok)

	val, ok := optionutil.Value[string](ctx2)
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
	ctx.With(nil, "some value")
}
