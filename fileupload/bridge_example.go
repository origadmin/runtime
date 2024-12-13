package fileupload

import (
	"context"
	"log"
	"net/http"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
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
			Timeout:         durationpb.New(3 * time.Minute),
			ShutdownTimeout: durationpb.New(3 * time.Minute),

			ReadTimeout: durationpb.New(3 * time.Minute),

			WriteTimeout: durationpb.New(3 * time.Minute),

			IdleTimeout: durationpb.New(3 * time.Minute),
			Endpoint:    "",
		},
		Http: &configv1.Service_HTTP{
			Network:  "",
			Addr:     "",
			UseTls:   false,
			CertFile: "",
			KeyFile:  "",
			Timeout:  durationpb.New(3 * time.Minute),

			ShutdownTimeout: durationpb.New(3 * time.Minute),

			ReadTimeout: durationpb.New(3 * time.Minute),

			WriteTimeout: durationpb.New(3 * time.Minute),
			IdleTimeout:  durationpb.New(3 * time.Minute),

			Endpoint: "",
		},
		Gins: &configv1.Service_GINS{
			Network:  "",
			Addr:     "",
			UseTls:   false,
			CertFile: "",
			KeyFile:  "",
			Timeout:  durationpb.New(3 * time.Minute),

			ShutdownTimeout: durationpb.New(3 * time.Minute),

			ReadTimeout: durationpb.New(3 * time.Minute),

			WriteTimeout: durationpb.New(3 * time.Minute),

			IdleTimeout: durationpb.New(3 * time.Minute),

			Endpoint: "",
		},
		Websocket: &configv1.WebSocket{
			Network: "",
			Addr:    "",
			Path:    "",
			Codec:   "",
			Timeout: durationpb.New(3 * time.Minute),
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
		Middleware: &configv1.Middleware{
			EnableLogging:        false,
			EnableRecovery:       false,
			EnableTracing:        false,
			EnableValidate:       false,
			EnableCircuitBreaker: false,
			EnableMetadata:       false,
			RateLimiter: &configv1.Middleware_RateLimiter{
				Name:                "",
				Period:              0,
				XRatelimitLimit:     0,
				XRatelimitRemaining: 0,
				XRatelimitReset:     0,
				RetryAfter:          0,
				Memory: &configv1.Middleware_RateLimiter_Memory{
					Expiration:      durationpb.New(3 * time.Minute),
					CleanupInterval: durationpb.New(3 * time.Minute),
				},
				Redis: &configv1.Middleware_RateLimiter_Redis{
					Addr:     "",
					Username: "",
					Password: "",
					Db:       0,
				},
			},
			Metadata: &configv1.Middleware_Metadata{
				Prefix: "",
				Data:   nil,
			},
			Metrics: &configv1.Middleware_Metrics{
				SupportedMetrics: nil,
				UserMetrics:      nil,
			},
			Validator: &configv1.Middleware_Validator{
				Version:  0,
				FailFast: false,
			},
			Security: &configv1.Security{
				PublicPaths: nil,
				Authz: &configv1.AuthZConfig{
					Disabled:    false,
					PublicPaths: nil,
					Type:        "",
					Casbin: &configv1.AuthZConfig_CasbinConfig{
						PolicyFile: "",
						ModelFile:  "",
					},
					Opa: &configv1.AuthZConfig_OpaConfig{
						PolicyFile: "",
						DataFile:   "",
						ServerUrl:  "",
						RegoFile:   "",
					},
					Zanzibar: &configv1.AuthZConfig_ZanzibarConfig{
						ApiEndpoint:      "",
						Namespace:        "",
						ReadConsistency:  "",
						WriteConsistency: "",
					},
				},
				Authn: &configv1.AuthNConfig{
					Disabled:    false,
					PublicPaths: nil,
					Type:        "",
					Jwt: &configv1.AuthNConfig_JWTConfig{
						Algorithm:     "",
						SigningKey:    "",
						OldSigningKey: "",
						ExpireTime:    durationpb.New(3 * time.Minute),
						RefreshTime:   durationpb.New(3 * time.Minute),
						CacheName:     "",
					},
					Oidc: &configv1.AuthNConfig_OIDCConfig{
						IssuerUrl: "",
						Audience:  "",
						Algorithm: "",
					},
					PreSharedKey: &configv1.AuthNConfig_PreSharedKeyConfig{
						SecretKeys: nil,
					},
				},
			},
		},
		Selector: &configv1.Service_Selector{
			Version: "",
			Builder: "",
		},
		HostName: "",
		HostIp:   "",
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
