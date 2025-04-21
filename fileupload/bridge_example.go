package fileupload

import (
	"context"
	"log"
	"net/http"
	"time"

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
			Timeout:         int64(3 * time.Minute),
			ShutdownTimeout: int64(3 * time.Minute),
			ReadTimeout:     int64(3 * time.Minute),
			WriteTimeout:    int64(3 * time.Minute),
			IdleTimeout:     int64(3 * time.Minute),
			Endpoint:        "",
		},
		Http: &configv1.Service_HTTP{
			Network:         "",
			Addr:            "",
			UseTls:          false,
			CertFile:        "",
			KeyFile:         "",
			Timeout:         int64(3 * time.Minute),
			ShutdownTimeout: int64(3 * time.Minute),
			ReadTimeout:     int64(3 * time.Minute),
			WriteTimeout:    int64(3 * time.Minute),
			IdleTimeout:     int64(3 * time.Minute),
			Endpoint:        "",
		},
		Websocket: &configv1.WebSocket{
			Network: "",
			Addr:    "",
			Path:    "",
			Codec:   "",
			Timeout: int64(3 * time.Minute),
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
