# Makefile for the OrigAdmin Framework Runtime

# ============================================================================ #
#                              CONFIGURATION
# ============================================================================ #

# --------------------------- Basic Configuration ---------------------------- #
GOHOSTOS         ?= $(shell go env GOHOSTOS)
ENV              ?= dev
PROJECT_ORG      := OrigAdmin
THIRD_PARTY_PATH := ./third_party
BUILT_BY         := $(PROJECT_ORG)

# ---------------------------- Git Information ----------------------------- #
GIT_COMMIT      := $(shell git rev-parse --short HEAD)
GIT_BRANCH      := $(shell git rev-parse --abbrev-ref HEAD)
GIT_VERSION     := $(shell git describe --tags --always)

# ------------------------- OS-Specific Variables -------------------------- #
ifeq ($(GOHOSTOS), windows)
    SHELL          := powershell.exe
    .SHELLFLAGS    := -NoProfile -Command
    GIT_HEAD_TAG   := $(shell git tag --points-at HEAD 2>$$null)
    BUILD_DATE     := $(shell powershell -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssK'")
    GIT_TREE_STATE := $(shell powershell -Command "if ((git status --porcelain)) { 'dirty' } else { 'clean' }")
    GIT_TAG        := $(shell powershell -Command "if ('${GIT_HEAD_TAG}') { '${GIT_HEAD_TAG}' } else { '${GIT_COMMIT}' }")
else
    SHELL          := /bin/bash
    GIT_HEAD_TAG   := $(shell git tag --points-at HEAD 2>/dev/null)
    BUILD_DATE     := $(shell TZ=Asia/Shanghai date +%FT%T%z)
    GIT_TREE_STATE := $(if $(shell git status --porcelain),dirty,clean)
    GIT_TAG        := $(if $(GIT_HEAD_TAG),$(GIT_HEAD_TAG),$(GIT_COMMIT))
endif

# Append -dirty suffix if the working directory is not clean
ifneq ($(GIT_TREE_STATE), clean)
    GIT_VERSION := $(GIT_VERSION)-dirty
endif

# ----------------------------- Build Flags ------------------------------ #
ifeq ($(ENV), release)
    LDFLAGS = -s -w
endif

# ------------------------ Protobuf Configuration ------------------------ #
# Common protoc include paths, used for both buf and protoc fallback
PROTOC_INCLUDES := -I. -I./api/proto -I$(THIRD_PARTY_PATH)

# Protoc plugin definitions, used for protoc fallback
PROTOC_GO_OUT       := --go_out=paths=source_relative
PROTOC_GRPC_OUT     := --go-grpc_out=paths=source_relative
PROTOC_HTTP_OUT     := --go-http_out=paths=source_relative
PROTOC_ERRORS_OUT   := --go-errors_out=paths=source_relative
PROTOC_VALIDATE_OUT := --validate_out=lang=go,paths=source_relative

# A single variable for all proto plugins used in the main API generation
PLUGINS := $(PROTOC_GO_OUT):./api/gen/go \
		$(PROTOC_GRPC_OUT):./api/gen/go \
		$(PROTOC_HTTP_OUT):./api/gen/go \
		$(PROTOC_ERRORS_OUT):./api/gen/go \
		$(PROTOC_VALIDATE_OUT):./api/gen/go

# A single variable for all proto plugins used in example generation (output to current dir)
EXAMPLE_PLUGINS := $(PROTOC_GO_OUT):. \
		$(PROTOC_GRPC_OUT):. \
		$(PROTOC_HTTP_OUT):. \
		$(PROTOC_ERRORS_OUT):. \
		$(PROTOC_VALIDATE_OUT):.

# Proto file discovery
ifeq ($(GOHOSTOS), windows)
    API_PROTO_FILES     := $(subst \,/, $(shell powershell -Command "(Get-ChildItem -Recurse ./api/proto -Filter *.proto | Resolve-Path -Relative) -join ' '"))
    TEST_PROTO_DIRS     := $(subst \,/, $(shell powershell -Command "(Get-ChildItem -Recurse ./test/integration -Directory -Filter proto | Resolve-Path -Relative) -join ' '"))
    EXAMPLE_PROTO_FILES := $(subst \,/, $(shell powershell -Command "(Get-ChildItem -Recurse ./examples/protos -Filter *.proto | Resolve-Path -Relative) -join ' '"))
else
    API_PROTO_FILES     := $(shell find ./api/proto -name '*.proto')
    TEST_PROTO_DIRS     := $(shell find ./test/integration -maxdepth 2 -type d -name "proto")
    EXAMPLE_PROTO_FILES := $(shell find ./examples/protos -name '*.proto')
endif


# ============================================================================ #
#                           LIFECYCLE TARGETS
# ============================================================================ #

.PHONY: all init deps update update-tools protos generate-examples-protos generate-test-protos generate test clean clean-api-gen clean-examples-protos clean-test-protos

all: init deps protos generate-examples-protos generate-test-protos generate ## ✅ Run the full build process

init: ## 🔧 Install tools from tools.go, ensuring reproducible builds
	@echo "Ensuring tool dependencies are in go.mod..."
	@go mod tidy
	@echo "Installing tools listed in tools.go..."
	@go install github.com/bufbuild/buf/cmd/buf
	@go install github.com/envoyproxy/protoc-gen-validate
	@go install github.com/go-kratos/kratos/cmd/kratos/v2
	@go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2
	@go install github.com/google/gnostic/cmd/protoc-gen-openapi
	@go install github.com/google/wire/cmd/wire
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@go install google.golang.org/protobuf/cmd/protoc-gen-go

deps: ## 📦 Export and install all third-party protobuf dependencies
	@echo "Exporting protobuf dependencies to $(THIRD_PARTY_PATH)..."
	@buf export buf.build/kratos/apis -o $(THIRD_PARTY_PATH)
	@buf export buf.build/bufbuild/protovalidate -o $(THIRD_PARTY_PATH)
	@buf export buf.build/googleapis/googleapis -o $(THIRD_PARTY_PATH)
	@buf export buf.build/protocolbuffers/wellknowntypes -o $(THIRD_PARTY_PATH)
	@buf export buf.build/envoyproxy/protoc-gen-validate -o $(THIRD_PARTY_PATH)
	@buf export buf.build/gnostic/gnostic -o $(THIRD_PARTY_PATH)

update: ## 🔄 Update Go module dependencies (libraries only)
	@echo "Updating Go dependencies..."
	go get -u github.com/goexts/generic@latest
	go get -u github.com/origadmin/toolkits@latest
	go mod tidy

update-tools: ## ⚠️  Update all Go tools in tools.go to latest. High-risk, use with caution!
	@echo "Updating all tools in tools.go to @latest..."
	@go get -u github.com/bufbuild/buf/cmd/buf
	@go get -u github.com/envoyproxy/protoc-gen-validate
	@go get -u github.com/go-kratos/kratos/cmd/kratos/v2
	@go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2
	@go get -u github.com/google/gnostic/cmd/protoc-gen-openapi
	@go get -u github.com/google/wire/cmd/wire
	@go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	@go mod tidy

protos: ## 🧬 Generate API protos (smartly uses buf, falls back to protoc)
	@echo "Generating API protos..."
ifeq ($(GOHOSTOS), windows)
	@if (Get-Command buf -ErrorAction SilentlyContinue) { echo '--> `buf` command found. Using buf to generate protos (recommended).'; buf generate } else { echo '--> WARNING: `buf` command not found. Falling back to protoc.'; echo '--> Please run ''make init'' to install buf for reproducible builds.'; protoc $(PROTOC_INCLUDES) $(PLUGINS) $(API_PROTO_FILES) }
else
	@if command -v buf >/dev/null 2>&1; then \
		echo "--> 'buf' command found. Using buf to generate protos (recommended)."; \
		buf generate; \
	else \
		echo "--> WARNING: 'buf' command not found. Falling back to protoc."; \
		echo "--> Please run 'make init' to install buf for reproducible builds."; \
		protoc $(PROTOC_INCLUDES) $(PLUGINS) $(API_PROTO_FILES); \
	fi
endif

generate-examples-protos: ## 🧬 Generate all example protos
	@echo "Generating example protos..."
	@protoc $(PROTOC_INCLUDES) $(EXAMPLE_PLUGINS) $(EXAMPLE_PROTO_FILES)

generate-test-protos: ## Generate protos for integration tests (cross-platform)
	@echo "Generating protos for integration tests..."
ifeq ($(GOHOSTOS), windows)
	@foreach ($$dir in '$(TEST_PROTO_DIRS)'.Split(' ')) { if ($$dir) { Write-Host "  Processing $$dir"; protoc $(PROTOC_INCLUDES) --go_out=paths=source_relative:. "$$dir/*.proto" } }
else
	@for dir in $(TEST_PROTO_DIRS); do \
		echo "  Processing $$dir"; \
		protoc $(PROTOC_INCLUDES) --go_out=paths=source_relative:. $$dir/*.proto; \
	done
endif

generate: ## 🧬 Run go generate to generate code (e.g., wire)
	@echo "Running go generate..."
	go generate ./...

test: generate-test-protos ## 🧪 Run all Go tests
	go test ./...

clean: clean-api-gen clean-examples-protos clean-test-protos ## 🧹 Clean up all generated files

clean-api-gen: ## 🧹 Clean up API generated files
	@echo "Cleaning up API generated files..."
	@rm -rf ./api/gen

clean-examples-protos: ## 🧹 Clean up example generated protos
	@echo "Cleaning up example generated protos..."
ifeq ($(GOHOSTOS), windows)
	@Get-ChildItem -Recurse ./examples/protos -Filter *.pb.go | Remove-Item -Force
else
	@find ./examples/protos -name '*.pb.go' -delete
endif

clean-test-protos: ## 🧹 Clean up test generated protos
	@echo "Cleaning up test generated protos..."
ifeq ($(GOHOSTOS), windows)
	@Get-ChildItem -Recurse ./test/integration -Filter *.pb.go | Remove-Item -Force
else
	@find ./test/integration -name '*.pb.go' -delete
endif


# ============================================================================ #
#                                     HELP
# ============================================================================ #

.PHONY: help

help: ## ✨ Show this help message
	@echo ''
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  \033[36m%-22s\033[0m %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
