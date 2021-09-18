export GO111MODULE=on
export GOPRIVATE=github.com/eugeneradionov*

BIN_NAME := $(or $(PROJECT_NAME),'froxy')
PKG_LIST := $(shell go list ./...)

.PHONY: dep lint

all: dep lint

dep: ## Download required dependencies
	go mod vendor
	go mod tidy

lint: ## Lint files
	golangci-lint run -c .golangci.yml

test: dep ## Run unit tests
	go test -cover -race -count=1 ${PKG_LIST}

build: dep ## Build the binary file
	go build -o ./bin/${BIN_NAME} -a -ldflags '-w -extldflags "-static"' ./cmd
