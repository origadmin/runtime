GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)

ENV=dev
PROJECT_ORG=OrigAdmin
THIRD_PARTY_PATH=third_party

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	#GIT_BASH=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell which git))))

	VERSION=$(shell git describe --tags --always)
	BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
	HEAD_TAG=$(shell git tag --points-at '${gitHash}')
	# gitHash Current commit id, same as gitCommit result
	gitHash = $(shell git rev-parse HEAD)

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
	VERSION=$(shell git describe --tags --always)
	BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
	HEAD_TAG=$(shell git tag --points-at '${gitHash}')
	# gitHash Current commit id, same as gitCommit result
    gitHash = $(shell git rev-parse HEAD)

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
MODULE_PATH=github.com/origadmin/toolkits/version
LDFLAGS := -X $(MODULE_PATH).gitTag=$(TAG) \
           -X $(MODULE_PATH).buildDate=$(BUILT_DATE) \
           -X $(MODULE_PATH).gitCommit=$(COMMIT) \
           -X $(MODULE_PATH).gitTreeState=$(TREE_STATE) \
           -X $(MODULE_PATH).gitBranch=$(BRANCH) \
           -X $(MODULE_PATH).version=$(VERSION)

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
# generate examples proto or use ./examples/generate.go
examples:
	cd examples && protoc \
	--proto_path=./proto \
	--proto_path=../third_party \
	--go_out=paths=source_relative:./proto \
	--go-gins_out=paths=source_relative:./proto \
	--go-grpc_out=paths=source_relative:./proto \
	--go-http_out=paths=source_relative:./proto \
	--go-errors_out=paths=source_relative:./proto \
	--openapi_out=paths=source_relative:./proto \
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
# generate internal proto or use ./internal/generate.go
runtime:
	protoc \
	--proto_path=./proto \
	--proto_path=./third_party \
	--go_out=paths=source_relative:./gen/go \
	--validate_out=lang=go,paths=source_relative:./gen/go \
	$(RUNTIME_PROTO_FILES)

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

.PHONY: version
# run version example
version:
	go run -ldflags "$(LDFLAGS)" $(BUILD_FLAGS) -gcflags=all="-N -l" ./cmd/version/main.go

.PHONY: build-gins
# build protoc-gen-go-gins with current snapshot to dist
build-gins:
	goreleaser build --single-target --clean --snapshot -f ./cmd/protoc-gen-go-gins/.goreleaser.yaml

.PHONY: release-gins
# release
release:
	goreleaser release --clean -f ./cmd/protoc-gen-go-gins/.goreleaser.yaml

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
