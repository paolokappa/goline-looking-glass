FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o looking-glass main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/looking-glass .
COPY --from=builder /app/public ./public/

RUN mkdir -p logs config
VOLUME ["/root/config", "/root/logs"]

EXPOSE 3002
CMD ["./looking-glass"]
