module github.com/origadmin/runtime/tools/protoc-gen-go-security

go 1.24.6

require (
	google.golang.org/genproto/googleapis/api v0.0.0-20251103181224-f26f9409b101
	google.golang.org/protobuf v1.36.10
	github.com/origadmin/runtime v0.2.13-dev
)

replace (
	github.com/origadmin/runtime v0.2.13-dev => ../../
)