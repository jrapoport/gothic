CUR_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
include $(CUR_DIR)/protoc.mk

PKG := github.com/jrapoport/gothic
EXE := gothic
CLI := gadmin

CMD_DIR := ./cmd
BUILD_DIR := build
DEBUG_DIR := $(BUILD_DIR)/debug
RELEASE_DIR := $(BUILD_DIR)/release
OUT_DIR = $(DEBUG_DIR)
IN_EXE = $(CMD_DIR)/exe
OUT_EXE = $(OUT_DIR)/$(EXE)
IN_CLI = $(CMD_DIR)/cli
OUT_CLI = $(OUT_DIR)/$(CLI)

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
GO_TEST := $(GO) test -v
GO_LINT := $(GO_BIN)/golangci-lint
GO_SEC := $(GO_BIN)/gosec
GO_STATIC := $(GO_BIN)/staticcheck
GO_GET := $(GO) get

GO_LINT_REPO := github.com/golangci/golangci-lint/cmd/golangci-lint
GO_SEC_REPO := github.com/securego/gosec/cmd/gosec
GO_STATIC_REPO := honnef.co/go/tools/cmd/staticcheck

GRPC_DIR := api/grpc
GRPC_PREFIX := github.com/jrapoport/gothic/api/grpc/rpc
PROTO_INCLUDES := -I=api/service/proto $(PROTO_INCLUDES)

DEBUG_TAGS := -tags "debug"
RELEASE_TAGS := -tags "osusergo,netgo,release"
BUILD_TAGS := $(DEBUG_TAGS) -tags "sqlite_json"

#VERSION_NUM := $(shell git describe --abbrev=0 --tags 2> /dev/null)
#ifeq (, $(VERSION_NUM))
	VERSION_NUM := 0.5.0
#endif
# BUILD_MN := $(shell git log -1 --format=%cd --date=format:'%m')
# BUILD_YR := $(shell git log -1 --format=%cd --date=format:'%y%d')
# BUILD_NUM := $(shell printf '%b%s' \\$(shell printf %o $(shell expr $(BUILD_MN) + 64)) $(BUILD_YR))
BUILD_NUM := $(shell printf '%b%s' \\$(shell printf %o $(shell expr $(shell date +%m) + 64)) $(shell date +%y%d))
VER_PKG := $(PKG)/config
# make sure this is = and not := so it gets expanded properly
VER_FLAGS = -X '${VER_PKG}.Version=${VERSION_NUM}' -X '${VER_PKG}.Build=${BUILD_NUM}' -X '${VER_PKG}.ExeName=${EXE}'

TEST_FLAGS :=-failfast -count=1 $(TEST_FLAGS)
COVERAGE_FILE=coverage.txt
COVERAGE_FLAGS=-race -covermode=atomic -coverpkg=./... -coverprofile=$(COVERAGE_FILE)
COVERAGE=0
ifeq ($(COVERAGE),1)
	TEST_FLAGS := $(TEST_FLAGS) $(COVERAGE_FLAGS)
endif

DOCKER_FLAGS := DOCKER_BUILDKIT=0
DOCKER := $(DOCKER_FLAGS) docker
DOCKER_COMPOSE := $(DOCKER_FLAGS) docker-compose # $(DOCKER) compose
COMPOSE_FILE?=docker-compose.yaml
COMPDEV_FILE:=docker-compose-dev.yaml

$(GO_LINT):
	$(GO_INSTALL) $(GO_LINT_REPO)

$(GO_SEC):
	$(GO_INSTALL) $(GO_SEC_REPO)

$(GO_STATIC):
	$(GO_INSTALL) $(GO_STATIC_REPO)

help: ## Show this help
	echo $(BUILD_NUM)
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

fmt: ## Format code
	$(GO_FMT) ./...

vet: ## Run vet
	$(GO_VET) ./...

lint: $(GO_LINT) ## Run linter
	$(GO_LINT) run ./...

audit: $(GO_SEC) ## Run audit
	$(GO_SEC) ./...

static: $(GO_STATIC) ## Run static analysis
	$(GO_STATIC) ./...

tidy: ## Tidy module
	$(GO_MOD) tidy

deps: tidy ## Install dependencies
	$(GO_MOD) download

update:
	$(GO_GET) -u=patch all

rpcw:: PROTO_FILES = \
	./api/service/proto/account.proto \
	./api/service/proto/auth.proto \
	./api/service/proto/response.proto \
	./api/service/proto/user.proto

test: ## Run tests
ifeq (, $(shell which docker))
	curl -fsSL https://get.docker.com -o get-docker.sh
	sh get-docker.sh
endif
	$(GO_TEST) $(BUILD_TAGS) $(TEST_FLAGS) ./...

cover: TEST_FLAGS := $(TEST_FLAGS) $(COVERAGE_FLAGS)
cover: test ## Run tests w/ coverage
	curl -fsSL https://codecov.io/bash | bash
	$(RM) $(COVERAGE_FILE)

build: ## Debug build
	$(GO_BUILD) -o $(OUT_EXE) $(BUILD_TAGS) -ldflags="$(LD_FLAGS) $(VER_FLAGS)" $(IN_EXE)
	$(GO_BUILD) -o $(OUT_CLI) $(BUILD_TAGS) -ldflags="$(LD_FLAGS) $(VER_FLAGS)" $(IN_CLI)

release: BUILD_TAGS := $(RELEASE_TAGS)
release: OUT_DIR := $(RELEASE_DIR)
release: LD_FLAGS := -s -w
release: CGO_ENABLED=0
release: build ## Production build

install: OUT_DIR := $(RELEASE_DIR)
install: release ## Install gothic
	mv $(OUT_EXE) $(GO_BIN)/$(EXE)
	mv $(OUT_CLI) $(GO_BIN)/$(CLI)

clean: ## Clean
	$(RM) -r $(BUILD_DIR)

# lint temp disable lint as we are getting false positives
all: vet test release ## Lint, vet, test, & release

image: ## Build the Docker image
	$(DOCKER) build .

gothic:  ## Start gothic
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) up -d gothic

envoy: ## Start envoy
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) up -d envoy

mysql: COMPOSE_FILE+= -f $(COMPDEV_FILE)
mysql: ## Start mysql
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) up -d mysql

pg: COMPOSE_FILE+= -f $(COMPDEV_FILE)
pg: ## Start postgres
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) up -d pg

db: mysql ## Start mysql db

sim: COMPOSE_FILE+= -f $(COMPDEV_FILE)
sim: COMPOSE_ENV:=./env/sim.env
sim: db envoy gothic

.PHONY: help fmt vet lint audit static tidy deps rpc test build \
		release install all image gothic envoy mysql pg db cover sim

.DEFAULT_GOAL := build
