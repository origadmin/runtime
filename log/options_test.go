package log

import (
	"testing"

	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"

	"github.com/origadmin/runtime/optionutil"
)

// MockLogger implements the Logger interface for testing purposes.
type MockLogger struct {
	name string
}

func (m *MockLogger) Log(level kratoslog.Level, keyvals ...interface{}) error {
	// For testing, we can just print or store the log.
	// fmt.Printf("[%s] %s: %v\n", m.name, level.String(), keyvals)
	return nil
}

func TestFromOptions(t *testing.T) {
	t.Run("should return DefaultLogger when no options are provided", func(t *testing.T) {
		logger := FromOptions()
		assert.Equal(t, DefaultLogger, logger)
	})

	t.Run("should return the logger provided by WithLogger option", func(t *testing.T) {
		mockLogger := &MockLogger{name: "test-logger-from-options"}
		logger := FromOptions(WithLogger(mockLogger))
		assert.Equal(t, mockLogger, logger)
	})

	t.Run("should return the last logger when multiple WithLogger options are provided", func(t *testing.T) {
		mockLogger1 := &MockLogger{name: "test-logger-1"}
		mockLogger2 := &MockLogger{name: "test-logger-2"}
		logger := FromOptions(WithLogger(mockLogger1), WithLogger(mockLogger2))
		assert.Equal(t, mockLogger2, logger)
	})
}

func TestFromContext(t *testing.T) {
	t.Run("should return DefaultLogger when context is empty", func(t *testing.T) {
		emptyOptCtx := optionutil.Empty()
		logger := FromContext(emptyOptCtx)
		assert.Equal(t, DefaultLogger, logger)
	})

	t.Run("should return the logger stored in the context", func(t *testing.T) {
		mockLogger := &MockLogger{name: "test-logger-from-context"}

		// Manually create and store loggerContext in interfaces.Context
		lc := &loggerContext{Logger: mockLogger}
		optCtx := optionutil.WithValue(optionutil.Empty(), optionutil.Key[*loggerContext]{}, lc)

		logger := FromContext(optCtx)
		assert.Equal(t, mockLogger, logger)
	})

	t.Run("should return DefaultLogger if loggerContext is in context but Logger field is nil", func(t *testing.T) {
		lc := &loggerContext{Logger: nil} // Logger field is nil
		optCtx := optionutil.WithValue(optionutil.Empty(), optionutil.Key[*loggerContext]{}, lc)

		logger := FromContext(optCtx)
		assert.Equal(t, DefaultLogger, logger)
	})
}

func TestWithLogger(t *testing.T) {
	t.Run("WithLogger should correctly update logger in context", func(t *testing.T) {
		mockLogger1 := &MockLogger{name: "initial-logger"}
		mockLogger2 := &MockLogger{name: "updated-logger"}

		// Start with a context containing mockLogger1
		lc := &loggerContext{Logger: mockLogger1}
		optCtx := optionutil.WithValue(optionutil.Empty(), optionutil.Key[*loggerContext]{}, lc)

		// Apply the WithLogger option to update the logger to mockLogger2
		option := WithLogger(mockLogger2)
		option(optCtx) // Apply the option to the context

		// Retrieve the logger from the updated context
		retrievedLogger := FromContext(optCtx)
		assert.Equal(t, mockLogger2, retrievedLogger)
	})

	t.Run("WithLogger should not panic if loggerContext is not initially in context", func(t *testing.T) {
		mockLogger := &MockLogger{name: "new-logger"}
		emptyOptCtx := optionutil.Empty()

		// Applying WithLogger to an empty context (where loggerContext is not yet present)
		// optionutil.Update will not create it, so Extract should still return DefaultLogger
		option := WithLogger(mockLogger)
		option(emptyOptCtx)

		retrievedLogger := FromContext(emptyOptCtx)
		assert.Equal(t, DefaultLogger, retrievedLogger) // Expect DefaultLogger as Update doesn't add if not present
	})
}
