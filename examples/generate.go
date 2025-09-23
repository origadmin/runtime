// Package examples implements the functions, types, and interfaces for the module.
package examples

//go:generate protoc --proto_path=. --proto_path=../api/proto --proto_path=../third_party --go_out=paths=source_relative:. ./protos/api_gateway/bootstrap.proto
//go:generate protoc --proto_path=. --proto_path=../api/proto --proto_path=../third_party --go_out=paths=source_relative:. ./protos/http_server_grpc_client/bootstrap.proto
//go:generate protoc --proto_path=. --proto_path=../api/proto --proto_path=../third_party --go_out=paths=source_relative:. ./protos/simple_grpc_server/bootstrap.proto
//go:generate protoc --proto_path=. --proto_path=../api/proto --proto_path=../third_party --go_out=paths=source_relative:. ./protos/load_with_runtime/bootstrap.proto
