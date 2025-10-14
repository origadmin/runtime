# Makefile for the OrigAdmin Framework Runtime

# ---------------------------------------------------------------------------- #
#                             Build Configuration                              #
# ---------------------------------------------------------------------------- #

# Basic configuration
GOHOSTOS         ?= $(shell go env GOHOSTOS)
ENV              ?= dev
PROJECT_ORG      := OrigAdmin
THIRD_PARTY_PATH := ./third_party
BUILT_BY         := $(PROJECT_ORG)

# Common Git information
GIT_COMMIT      := $(shell git rev-parse --short HEAD)
GIT_BRANCH      := $(shell git rev-parse --abbrev-ref HEAD)
GIT_VERSION     := $(shell git describe --tags --always)

# OS-specific variables
ifeq ($(GOHOSTOS), windows)
    SHELL          := powershell.exe
    .SHELLFLAGS    := -NoProfile -Command
    GIT_HEAD_TAG   := $(shell git tag --points-at HEAD 2>$null)
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

# Build flags for release
ifeq ($(ENV), release)
    LDFLAGS = -s -w
endif


# ---------------------------------------------------------------------------- #
#                         Proto File & Plugin Config                         #
# ---------------------------------------------------------------------------- #

# Protoc plugin definitions
PROTOC_GO_OUT       := --go_out=paths=source_relative
PROTOC_GRPC_OUT     := --go-grpc_out=paths=source_relative
PROTOC_HTTP_OUT     := --go-http_out=paths=source_relative
PROTOC_ERRORS_OUT   := --go-errors_out=paths=source_relative
PROTOC_VALIDATE_OUT := --validate_out=lang=go,paths=source_relative

# Proto file discovery
ifeq ($(GOHOSTOS), windows)
    API_PROTO_FILES     := $(subst \,/, $(shell powershell -Command "(Get-ChildItem -Recurse ./api/proto -Filter *.proto | Resolve-Path -Relative) -join ' '"))
    TEST_PROTO_DIRS     := $(subst \,/, $(shell powershell -Command "(Get-ChildItem -Recurse ./test/integration -Directory -Filter proto | Resolve-Path -Relative) -join ' '"))
else
    API_PROTO_FILES     := $(shell find ./api/proto -name '*.proto')
    TEST_PROTO_DIRS     := $(shell find ./test/integration -maxdepth 2 -type d -name "proto")
endif


# ---------------------------------------------------------------------------- #
#                         Tooling & Dependency Management                      #
# ---------------------------------------------------------------------------- #

.PHONY: init deps update

init: ## 🔧 Install all required Go tools for development
	@echo "Installing development tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/bufbuild/buf/cmd/buf@latest

deps: ## 📦 Export and install all third-party protobuf dependencies
	@echo "Exporting protobuf dependencies to $(THIRD_PARTY_PATH)..."
	@buf export buf.build/kratos/apis -o $(THIRD_PARTY_PATH)
	@buf export buf.build/bufbuild/protovalidate -o $(THIRD_PARTY_PATH)
	@buf export buf.build/googleapis/googleapis -o $(THIRD_PARTY_PATH)
	@buf export buf.build/protocolbuffers/wellknowntypes -o $(THIRD_PARTY_PATH)
	@buf export buf.build/envoyproxy/protoc-gen-validate -o $(THIRD_PARTY_PATH)
	@buf export buf.build/gnostic/gnostic -o $(THIRD_PARTY_PATH)

update: ## 🔄 Update Go module dependencies
	@echo "Updating Go dependencies..."
	go get -u github.com/goexts/generic@latest
	go get -u github.com/origadmin/toolkits@latest
	go mod tidy


# ---------------------------------------------------------------------------- #
#                             Proto Generation                               #
# ---------------------------------------------------------------------------- #

.PHONY: protos generate-test-protos

protos: ## 🧬 Generate all API protos into ./api/gen/go
	@echo "Generating API protos..."
ifeq ($(GOHOSTOS), windows)
	@protoc --proto_path=./api/proto --proto_path=$(THIRD_PARTY_PATH) $(PROTOC_GO_OUT):./api/gen/go $(PROTOC_GRPC_OUT):./api/gen/go $(PROTOC_HTTP_OUT):./api/gen/go $(PROTOC_ERRORS_OUT):./api/gen/go $(PROTOC_VALIDATE_OUT):./api/gen/go $(API_PROTO_FILES)
else
	@protoc \
		--proto_path=./api/proto \
		--proto_path=$(THIRD_PARTY_PATH) \
		$(PROTOC_GO_OUT):./api/gen/go \
		$(PROTOC_GRPC_OUT):./api/gen/go \
		$(PROTOC_HTTP_OUT):./api/gen/go \
		$(PROTOC_ERRORS_OUT):./api/gen/go \
		$(PROTOC_VALIDATE_OUT):./api/gen/go \
		$(API_PROTO_FILES)
endif

generate-test-protos: ## Generate protos for integration tests (cross-platform)
	@echo "Generating protos for integration tests..."
ifeq ($(GOHOSTOS), windows)
	@foreach ($$dir in '$(TEST_PROTO_DIRS)'.Split(' ')) { if ($$dir) { Write-Host "  Processing $$dir"; protoc -I. -I$(THIRD_PARTY_PATH) -I./api/proto --go_out=paths=source_relative:. "$$dir/*.proto" } }
else
	@for dir in $(TEST_PROTO_DIRS); do \
		echo "  Processing $$dir"; \
		protoc -I. -I$(THIRD_PARTY_PATH) -I./api/proto --go_out=paths=source_relative:. $$dir/*.proto; \
	done
endif


# ---------------------------------------------------------------------------- #
#                           Main Lifecycle Targets                           #
# ---------------------------------------------------------------------------- #

.PHONY: test generate all

test: generate-test-protos ## 🧪 Run all Go tests
	go test ./...

generate: ## 🧬 Run go generate to generate code (e.g., wire)
	go generate ./...

all: init deps protos generate-test-protos generate ## ✅ Run the full build process


# ---------------------------------------------------------------------------- #
#                                     Help                                     #
# ---------------------------------------------------------------------------- #

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
