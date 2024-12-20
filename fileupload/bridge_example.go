package fileupload

import (
	"context"
	"log"
	"net/http"
	"time"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	metricsv1 "github.com/origadmin/runtime/gen/go/middleware/metrics/v1"
	ratelimitv1 "github.com/origadmin/runtime/gen/go/middleware/ratelimit/v1"
	middlewarev1 "github.com/origadmin/runtime/gen/go/middleware/v1"
	validatorv1 "github.com/origadmin/runtime/gen/go/middleware/validator/v1"
	"github.com/origadmin/toolkits/fileupload"
)

func main() {

	// Create a gRPC Builder
	grpcBuilder := NewBuilder(
		WithServiceType(fileupload.ServiceTypeGRPC),
		WithURI("grpc-server-address:port"),
	)
	// Create a bridge uploader
	grpcUploader, err := NewGRPCUploader(context.Background(), &configv1.Service{
		Name:            "",
		DynamicEndpoint: false,
		Grpc: &configv1.Service_GRPC{
			Network:         "",
			Addr:            "",
			UseTls:          false,
			CertFile:        "",
			KeyFile:         "",
			Timeout:         int64(3 * time.Minute),
			ShutdownTimeout: int64(3 * time.Minute),

			ReadTimeout: int64(3 * time.Minute),

			WriteTimeout: int64(3 * time.Minute),

			IdleTimeout: int64(3 * time.Minute),
			Endpoint:    "",
		},
		Http: &configv1.Service_HTTP{
			Network:  "",
			Addr:     "",
			UseTls:   false,
			CertFile: "",
			KeyFile:  "",
			Timeout:  int64(3 * time.Minute),

			ShutdownTimeout: int64(3 * time.Minute),

			ReadTimeout: int64(3 * time.Minute),

			WriteTimeout: int64(3 * time.Minute),
			IdleTimeout:  int64(3 * time.Minute),

			Endpoint: "",
		},
		Websocket: &configv1.WebSocket{
			Network: "",
			Addr:    "",
			Path:    "",
			Codec:   "",
			Timeout: int64(3 * time.Minute),
		},
		Message: &configv1.Message{
			Type: "",
			Name: "",
			Mqtt: &configv1.Message_MQTT{
				Endpoint: "",
				Codec:    "",
			},
			Kafka: &configv1.Message_Kafka{
				Endpoint: "",
				Codec:    "",
			},
			Rabbitmq: &configv1.Message_RabbitMQ{
				Endpoint: "",
				Codec:    "",
			},
			Activemq: &configv1.Message_ActiveMQ{
				Endpoint: "",
				Codec:    "",
			},
			Nats: &configv1.Message_NATS{
				Endpoint: "",
				Codec:    "",
			},
			Nsq: &configv1.Message_NSQ{
				Endpoint: "",
				Codec:    "",
			},
			Pulsar: &configv1.Message_Pulsar{
				Endpoint: "",
				Codec:    "",
			},
			Redis: &configv1.Message_Redis{
				Endpoint: "",
				Codec:    "",
			},
			Rocketmq: &configv1.Message_RocketMQ{
				Endpoint:         "",
				Codec:            "",
				EnableTrace:      false,
				NameServers:      nil,
				NameServerDomain: "",
				AccessKey:        "",
				SecretKey:        "",
				SecurityToken:    "",
				Namespace:        "",
				InstanceName:     "",
				GroupName:        "",
			},
		},
		Task: &configv1.Task{
			Type: "",
			Name: "",
			Asynq: &configv1.Task_Asynq{
				Endpoint: "",
				Password: "",
				Db:       0,
				Location: "",
			},
			Machinery: &configv1.Task_Machinery{
				Brokers:  nil,
				Backends: nil,
			},
			Cron: &configv1.Task_Cron{
				Addr: "",
			},
		},
		Middleware: &middlewarev1.Middleware{
			Logging:        true,
			Recovery:       true,
			Tracing:        true,
			CircuitBreaker: true,
			Metadata: &middlewarev1.Middleware_Metadata{
				Prefix: "",
				Data:   nil,
			},
			RateLimiter: &ratelimitv1.RateLimiter{
				Name:                "",
				Period:              0,
				XRatelimitLimit:     0,
				XRatelimitRemaining: 0,
				XRatelimitReset:     0,
				RetryAfter:          0,
				Memory: &ratelimitv1.RateLimiter_Memory{
					Expiration:      int64(3 * time.Minute),
					CleanupInterval: int64(3 * time.Minute),
				},
				Redis: &ratelimitv1.RateLimiter_Redis{
					Addr:     "",
					Username: "",
					Password: "",
					Db:       0,
				},
			},
			Metrics: &metricsv1.Metrics{
				SupportedMetrics: nil,
				UserMetrics:      nil,
			},
			Validator: &validatorv1.Validator{
				Version:  0,
				FailFast: false,
			},
			Jwt:      nil,
			Selector: nil,
		},
		Selector: &configv1.Service_Selector{
			Version: "",
			Builder: "",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	bridgeUploader := NewBridgeUploader(grpcBuilder, grpcUploader)

	// Register the HTTP processor
	http.HandleFunc("/upload", bridgeUploader.ServeHTTP)
	// Start the HTTP server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
