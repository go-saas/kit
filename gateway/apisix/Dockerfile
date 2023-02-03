FROM golang:1.19 AS builder

COPY . /src
WORKDIR /src/gateway/apisix

RUN make -f ./Makefile build


FROM apache/apisix:2.15.1-centos

COPY --from=builder /src/gateway/apisix/bin /app
