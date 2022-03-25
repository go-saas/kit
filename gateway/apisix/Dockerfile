FROM golang:1.18 AS builder

COPY . /src
WORKDIR /src/gateway/apisix

RUN GOPROXY=https://goproxy.cn make -f ./Makefile build


FROM apache/apisix:2.13.0-centos

COPY --from=builder /src/gateway/apisix/bin /app
