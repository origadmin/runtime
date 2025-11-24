package optionutil

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/origadmin/runtime/interfaces/options"
)

// --- Test Fixtures ---

type serverConfig struct {
	Host string
	Port int
}

func withHost(host string) options.Option {
	return Update(func(c *serverConfig) {
		c.Host = host
	})
}

func withPort(port int) options.Option {
	return Update(func(c *serverConfig) {
		c.Port = port
	})
}

// --- Core API Tests ---

func TestNew(t *testing.T) {
	ctx, cfg := New[serverConfig](
		withHost("new.example.com"),
		withPort(3000),
	)

	assert.NotNil(t, cfg)
	assert.Equal(t, "new.example.com", cfg.Host)
	assert.Equal(t, 3000, cfg.Port)

	retrievedCfg, ok := Value[*serverConfig](ctx)
	assert.True(t, ok)
	assert.Same(t, cfg, retrievedCfg)
}

func TestNewT(t *testing.T) {
	cfg := NewT[serverConfig](
		withHost("newt.example.com"),
		withPort(3001),
	)

	assert.NotNil(t, cfg)
	assert.Equal(t, "newt.example.com", cfg.Host)
	assert.Equal(t, 3001, cfg.Port)
}

func TestNewContext(t *testing.T) {
	ctx := NewContext[serverConfig](
		withHost("newcontext.example.com"),
		withPort(3002),
	)

	assert.NotNil(t, ctx)

	retrievedCfg, ok := Value[*serverConfig](ctx)
	assert.True(t, ok)
	assert.NotNil(t, retrievedCfg)
	assert.Equal(t, "newcontext.example.com", retrievedCfg.Host)
	assert.Equal(t, 3002, retrievedCfg.Port)
}

func TestApply(t *testing.T) {
	cfg := &serverConfig{Host: "localhost", Port: 8080}

	ctx := Apply(cfg,
		withHost("example.com"),
		withPort(9090),
	)

	assert.Equal(t, "example.com", cfg.Host)
	assert.Equal(t, 9090, cfg.Port)

	retrievedCfg, ok := Value[*serverConfig](ctx)
	assert.True(t, ok)
	assert.Same(t, cfg, retrievedCfg)
}

// --- Implicit Key Data Access ---

func TestValue(t *testing.T) {
	ctx := Empty()
	ctx = WithValue(ctx, "hello")
	ctx = WithValue(ctx, 123)

	valStr, okStr := Value[string](ctx)
	assert.True(t, okStr)
	assert.Equal(t, "hello", valStr)

	valInt, okInt := Value[int](ctx)
	assert.True(t, okInt)
	assert.Equal(t, 123, valInt)
}

func TestValueOr(t *testing.T) {
	ctx := WithValue(Empty(), "actual")

	val1 := ValueOr(ctx, "default")
	assert.Equal(t, "actual", val1)

	val2 := ValueOr(ctx, 999)
	assert.Equal(t, 999, val2)
}

func TestSliceAndAppend(t *testing.T) {
	ctx := Empty()
	ctx = Append(ctx, "a", "b")
	ctx = Append(ctx, "c")

	slice := Slice[string](ctx)
	assert.Equal(t, []string{"a", "b", "c"}, slice)
}

func TestSliceOr(t *testing.T) {
	ctx := Append(Empty(), "a")
	defaultValue := []string{"default"}

	val1 := SliceOr(ctx, defaultValue)
	assert.Equal(t, []string{"a"}, val1)

	val2 := SliceOr(ctx, []int{99})
	assert.Equal(t, []int{99}, val2)
}

// --- Explicit Key (ByKey) Data Access ---

func TestByKeyFunctions(t *testing.T) {
	ctx := Empty()

	// Test WithValue
	ctx = WithValue(ctx, "value A")
	valA, okA := Value[string](ctx)
	assert.True(t, okA)
	assert.Equal(t, "value A", valA)

	ctx = WithValue(ctx, "value B")
	valB, okB := Value[string](ctx)
	assert.True(t, okB)
	assert.Equal(t, "value B", valB)

	// Test appendValues
	ctx = Append(ctx, 1, 2)
	ctx = Append(ctx, 3)

	// Test slice
	slice := Slice[int](ctx)
	assert.Equal(t, []int{1, 2, 3}, slice)
}

// --- Update Function Test ---

func TestUpdate(t *testing.T) {
	// This test now implicitly tests that Update calls UpdateByKey correctly
	opt := Update(
		func(c *serverConfig) { c.Host = "updated" },
		func(c *serverConfig) { c.Port = 443 },
	)

	ctx, cfg := New[serverConfig](opt)

	assert.Equal(t, "updated", cfg.Host)
	assert.Equal(t, 443, cfg.Port)

	retrieved, _ := Value[*serverConfig](ctx)
	assert.Same(t, cfg, retrieved)
}

// --- Utility Option Tests ---

func TestWithCond(t *testing.T) {
	opt := func(ctx options.Context) options.Context {
		return WithValue(ctx, "applied")
	}

	t.Run("Condition is true", func(t *testing.T) {
		conditionalOpt := WithCond(true, opt)
		ctx := conditionalOpt(Empty())
		val, ok := Value[string](ctx)
		assert.True(t, ok)
		assert.Equal(t, "applied", val)
	})

	t.Run("Condition is false", func(t *testing.T) {
		conditionalOpt := WithCond(false, opt)
		ctx := conditionalOpt(Empty())
		_, ok := Value[string](ctx)
		assert.False(t, ok)
	})
}

func TestWithGroup(t *testing.T) {
	opt1 := func(ctx options.Context) options.Context {
		return WithValue(ctx, "from group 1")
	}
	opt2 := func(ctx options.Context) options.Context {
		return WithValue(ctx, 456)
	}

	groupedOpt := WithGroup(opt1, opt2)
	ctx := groupedOpt(Empty())

	val1, ok1 := Value[string](ctx)
	assert.True(t, ok1)
	assert.Equal(t, "from group 1", val1)

	val2, ok2 := Value[int](ctx)
	assert.True(t, ok2)
	assert.Equal(t, 456, val2)
}

// TestOptionCarriesContextData verifies that data stored in the context via an options.Option
// can be retrieved after applying options with NewContext.
func TestOptionCarriesContextData(t *testing.T) {
	type customData struct {
		ID   string
		Name string
	}

	expectedData := customData{ID: "test-id", Name: "Test Name"}

	// Define an options.Option that puts customData into the context
	withCustomContextData := func(data customData) options.Option {
		return func(ctx options.Context) options.Context {
			return WithValue(ctx, data)
		}
	}

	// Apply the option using NewContext
	ctx := NewContext[serverConfig](
		withHost("context.data.example.com"),
		withCustomContextData(expectedData),
	)

	assert.NotNil(t, ctx)

	// Retrieve the custom data from the context
	retrievedData, ok := Value[customData](ctx)
	assert.True(t, ok, "Expected customData to be found in context")
	assert.Equal(t, expectedData, retrievedData, "Retrieved customData should match expectedData")

	// Verify other options still work
	retrievedCfg, ok := Value[*serverConfig](ctx)
	assert.True(t, ok)
	assert.NotNil(t, retrievedCfg)
	assert.Equal(t, "context.data.example.com", retrievedCfg.Host)
}

// --- Other Tests ---

func TestImmutability(t *testing.T) {
	ctx1 := Empty()
	ctx2 := WithValue(ctx1, "hello")

	_, ok := Value[string](ctx1)
	assert.False(t, ok)

	val, ok := Value[string](ctx2)
	assert.True(t, ok)
	assert.Equal(t, "hello", val)
}

func TestWithNilKeyPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	var ctx options.Context = Empty()
	ctx.With(nil, "some value")
}
