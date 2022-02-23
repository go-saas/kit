GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
BUF_VERSION=v1.0.0

.PHONY: init
# init env
init:
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/envoyproxy/protoc-gen-validate@v0.6.2
	go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@$(BUF_VERSION)
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@$(BUF_VERSION)

.PHONY: user
user:
	cd user && $(MAKE) all

.PHONY: saas
saas:
	cd saas && $(MAKE) all

all:
	make user
	make saas
	make api

.PHONY: api
# generate api proto
api:
	buf generate ./proto
