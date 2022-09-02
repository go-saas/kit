GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)


.PHONY: api
# generate api proto
api:
	cd .. && buf generate --path ./buf/realtime --template ./realtime/buf.gen.yaml

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -buildvcs=false -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...


.PHONY: all
# generate all
all:
	make api;


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
