# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOFLAGS=-trimpath \
    GOPROXY=https://proxy.golang.org,direct
    
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine:3.20
WORKDIR /srv

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server /usr/local/bin/server

EXPOSE 8080
ENV PORT=8080

ENTRYPOINT ["/usr/local/bin/server"]
