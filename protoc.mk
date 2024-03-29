OS := $(shell uname -s)

GO_PATH := $(shell go env GOPATH)
GO_BIN := $(GO_PATH)/bin
GO := go
GO_GET := $(GO) get -v

BASE_DIR := .
BUILD_DIR := build

ifeq ($(OS),Linux)
	PLUGIN_BIN := ./bin/linux
else ifeq ($(OS),Darwin)
	PLUGIN_BIN := ./bin/macos
endif

GRPC_DIR := $(BUILD_DIR)/grpc
GRPC_PLUGIN_WEB := $(PLUGIN_BIN)/protoc-gen-grpc-web
GRPC_PLUGIN_CSHARP := $(PLUGIN_BIN)/grpc_csharp_plugin
GRPC_PREFIX_OPT = $(if $(GRPC_PREFIX), --go_opt=module=$(GRPC_PREFIX) --go-grpc_opt=module=$(GRPC_PREFIX),)
GRPC_RELATIVE_OPT = $(if $(GRPC_RELATIVE), --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative,)

PROTOC = PATH="${GO_BIN}:${PATH}" $(shell which protoc)
PROTO_WILDCARD ?= "*.proto"
PROTO_FIND = find $(BASE_DIR) -type f \( -iname $(PROTO_WILDCARD) ! -iname "_*" \)
ifneq (BASE_DIR,.)
	PROTO_FIND += $(if $(PROTO_RELATIVE),| sed -n 's|^${BASE_DIR}|.|p',)
endif
PROTO_FILES = $(shell $(PROTO_FIND))
PROTO_INCLUDES = -I=$(BASE_DIR)

PROTOC_GO := $(GO_BIN)/protoc-gen-go
PROTOC_REPO_GO := google.golang.org/protobuf/cmd/protoc-gen-go@latest
$(PROTOC_GO):
	$(GO_GET) $(PROTOC_GO)

PROTOC_GRPC := $(GO_BIN)/protoc-gen-go-grpc
PROTOC_REPO_GRPC := google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
$(PROTOC_GRPC):
	$(GO_GET) $(PROTOC_REPO_GRPC)

proto: $(PROTOC_GO) $(PROTOC_GRPC)
ifeq (, $(PROTOC))
ifeq ($(OS),Linux)
	apt install -y protobuf-compiler
endif
ifeq ($(OS),Darwin)
	brew install protobuf
endif
endif

GRPC_RPC_DIR = $(GRPC_DIR)/rpc
rpc:: proto ## Protobuf gRPC
	$(RM) -rf $(GRPC_RPC_DIR)
	mkdir -p $(GRPC_RPC_DIR)
	$(PROTOC) $(PROTO_INCLUDES) \
	--go_out=$(GRPC_RPC_DIR) \
	--go-grpc_out=$(GRPC_RPC_DIR) \
	$(GRPC_RELATIVE_OPT) $(GRPC_PREFIX_OPT) \
	$(PROTO_FILES)

GRPC_WEB_DIR = $(GRPC_DIR)/web
rpcw:: proto ## Protobuf gRPC-Web
	$(RM) -rf $(GRPC_WEB_DIR)
	mkdir -p $(GRPC_WEB_DIR)
	$(PROTOC) $(PROTO_INCLUDES) \
	--js_out=import_style=commonjs:$(GRPC_WEB_DIR) \
	--grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:$(GRPC_WEB_DIR) \
	--plugin=protoc-gen-grpc-web=$(GRPC_PLUGIN_WEB) $(PROTO_FILES)

GRPC_CSHARP = $(GRPC_DIR)/csharp
rpc-cs:: proto ## Protobuf gRPC C#
	mkdir -p $(GRPC_CSHARP)
	$(PROTOC) $(PROTO_INCLUDES) \
	--plugin=protoc-gen-grpc=$(GRPC_PLUGIN_CSHARP) \
	--csharp_out=./$(GRPC_CSHARP) \
	--grpc_out=./$(GRPC_CSHARP) \
	--grpc_opt=no_server \
	$(PROTO_FILES)

grpc: rpc rpcw # Protobuf gRPC All

.PHONY: proto rpc rpcw rpc-cs grpc
