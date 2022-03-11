FROM golang:1.18rc1 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make -f ./sys/Makefile build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app

WORKDIR /app

VOLUME /data/conf

CMD ["./sys", "-conf", "/data/conf"]
