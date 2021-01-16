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

$(GO_LINT):
	$(GO_GET) $(GO_LINT_REPO)

$(GO_SEC):
	$(GO_GET) $(GO_SEC_REPO)

$(GO_STATIC):
	$(GO_GET) $(GO_STATIC_REPO)

BUILD_DIR := build
DEBUG_DIR := $(BUILD_DIR)/debug
RELEASE_DIR := $(BUILD_DIR)/release
OUT_DIR := $(DEBUG_DIR)
OUT_EXE = -o $(OUT_DIR)/$(EXE)

PKG := github.com/jrapoport/gothic/conf
EXE := gothic

VERSION_NUM := $(shell git describe --abbrev=0 --tags 2> /dev/null)
ifeq (, $(VERSION_NUM))
	VERSION_NUM := 0.0.1
endif
BUILD_MN := $(shell git log -1 --format=%cd --date=format:'%m')
BUILD_YR := $(shell git log -1 --format=%cd --date=format:'%y%d')
BUILD_NUM := $(shell printf '%b%s' \\$(shell printf %o $(shell expr $(shell date +%m) + 64)) $(shell date +%y%d))
VER_PKG := $(PKG)
VER_FLAGS = -X '${VER_PKG}.Version=${VERSION_NUM}' -X '${VER_PKG}.Build=${BUILD_NUM}'
DEBUG_TAGS := -tags="debug"
RELEASE_TAGS := -tags="osusergo,netgo,release"
BUILD_TAGS := $(DEBUG_TAGS)

help: ## Show this help.
	echo $(BUILD_NUM)
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

fmt:
	$(GO_FMT) ./...

vet:
	$(GO_VET) ./...

lint: $(GO_LINT)
	$(GO_LINT) ./...

audit: $(GO_SEC)
	$(GO_SEC) ./...

static: $(GO_STATIC)
	$(GO_STATIC) ./...

tidy:
	$(GO_MOD) tidy

deps: tidy
	$(GO_MOD) download

rpc:
	$(GO_GEN) ./...

test:
	$(GO_TEST) $(BUILD_TAGS) ./...

build:
	$(GO_BUILD) $(OUT_EXE) $(BUILD_TAGS) -ldflags="$(LD_FLAGS) $(VER_FLAGS)" $(IN_EXE)

# RELEASE
release: BUILD_TAGS := $(RELEASE_TAGS)
release: OUT_DIR := $(RELEASE_DIR)
release: LD_FLAGS := -s -w
release: CGO_ENABLED=0
release: build

# INSTALL
install: OUT_EXE :=
install: OUT_CLI :=
install: GO_BUILD = $(GOINSTALL)
install: release

all: lint vet test release

image: ## Build the Docker image.
	docker build .

auth:
	docker-compose -f docker-compose.yaml up -d auth

envoy:
	docker-compose -f docker-compose.yaml up -d envoy

db:
	docker-compose -f docker-compose-dev.yaml up -d

.PHONY: help fmt vet lint audit static tidy deps rpc test \
		build release install all image auth envoy db

.DEFAULT_GOAL := build
