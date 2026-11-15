# --- DEV (Has Go, Source Code, and Swaggo) ---
FROM golang:1.26.8-alpine AS dev
RUN apk add --no-cache git build-base
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Pre-download and compile swag into the dev image layer
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

# --- DEV BUILDER ---
FROM dev AS builder
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/server

# --- PROD (Tiny, No Go, No Docs) ---
FROM alpine:latest AS prod
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080