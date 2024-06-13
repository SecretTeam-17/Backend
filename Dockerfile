FROM golang:1.22.3-alpine as builder
# FROM huecker.io/library/golang:1.22.3-alpine as builder

ENV CGO_ENABLED=1

ENV CONFIG_PATH=./configs/dev.yaml

# ENV DB_PATH=mongodb://95.164.3.230:10515/

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
# FROM huecker.io/library/alpine:latest as runner

ENV CGO_ENABLED=1

ENV CONFIG_PATH=/root/configs/dev.yaml

# ENV DB_PATH=mongodb://95.164.3.230:10515/

RUN apk add --no-cache \
    ca-certificates \
    gcc \
    musl-dev

STOPSIGNAL SIGTERM

WORKDIR /root

RUN mkdir -p /root/storage /root/configs /root/internal/templates

COPY --from=builder /go/src/sitterserver/configs ./configs

COPY --from=builder /go/src/sitterserver/internal/templates ./internal/templates

COPY --from=builder /go/bin/server .

ENTRYPOINT /root/server

LABEL Name=petsittersgameserver Version=1.9.0

EXPOSE 8083