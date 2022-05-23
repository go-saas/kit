FROM golang:1.18 AS builder

COPY . /src
WORKDIR /src/gateway/apisix

RUN make -f ./Makefile build


FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/gateway/apisix/bin /app

WORKDIR /app

VOLUME /data/conf

CMD ["./go-runner","run", "-c", "/data/conf"]
