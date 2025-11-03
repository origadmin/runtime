package builders

import (
	"google.golang.org/protobuf/types/known/durationpb"

	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/runtime/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/runtime/logger/v1"
	corsv1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/cors/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	selectorv1 "github.com/origadmin/runtime/api/gen/go/runtime/selector/v1"
	tracev1 "github.com/origadmin/runtime/api/gen/go/runtime/trace/v1"
	grpcv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/grpc/v1"
	httpv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/http/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
)

func NewDefaultApp() *appv1.App {
	return &appv1.App{
		Id:      "test-app-id",
		Name:    "TestApp",
		Version: "1.0.0",
		Env:     "test",
		Metadata: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}
}

func NewDefaultServers() *transportv1.Servers {
	return &transportv1.Servers{
		Servers: []*transportv1.Server{
			{
				Name:     "grpc_server",
				Protocol: "grpc",
				Grpc: &grpcv1.Server{
					Network: "tcp",
					Addr:    ":9000",
					Timeout: durationpb.New(1000000000), // 1s
				},
			},
			{
				Name:     "http_server",
				Protocol: "http",
				Http: &httpv1.Server{
					Network: "tcp",
					Addr:    ":8000",
					Timeout: durationpb.New(2000000000), // 2s
				},
			},
		},
	}
}

func NewDefaultClient() *transportv1.Client {
	return &transportv1.Client{
		Grpc: &grpcv1.Client{
			Endpoint: "discovery:///user-service",
			Timeout:  durationpb.New(3000000000), // 3s
			Selector: &selectorv1.SelectorConfig{
				Version: "v1.0.0",
			},
		},
	}
}

func NewDefaultLogger() *loggerv1.Logger {
	return &loggerv1.Logger{
		Level:  "info",
		Format: "json",
		Stdout: true,
	}
}

func NewDefaultDiscoveries() *discoveryv1.Discoveries {
	return &discoveryv1.Discoveries{
		Discoveries: []*discoveryv1.Discovery{
			{
				Name: "internal-consul",
				Type: "consul",
				Consul: &discoveryv1.Consul{
					Address: "consul.internal:8500",
				},
			},
			{
				Name: "legacy-etcd",
				Type: "etcd", Etcd: &discoveryv1.ETCD{
				Endpoints: []string{"etcd.legacy:2379"},
			},
			},
		},
	}
}

func NewDefaultTrace() *tracev1.Trace {
	return &tracev1.Trace{
		Name:     "jaeger",
		Endpoint: "http://jaeger:14268/api/traces",
	}
}

func NewDefaultMiddlewares() *middlewarev1.Middlewares {
	return &middlewarev1.Middlewares{
		Middlewares: []*middlewarev1.Middleware{
			{
				Name:    "cors-middleware",
				Type:    "cors",
				Enabled: true,
				Cors: &corsv1.Cors{
					AllowedOrigins: []string{"*"},
					AllowedMethods: []string{"GET", "POST"},
				},
			},
		},
	}
}
