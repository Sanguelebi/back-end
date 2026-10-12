FROM golang:1.26.8-alpine AS builder
RUN apk add --no-cache git build-base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Building the binary
# - GOOS=linux ensures the binary targets Linux
# - ldflags "-s -w" strips debugging information to reduce binary size.
RUN GOOS=linux go build -ldflags="-s -w" -o main ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Copy the compiled the binary from the builder stage.
COPY --from=builder /app/main .
# The container listens on port 8080 at runtime. 
EXPOSE 8080
# Set Gin to production mode.
ENV GIN_MODE=release

# Run the binary
CMD ["./main"]
