UNAME_S := $(shell uname -s)

GO_PATH := $(shell go env GOPATH)
GO_BIN := $(GO_PATH)/bin

ifeq ($(UNAME_S),Linux)
	PLUGIN_BIN := ./bin/linux
else ifeq ($(UNAME_S),Darwin)
	PLUGIN_BIN := ./bin/darwin
endif

BASE_DIR := .
BUILD_DIR := build

GRPC_DIR := $(BUILD_DIR)/grpc
GRPC_PLUGIN_WEB := $(PLUGIN_BIN)/protoc-gen-grpc-web
GRPC_PLUGIN_CSHARP := $(PLUGIN_BIN)/grpc_csharp_plugin
GRPC_PREFIX_OPT = $(if $(GRPC_PREFIX), --go_opt=module=$(GRPC_PREFIX),)
GRPC_RELATIVE_OPT = $(if $(GRPC_RELATIVE), --go_opt=paths=source_relative,)

PROTOC := protoc
PROTO_WILDCARD ?= "*.proto"
PROTO_FIND = find $(BASE_DIR) -type f \( -iname $(PROTO_WILDCARD) ! -iname "_*" \)
ifneq (BASE_DIR,.)
	PROTO_FIND += $(if $(PROTO_RELATIVE),| sed -n 's|^${BASE_DIR}|.|p',)
endif
PROTO_FILES = $(shell $(PROTO_FIND))
PROTO_INCLUDES = -I=$(BASE_DIR)
PROTOC_GEN := $(GO_BIN)/protoc-gen-go
PROTOC_REPO := github.com/golang/protobuf/protoc-gen-go

$(PROTOC_GEN):
ifeq (, $(shell which protoc))
ifeq ($(UNAME_S),Linux)
	apt install -y protobuf-compiler
endif
ifeq ($(UNAME_S),Darwin)
	brew install protobuf
endif
endif
	$(GOGET) $(PROTOC_REPO)

GRPC_RPC_DIR = $(GRPC_DIR)/rpc
rpc:: $(PROTOC_GEN) ## Protobuf gRPC
	mkdir -p $(GRPC_RPC_DIR)
	$(PROTOC) $(PROTO_INCLUDES) \
	--go_out=plugins=grpc:$(GRPC_RPC_DIR) \
	$(GRPC_RELATIVE_OPT) $(GRPC_PREFIX_OPT) \
	$(PROTO_FILES)

GRPC_WEB_DIR = $(GRPC_DIR)/web
rpcw:: $(PROTOC_GEN) ## Protobuf gRPC-Web
	mkdir -p $(GRPC_WEB_DIR)
	$(PROTOC) $(PROTO_INCLUDES) \
	--js_out=import_style=commonjs:$(GRPC_WEB_DIR) \
	--grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:$(GRPC_WEB_DIR) \
	--plugin=protoc-gen-grpc-web=$(GRPC_PLUGIN_WEB) $(PROTO_FILES)

GRPC_CSHARP = $(GRPC_DIR)/csharp
rpc-cs:: $(PROTOC_GEN) ## Protobuf gRPC C#
	mkdir -p $(GRPC_CSHARP)
	$(PROTOC) $(PROTO_INCLUDES) \
	--plugin=protoc-gen-grpc=$(GRPC_PLUGIN_CSHARP) \
	--csharp_out=./$(GRPC_CSHARP) \
	--grpc_out=./$(GRPC_CSHARP) \
	--grpc_opt=no_server \
	$(PROTO_FILES)

grpc: rpc rpcw # Protobuf gRPC All

.PHONY: rpc rpcw rpc-cs grpc
