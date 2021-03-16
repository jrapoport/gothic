CUR_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
include $(CUR_DIR)/protoc.mk

PKG := github.com/jrapoport/gothic
EXE := gothic
CLI := gcli

CMD_DIR := ./cmd
BUILD_DIR := build
DEBUG_DIR := $(BUILD_DIR)/debug
RELEASE_DIR := $(BUILD_DIR)/release

GO_LINT_REPO := golang.org/x/lint/golint
GO_SEC_REPO := github.com/securego/gosec/cmd/gosec
GO_STATIC_REPO := honnef.co/go/tools/cmd/staticcheck

GO := go
GO_PATH := $(shell $(GO) env GOPATH)
GO_BIN := $(GO_PATH)/bin
GO_MOD := $(GO) mod
GO_VET := $(GO) vet
GO_CLEAN := $(GO) clean -v
GO_GEN := $(GO) generate -v
GO_BUILD := $(GO) build -v
GO_INSTALL := $(GO) install -v
GO_FMT := $(GO) fmt
GO_GET := $(GO) get -v
GO_TEST := $(GO) test -v
GO_LINT := $(GO_BIN)/golint
GO_SEC := $(GO_BIN)/gosec
GO_STATIC := $(GO_BIN)/staticcheck

GRPC_PREFIX := github.com/jrapoport/gothic/hosts/rpc
PROTO_INCLUDES := -I=hosts/rpc $(PROTO_INCLUDES)

TEST_FLAGS :=-failfast
COVERAGE_FILE=coverage.txt
COVERAGE_FLAGS=-race -covermode=atomic -coverpkg=./... -coverprofile=$(COVERAGE_FILE)
COVERAGE=0
ifeq ($(COVERAGE),1)
	TEST_FLAGS := $(TEST_FLAGS) $(COVERAGE_FLAGS)
endif

$(GO_LINT):
	$(GO_GET) $(GO_LINT_REPO)

$(GO_SEC):
	$(GO_GET) $(GO_SEC_REPO)

$(GO_STATIC):
	$(GO_GET) $(GO_STATIC_REPO)

help: ## Show this help
	echo $(BUILD_NUM)
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

fmt: ## Format code
	$(GO_FMT) ./...

vet: ## Run vet
	$(GO_VET) ./...

lint: $(GO_LINT) ## Run linter
	$(GO_LINT) ./...

audit: $(GO_SEC) ## Run audit
	$(GO_SEC) ./...

static: $(GO_STATIC) ## Run static analysis
	$(GO_STATIC) ./...

tidy: ## Tidy module
	$(GO_MOD) tidy

deps: tidy ## Install dependencies
	$(GO_MOD) download

rpc::
	$(GO_GEN) ./...

rpcw:: GRPC_PREFIX := github.com/jrapoport/gothic/hosts
# rpcw:: PROTO_WILDCARD := *_web.proto
# rpcw:: PROTO_FILES += ./hosts/rpc/response.proto

test: ## Run tests
ifeq (, $(shell which docker))
	curl -fsSL https://get.docker.com -o get-docker.sh
	sh get-docker.sh
endif
	$(GO_TEST) $(BUILD_TAGS) $(TEST_FLAGS) ./...

cover: TEST_FLAGS := $(TEST_FLAGS) $(COVERAGE_FLAGS)
cover: clean test
	curl -fsSL https://codecov.io/bash | bash
	$(RM) $(COVERAGE_FILE)

VERSION_NUM := $(shell git describe --abbrev=0 --tags 2> /dev/null)
ifeq (, $(VERSION_NUM))
	VERSION_NUM := 0.0.1
endif
BUILD_MN := $(shell git log -1 --format=%cd --date=format:'%m')
BUILD_YR := $(shell git log -1 --format=%cd --date=format:'%y%d')
BUILD_NUM := $(shell printf '%b%s' \\$(shell printf %o $(shell expr $(shell date +%m) + 64)) $(shell date +%y%d))
VER_PKG := $(PKG)/config
VER_FLAGS = -X '${VER_PKG}.Version=${VERSION_NUM}' -X '${VER_PKG}.Build=${BUILD_NUM}'
DEBUG_TAGS := -tags "debug"
RELEASE_TAGS := -tags "osusergo,netgo,release"
BUILD_TAGS := $(DEBUG_TAGS) -tags "sqlite_json"
OUT_DIR := $(DEBUG_DIR)
IN_EXE = $(CMD_DIR)/gothic
OUT_EXE = -o $(OUT_DIR)/$(EXE)
IN_CLI = $(CMD_DIR)/cli
OUT_CLI = -o $(OUT_DIR)/$(CLI)
build: ## Debug build
	echo $(VER_FLAGS)
	$(GO_BUILD) $(OUT_EXE) $(BUILD_TAGS) -ldflags="$(LD_FLAGS) $(VER_FLAGS)" $(IN_EXE)
	$(GO_BUILD) $(OUT_CLI) $(BUILD_TAGS) -ldflags="$(LD_FLAGS) $(VER_FLAGS)" $(IN_CLI)

release: BUILD_TAGS := $(RELEASE_TAGS)
release: OUT_DIR := $(RELEASE_DIR)
release: LD_FLAGS := -s -w
release: CGO_ENABLED=0
release: build ## Production build

install: OUT_EXE :=
install: OUT_CLI :=
install: GO_BUILD = $(GOINSTALL)
install: release ## Install gothic

clean: ## Clean
	$(RM) -r $(BUILD_DIR)

all: lint vet test release ## Lint, vet, test, & release

image: ## Build the Docker image
	docker build .

gothic:  ## Start gothic
	docker-compose -f docker-compose.yaml up -d gothic

envoy: ## Start envoy
	docker-compose -f docker-compose.yaml up -d envoy

mysql: ## Start mysql
	docker-compose -f docker-compose-dev.yaml up -d mysql

pg: ## Start postgres
	docker-compose -f docker-compose-dev.yaml up -d pg

db: mysql ## Start mysql db

.PHONY: help fmt vet lint audit static tidy deps rpc test build \
		release install all image gothic envoy mysql pg db cover

.DEFAULT_GOAL := build
