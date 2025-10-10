package http

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/goexts/generic/configure"
	"github.com/rs/cors"

	corsv1 "github.com/origadmin/runtime/api/gen/go/middleware/v1/cors"
)

// NewCorsHandler creates a Kratos filter for CORS by merging the base configuration from proto
// with advanced, code-based functional options.
func NewCorsHandler(config *corsv1.Cors, codeOpts ...CorsOption) (func(http.Handler) http.Handler, error) {
	// If CORS is not configured or disabled, return a no-op filter.
	if config == nil || !config.GetEnabled() {
		return func(h http.Handler) http.Handler {
			return h
		}, nil
	}

	// 1. Apply base configuration from the proto file.
	opts := cors.Options{
		AllowedMethods:       config.GetAllowedMethods(),
		AllowedHeaders:       config.GetAllowedHeaders(),
		ExposedHeaders:       config.GetExposedHeaders(),
		MaxAge:               int(config.GetMaxAge()),
		AllowCredentials:     config.GetAllowCredentials(),
		OptionsPassthrough:   config.GetOptionsPassthrough() || config.GetPreflightContinue(),
		OptionsSuccessStatus: int(config.GetOptionsSuccessStatus()),
		Debug:                config.GetDebug(),
	}

	configure.Apply(&opts, codeOpts)

	// --- Origin Control Logic (Backward-Compatible) ---
	allOrigins := config.GetAllowedOrigins()
	if len(config.GetAllowedOriginPatterns()) > 0 {
		allOrigins = append(allOrigins, config.GetAllowedOriginPatterns()...)
	}
	opts.AllowedOrigins = allOrigins

	if config.GetAllowAnyOrigin() {
		opts.AllowedOrigins = []string{"*"}
	}

	if config.GetAllowOriginRegex() != "" {
		re, err := regexp.Compile(config.GetAllowOriginRegex())
		if err != nil {
			return nil, fmt.Errorf("invalid CORS allow_origin_regex: %w", err)
		}
		opts.AllowOriginFunc = func(origin string) bool {
			return re.MatchString(origin)
		}
	}

	c := cors.New(opts)

	return c.Handler, nil
}
