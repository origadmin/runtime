package interfaces

//type MiddlewareProvider interface {
//	// BuildClient builds client-side middleware using the provided ConfigDecoder and config path.
//	BuildClient(decoder ConfigDecoder, path string, opts ...interface{}) []kratosmiddleware.Middleware
//	// BuildServer builds server-side middleware using the provided ConfigDecoder and config path.
//	BuildServer(decoder ConfigDecoder, path string, opts ...interface{}) []kratosmiddleware.Middleware
//}
//
//// ServerBuilder is an interface for building transport servers (HTTP, gRPC).
//type ServerBuilder interface {
//	// DefaultBuild builds a default server using the provided ConfigDecoder and config path.
//	DefaultBuild(decoder ConfigDecoder, path string, opts ...interface{}) (transport.Server, error)
//	// Build builds a named server using the provided ConfigDecoder and config path.
//	Build(name string, decoder ConfigDecoder, path string, opts ...interface{}) (transport.Server, error)
//}
