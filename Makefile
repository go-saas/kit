GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
BUF_VERSION=v1.27.1
DIR=$(shell pwd)

SRV_PROTO_DIR = dtm event oidc user sys saas realtime gateway payment order product
PKG_PROTO_DIR = $(patsubst %/,%,$(shell cd pkg && ls -d */))
OTHER_PROTO_DIR = $(patsubst %/,%,$(shell cd proto && ls -d */))
THIRD_PARTY_PROTO_DIR = errors google lbs protoc-gen-openapiv2 validate

.PHONY: link
# link proto
link:
	for d in $(SRV_PROTO_DIR); do \
  		ln -sfn $(DIR)/$$d $(DIR)/buf/$$d; \
  	done
	for d in $(PKG_PROTO_DIR); do \
		ln -sfn $(DIR)/pkg/$$d $(DIR)/buf/$$d; \
    done
	for d in $(OTHER_PROTO_DIR); do \
		ln -sfn $(DIR)/proto/$$d $(DIR)/buf/$$d; \
    done

.PHONY: init
# init env
init:
	make link
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-saas/kit/cmd/protoc-gen-go-grpc-proxy@c2ded75bd3ee9f1229e50d7141966ecbde39a84f
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-saas/kit/cmd/protoc-gen-go-errors-i18n/v2@c2ded75bd3ee9f1229e50d7141966ecbde39a84f
	go install github.com/envoyproxy/protoc-gen-validate@v1.0.2
	go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@$(BUF_VERSION)
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@$(BUF_VERSION)

.PHONY: user
# user
user:
	cd user && $(MAKE) all

.PHONY: saas
# saas
saas:
	cd saas && $(MAKE) all

.PHONY: realtime
# realtime
realtime:
	cd realtime && $(MAKE) all

.PHONY: sys
# sys
sys:
	cd sys && $(MAKE) all

.PHONY: payment
# payment
payment:
	cd payment && $(MAKE) all

.PHONY: order
# order
order:
	cd order && $(MAKE) all

.PHONY: product
# product
product:
	cd product && $(MAKE) all

.PHONY: apisix
apisix:
	cd gateway/apisix && $(MAKE) all

.PHONY: api
# generate api proto
api:
	buf generate --exclude-path ./buf/errors --exclude-path ./buf/google --exclude-path ./buf/lbs --exclude-path ./buf/protoc-gen-openapiv2 --exclude-path ./buf/validate
	cd user && $(MAKE) api
	cd saas && $(MAKE) api
	cd sys && $(MAKE) api
	cd realtime && $(MAKE) api
	cd payment && $(MAKE) api
	cd order && $(MAKE) api
	cd product && $(MAKE) api

.PHONY: build
build:
	cd user && $(MAKE) build
	cd saas && $(MAKE) build
	cd sys && $(MAKE) build
	cd realtime && $(MAKE) build
	cd gateway/apisix && $(MAKE) build
	cd payment && $(MAKE) build
	cd order && $(MAKE) build
	cd product && $(MAKE) build

.PHONY: all
all:
	go mod tidy
	make api
	make user
	make saas
	make sys
	make realtime
	make apisix
	make payment
	make order
	make product

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help