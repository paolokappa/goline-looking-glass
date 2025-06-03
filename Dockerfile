# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o looking-glass main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl

WORKDIR /root/

# Copy binary and config
COPY --from=builder /app/looking-glass .
COPY --from=builder /app/config.json .
COPY --from=builder /app/public ./public/

# Create logs directory
RUN mkdir -p logs config

# Expose port
EXPOSE 3002

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:3002/api/health || exit 1

# Run the application
CMD ["./looking-glass"]
