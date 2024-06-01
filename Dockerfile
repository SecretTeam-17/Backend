FROM golang:1.22.3-alpine as builder

ENV CGO_ENABLED=1

ENV CONFIG_PATH=./configs/dev.yaml

RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

WORKDIR /go/src/sitterserver

COPY . .

RUN go mod download

RUN go mod tidy

RUN go install -ldflags='-s -w -extldflags "-static"' ./cmd/main.go

RUN mv /go/bin/main /go/bin/server

FROM alpine:latest as runner

ENV CGO_ENABLED=1

ENV CONFIG_PATH=/root/configs/dev.yaml

RUN apk add --no-cache \
    ca-certificates \
    gcc \
    musl-dev

STOPSIGNAL SIGTERM

WORKDIR /root

RUN mkdir -p /root/storage /root/configs

COPY --from=builder /go/src/sitterserver/configs ./configs

COPY --from=builder /go/bin/server .

ENTRYPOINT /root/server

LABEL Name=petsittersgameserver Version=1.3.0

EXPOSE 8082