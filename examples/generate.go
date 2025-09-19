// Package examples implements the functions, types, and interfaces for the module.
package examples

//go:generate protoc --proto_path=. --proto_path=../api/proto --proto_path=../third_party --go_out=paths=source_relative:. ./proto/api_gateway/bootstrap.proto
//go:generate protoc --proto_path=. --proto_path=../api/proto --proto_path=../third_party --go_out=paths=source_relative:. ./proto/http_server_grpc_client/bootstrap.proto
//go:generate protoc --proto_path=. --proto_path=../api/proto --proto_path=../third_party --go_out=paths=source_relative:. ./proto/simple_grpc_server/bootstrap.proto
