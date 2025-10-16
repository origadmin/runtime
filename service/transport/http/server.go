package http

import (
	"net/http/pprof"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
)

// NewServer creates a new concrete HTTP server instance based on the provided configuration.
// It returns *transhttp.Server, not the generic interfaces.Server.
func NewServer(httpConfig *transportv1.HttpServerConfig, serverOpts *ServerOptions) (*transhttp.Server, error) {
	// Initialize the Kratos HTTP server options using the adapter function.
	kratosOpts, err := initHttpServerOptions(httpConfig, serverOpts)
	if err != nil {
		return nil, err
	}

	// Create the HTTP server instance.
	srv := transhttp.NewServer(kratosOpts...)

	return srv, nil
}

// registerPprof registers the pprof handlers with the HTTP server.
func registerPprof(srv *transhttp.Server) {
	srv.HandleFunc("/debug/pprof", pprof.Index)
	srv.HandleFunc("/debug/cmdline", pprof.Cmdline)
	srv.HandleFunc("/debug/profile", pprof.Profile)
	srv.HandleFunc("/debug/symbol", pprof.Symbol)
	srv.HandleFunc("/debug/trace", pprof.Trace)
	srv.HandleFunc("/debug/allocs", pprof.Handler("allocs").ServeHTTP)
	srv.HandleFunc("/debug/block", pprof.Handler("block").ServeHTTP)
	srv.HandleFunc("/debug/goroutine", pprof.Handler("goroutine").ServeHTTP)
	srv.HandleFunc("/debug/heap", pprof.Handler("heap").ServeHTTP)
	srv.HandleFunc("/debug/mutex", pprof.Handler("mutex").ServeHTTP)
	srv.HandleFunc("/debug/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}
