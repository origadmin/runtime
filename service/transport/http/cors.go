package http

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/goexts/generic/configure"
	"github.com/rs/cors"

	corsv1 "github.com/origadmin/runtime/api/gen/go/middleware/v1/cors"
)

// NewCorsHandler creates a Kratos filter for CORS by merging framework defaults,
// the base configuration from proto, and advanced, code-based functional options.
func NewCorsHandler(config *corsv1.Cors, codeOpts ...CorsOption) (func(http.Handler) http.Handler, error) {
	// If CORS is not configured or disabled, return a no-op filter.
	if config == nil || !config.GetEnabled() {
		return func(h http.Handler) http.Handler {
			return h
		}, nil
	}

	// 1. Start with sensible framework defaults.
	opts := cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:       []string{"*"},
		AllowCredentials:     config.GetAllowCredentials(), // Default is false, which is safe.
		ExposedHeaders:       config.GetExposedHeaders(),
		MaxAge:               int(config.GetMaxAge()),
		OptionsPassthrough:   config.GetOptionsPassthrough() || config.GetPreflightContinue(),
		OptionsSuccessStatus: int(config.GetOptionsSuccessStatus()),
		Debug:                config.GetDebug(),
	}

	configure.Apply(&opts, codeOpts)

	// 2. Override defaults with proto configuration if provided.
	if len(config.GetAllowedMethods()) > 0 {
		opts.AllowedMethods = config.GetAllowedMethods()
	}
	if len(config.GetAllowedHeaders()) > 0 {
		opts.AllowedHeaders = config.GetAllowedHeaders()
	}

	// 3. Handle complex origin logic, applying a default if no origin config is set.
	allOrigins := config.GetAllowedOrigins()
	if len(config.GetAllowedOriginPatterns()) > 0 {
		allOrigins = append(allOrigins, config.GetAllowedOriginPatterns()...)
	}

	// If no origins are specified in any way (list, pattern, or regex), apply the default "*".
	if len(allOrigins) == 0 && config.GetAllowOriginRegex() == "" {
		opts.AllowedOrigins = []string{"*"}
	} else {
		opts.AllowedOrigins = allOrigins
	}

	// `allow_any_origin` acts as an explicit override to allow all origins.
	if config.GetAllowAnyOrigin() {
		opts.AllowedOrigins = []string{"*"}
	}

	// Regex provides the ultimate flexibility and overrides string-based origin lists.
	if config.GetAllowOriginRegex() != "" {
		re, err := regexp.Compile(config.GetAllowOriginRegex())
		if err != nil {
			return nil, fmt.Errorf("invalid CORS allow_origin_regex: %w", err)
		}
		// When AllowOriginFunc is set, AllowedOrigins is ignored by the library.
		opts.AllowOriginFunc = func(origin string) bool {
			return re.MatchString(origin)
		}
	}

	c := cors.New(opts)

	return c.Handler, nil
}
