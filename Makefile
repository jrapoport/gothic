.PHONY: all build deps image lint migrate test vet rpc db envoy auth

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: lint vet test build ## Run the tests and build the binary.

build: ## Build the binary.
	go build -ldflags "-X github.com/jrapoport/gothic/cmd.Version=`git rev-parse HEAD`"

rpc:
	@go generate ./...

deps: ## Install dependencies.
	@go get -u github.com/golang/protobuf/protoc-gen-go
	@go get -u golang.org/x/lint/golint
	@go mod download

image: ## Build the Docker image.
	docker build .

lint: ## Lint the code.
	golint ./...

test: ## Run tests.
	go test -p 1 -v ./...

vet: # Vet the code
	go vet ./...

auth:
	docker-compose -f docker-compose.yaml up -d auth

envoy:
	docker-compose -f docker-compose.yaml up -d envoy

db:
	docker-compose -f docker-compose-dev.yaml up -d