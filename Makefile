GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
BUF_VERSION=v1.0.0

.PHONY: init
# init env
init:
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
