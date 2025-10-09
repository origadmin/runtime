package cors

import (
	"net/http"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"

	corsv1 "github.com/origadmin/runtime/api/gen/go/middleware/v1/cors"
)

// NewGorillaCors creates a new CORS handler for Gorilla/Mux or standard HTTP servers
func NewGorillaCors(cfg *corsv1.Cors) func(http.Handler) http.Handler {
	if cfg == nil || !cfg.Enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	options := []handlers.CORSOption{
		handlers.AllowedOrigins(cfg.AllowOrigins),
		handlers.AllowedMethods(cfg.AllowMethods),
		handlers.AllowedHeaders(cfg.AllowHeaders),
		handlers.ExposedHeaders(cfg.ExposeHeaders),
	}

	if cfg.AllowCredentials {
		options = append(options, handlers.AllowCredentials())
	}

	if cfg.MaxAge > 0 {
		options = append(options, handlers.MaxAge(int(cfg.MaxAge)))
	}

	return handlers.CORS(options...)
}

// NewKratosCors creates a new CORS handler for Kratos HTTP servers
// Deprecated: Use the implementation in runtime/service/transport/http package instead
func NewKratosCors(cfg *corsv1.Cors) transhttp.FilterFunc {
	if cfg == nil || !cfg.Enabled {
		return nil
	}

	// If CORS is enabled but no configuration is provided, use defaults
	if len(cfg.AllowOrigins) == 0 && len(cfg.AllowMethods) == 0 && len(cfg.AllowHeaders) == 0 {
		return handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"}),
			handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
			handlers.MaxAge(86400), // 24 hours
		)
	}

	options := []handlers.CORSOption{
		handlers.AllowedOrigins(cfg.AllowOrigins),
		handlers.AllowedMethods(cfg.AllowMethods),
		handlers.AllowedHeaders(cfg.AllowHeaders),
		handlers.ExposedHeaders(cfg.ExposeHeaders),
	}

	if cfg.AllowCredentials {
		options = append(options, handlers.AllowCredentials())
	}

	if cfg.MaxAge > 0 {
		options = append(options, handlers.MaxAge(int(cfg.MaxAge)))
	} else if cfg.MaxAge == 0 {
		// Default max age if not specified
		options = append(options, handlers.MaxAge(86400)) // 24 hours
	}

	return handlers.CORS(options...)
}