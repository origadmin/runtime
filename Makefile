GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)

ENV=dev
PROJECT_ORG=OrigAdmin
THIRD_PARTY_PATH=third_party

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	#GIT_BASH=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell which git))))

	# gitHash Current commit id, same as gitCommit result
	gitHash = $(shell git rev-parse HEAD)
	VERSION=$(shell git describe --tags --always)
	BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
	HEAD_TAG=$(shell git tag --points-at '${gitHash}')

	# Use PowerShell to find .proto files, convert to relative paths, and replace \ with /
	RUNTIME_PROTO_FILES := $(shell powershell -Command "Get-ChildItem -Recurse proto -Filter *.proto | Resolve-Path -Relative")
	# TOOLKITS_PROTO_FILES := $(shell powershell -Command "Get-ChildItem -Recurse toolkits -Filter *.proto | Resolve-Path -Relative")
	API_PROTO_FILES := $(shell powershell -Command "Get-ChildItem -Recurse api -Filter *.proto | Resolve-Path -Relative")

	# Replace \ with /
	RUNTIME_PROTO_FILES := $(subst \,/, $(RUNTIME_PROTO_FILES))
    # TOOLKITS_PROTO_FILES := $(subst \,/, $(TOOLKITS_PROTO_FILES))
	API_PROTO_FILES := $(subst \,/, $(API_PROTO_FILES))

	BUILT_DATE = $(shell powershell -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssK'")
	TREE_STATE = $(shell powershell -Command "if ((git status) -match 'clean') { 'clean' } else { 'dirty' }")
	TAG = $(shell powershell -Command "if ((git tag --points-at '${gitHash}') -match '^v') { '$(HEAD_TAG)' } else { '${gitHash}' }")
	# buildDate = $(shell TZ=Asia/Shanghai date +%F\ %T%z | tr 'T' ' ')
	# same as gitHash previously
	COMMIT = $(shell git log --pretty=format:'%h' -n 1)
else
	# gitHash Current commit id, same as gitCommit result
    gitHash = $(shell git rev-parse HEAD)
	VERSION=$(shell git describe --tags --always)
	BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
	HEAD_TAG=$(shell git tag --points-at '${gitHash}')

	RUNTIME_PROTO_FILES=$(shell find runtime -name *.proto)
	# TOOLKITS_PROTO_FILES=$(shell find toolkits -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)

    BUILT_DATE = $(shell TZ=Asia/Shanghai date +%FT%T%z)
    TREE_STATE := $(if $(shell git status | grep -q 'clean'),clean,dirty)
    TAG = $(shell if git tag --points-at "${gitHash}" | grep -q '^v'; then echo $(HEAD_TAG); else echo ${gitHash}; fi)
	# buildDate = $(shell TZ=Asia/Shanghai date +%F\ %T%z | tr 'T' ' ')
	# same as gitHash previously
	COMMIT = $(shell git log --pretty=format:'%h' -n 1)
endif

BUILT_BY = $(PROJECT_ORG)

ifeq ($(ENV), dev)
#    BUILD_FLAGS = -race
endif

ifeq ($(ENV), release)
    LDFLAGS = -s -w
endif

# Protoc Plugin Output Flags - 定义所有 protoc 插件的输出标志为变量
PROTOC_GO_OUT          = --go_out=paths=source_relative
PROTOC_GRPC_OUT        = --go-grpc_out=paths=source_relative
PROTOC_HTTP_OUT        = --go-http_out=paths=source_relative
PROTOC_ERRORS_OUT      = --go-errors_out=paths=source_relative
PROTOC_VALIDATE_OUT    = --validate_out=lang=go,paths=source_relative
PROTOC_OPENAPI_OUT     = --openapi_out=paths=source_relative
PROTOC_GINS_OUT        = --go-gins_out=paths=source_relative # 这个插件目前只在 examples 中使用

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@latest
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@latest

.PHONY: deps
# export protobuf dependencies to ./third_party
deps:
	@echo Exporting errors/errors.proto dependencies to $(THIRD_PARTY_PATH)
	@buf export buf.build/kratos/apis -o $(THIRD_PARTY_PATH)

	@echo Exporting googleapis/googleapis.proto dependencies to $(THIRD_PARTY_PATH)
	@buf export buf.build/bufbuild/protovalidate -o $(THIRD_PARTY_PATH)

	@echo Exporting rpcerr/rpcerr.proto dependencies to $(THIRD_PARTY_PATH)
	@buf export buf.build/googleapis/googleapis -o $(THIRD_PARTY_PATH)

	@echo Exporting wellknowntypes/wellknowntypes.proto dependencies to $(THIRD_PARTY_PATH)
	@buf export buf.build/protocolbuffers/wellknowntypes -o $(THIRD_PARTY_PATH)

	@echo Exporting validate/validate.proto dependencies to $(THIRD_PARTY_PATH), please use buf.build/bufbuild/protovalidate instead
	@buf export buf.build/envoyproxy/protoc-gen-validate -o $(THIRD_PARTY_PATH)

	@echo Exporting google.golang.org/protobuf/protovalidate/protocolbuffers/go dependencies to $(THIRD_PARTY_PATH)
	buf export buf.build/gnostic/gnostic -o $(THIRD_PARTY_PATH)


.PHONY: examples
# generate examples proto
examples:
	cd examples && protoc \
	-I./proto \
	-I../third_party \
	$(PROTOC_GO_OUT):./proto \
	$(PROTOC_GINS_OUT):./proto \
	$(PROTOC_GRPC_OUT):./proto \
	$(PROTOC_HTTP_OUT):./proto \
	$(PROTOC_ERRORS_OUT):./proto \
	$(PROTOC_OPENAPI_OUT):./proto \
	./proto/helloworld/v1/helloworld.proto

#.PHONY: server
## server used generate a service at first
#server:
#	kratos proto server -t ./internal/mods/helloworld/service ./api/v1/protos/helloworld/greeter.proto
#
#.PHONY: client
## client used when proto file is in the same directory
#client:
#	kratos proto client ./api


.PHONY: update
#update
update:
	go get -u github.com/goexts/generic@latest
	go get -u github.com/origadmin/toolkits@latest
	go get -u github.com/origadmin/toolkits/codec@latest
	go get -u github.com/origadmin/toolkits/errors@latest
	go mod tidy

.PHONY: runtime
# generate internal proto
runtime:
	protoc \
	-I./api/proto \
	-I$(THIRD_PARTY_PATH) \
	$(PROTOC_GO_OUT):./api/gen/go \
	$(PROTOC_VALIDATE_OUT):./api/gen/go \
	$(RUNTIME_PROTO_FILES)

# Find all 'proto' directories under 'test/integration' for dynamic generation
TEST_INTEGRATION_PROTO_DIRS = $(shell find test/integration -maxdepth 2 -type d -name "proto")

.PHONY: generate-test-protos
# generate proto files for integration tests
generate-test-protos:
	@echo "Generating protos for integration tests..."
	@for dir in $(TEST_INTEGRATION_PROTO_DIRS); do \
		echo "  Processing $$dir"; \
		protoc \
		-I$$dir \
		-I./api/proto \
		-I$(THIRD_PARTY_PATH) \
		$(PROTOC_GO_OUT):$$dir \
		$$dir/*.proto; \
		done

.PHONY: test
# run Go unit and integration tests
test: generate-test-protos
	go test ./...

.PHONY: generate
# run go generate to generate code
generate:
	go generate ./...

.PHONY: all
# generate all
all:
	$(MAKE) init
	$(MAKE) update
	$(MAKE) runtime
	$(MAKE) generate

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
