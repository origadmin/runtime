package log

import (
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/logger/v1"
)

// NewLogger creates a new Kratos logger based on the provided configuration.
// If the configuration is nil, it returns a default standard output logger as a graceful fallback.
func NewLogger(cfg *loggerv1.Logger) (log.Logger, error) {
	// Graceful Fallback: If no configuration is provided, return a safe default.
	if cfg == nil {
		return log.NewStdLogger(os.Stdout), nil
	}

	// Create a specific logger based on the configuration type.
	// This is where different logging implementations (like zap, fluentd) would be handled.
	switch cfg.GetType() {
	case "default", "std", "": // Handle default cases
		// Here you could add more logic based on cfg, e.g., setting the log level.
		return log.NewStdLogger(os.Stdout), nil
	// case "zap":
	// 	 return newZapLogger(cfg.GetZap())
	default:
		return nil, fmt.Errorf("unsupported logger type: %s", cfg.GetType())
	}
}
