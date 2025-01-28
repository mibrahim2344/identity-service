# Build stage
FROM golang:1.22.0-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/identity

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy binary and config from builder
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config
COPY --from=builder /app/docs/swagger.json ./docs/swagger.json

EXPOSE 8080

CMD ["./main"]
