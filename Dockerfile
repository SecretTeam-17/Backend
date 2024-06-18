FROM golang:1.22.3-alpine as builder

ENV CGO_ENABLED=1

ENV CONFIG_PATH=./configs/dev.yaml

RUN apk add --no-cache \
    gcc \
    musl-dev

WORKDIR /go/src/sitterserver

RUN mkdir -p /go/src/sitterserver/cert

COPY . .

COPY /etc/letsencrypt/live/sf-hackathon.xyz/fullchain.pem /etc/letsencrypt/live/sf-hackathon.xyz/privkey.pem ./cert/

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

RUN mkdir -p /root/storage /root/configs /root/internal/templates /root/cert

COPY --from=builder /go/src/sitterserver/configs ./configs

COPY --from=builder /go/src/sitterserver/internal/templates ./internal/templates

COPY --from=builder /go/src/sitterserver/cert ./cert

COPY --from=builder /go/bin/server .

ENTRYPOINT /root/server

LABEL Name=PetsittersAPIgameserver Version=1.10.0

EXPOSE 80

EXPOSE 443