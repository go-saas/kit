GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
BUF_VERSION=v1.5.0
DIR=$(shell pwd)



.PHONY: link
# link proto
link:
	mkdir -p buf
	ln -sfn $(DIR)/user $(DIR)/buf/user
	ln -sfn $(DIR)/sys $(DIR)/buf/sys
	ln -sfn $(DIR)/saas $(DIR)/buf/saas

.PHONY: init
# init env
init:
	make link
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/goxiaoy/go-saas-kit/cmd/protoc-gen-go-grpc-proxy@main
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/goxiaoy/go-saas-kit/cmd/protoc-gen-go-errors-i18n/v2@main
	go install github.com/envoyproxy/protoc-gen-validate@v0.6.7
	go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@$(BUF_VERSION)
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@$(BUF_VERSION)

.PHONY: user
user:
	cd user && $(MAKE) all

.PHONY: saas
saas:
	cd saas && $(MAKE) all

.PHONY: sys
sys:
	cd sys && $(MAKE) all

.PHONY: apisix
apisix:
	cd gateway/apisix && $(MAKE) all


.PHONY: api
# generate api proto
api:
	buf generate

.PHONY: generate
# generate
generate:
	go generate ./pkg/...

.PHONY: examples
# generate
examples:
	cd examples/monolithic && $(MAKE) all

.PHONY: build
build:
	cd user && $(MAKE) build
	cd saas && $(MAKE) build
	cd sys && $(MAKE) build
	cd gateway/apisix && $(MAKE) build
	cd examples/monolithic && $(MAKE) build


.PHONY: all
all:
	go mod tidy
	make init
	make generate
	make api
	make user
	make saas
	make sys
	make apisix
	make examples