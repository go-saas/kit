FROM golang:1.18 AS builder

COPY . /src
WORKDIR /src/gateway/apisix

RUN make -f ./Makefile build


FROM apache/apisix:2.14.0-centos

COPY --from=builder /src/gateway/apisix/bin /app
