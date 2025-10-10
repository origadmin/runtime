package http

import (
	"net/http"
	"regexp"

	corsv1 "github.com/origadmin/runtime/api/gen/go/middleware/v1/cors"
	"github.com/rs/cors"
)

// NewCorsHandler creates a Kratos filter for CORS based on the provided configuration.
// It uses the github.com/rs/cors library, which offers a rich set of options
// that map well to our detailed Cors proto definition.
func NewCorsHandler(config *corsv1.Cors) func(http.Handler) http.Handler {
	// If CORS is not configured or disabled, return a no-op filter.
	if config == nil || !config.GetEnabled() {
		return func(h http.Handler) http.Handler {
			return h
		}
	}

	// Map the protobuf configuration to the rs/cors Options struct.
	opts := cors.Options{
		AllowedMethods:   config.GetAllowedMethods(),
		AllowedHeaders:   config.GetAllowedHeaders(),
		ExposedHeaders:   config.GetExposedHeaders(),
		MaxAge:           int(config.GetMaxAge()),
		AllowCredentials: config.GetAllowCredentials(),
		AllowWildcard:    config.GetAllowWildcard(),
		// The proto's `preflight_continue` and `options_passthrough` both map to this concept.
		OptionsPassthrough:   config.GetOptionsPassthrough() || config.GetPreflightContinue(),
		OptionsSuccessStatus: int(config.GetOptionsSuccessStatus()),
		Debug:                config.GetDebug(),
	}

	// Origin control is complex, handle it with priority.
	// 1. Allow any origin (highest priority).
	if config.GetAllowAnyOrigin() {
		opts.AllowedOrigins = []string{"*"}
	} else {
		opts.AllowedOrigins = config.GetAllowedOrigins()
	}

	// 2. Origin patterns (wildcards).
	if len(config.GetAllowedOriginPatterns()) > 0 {
		opts.AllowedOriginPatterns = config.GetAllowedOriginPatterns()
	}

	// 3. Origin regex (provides ultimate flexibility).
	if config.GetAllowOriginRegex() != "" {
		// Silently ignore invalid regex from config to prevent server startup failure.
		if regex, err := regexp.Compile(config.GetAllowOriginRegex()); err == nil {
			opts.AllowOriginFunc = func(origin string) bool {
				return regex.MatchString(origin)
			}
		}
	}

	// Create a new CORS handler with the specified options.
	c := cors.New(opts)

	// The `c.Handler` method returns a function with the signature `func(http.Handler) http.Handler`,
	// which is exactly what a Kratos HTTP filter needs.
	return c.Handler
}
