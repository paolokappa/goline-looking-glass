.PHONY: build install clean test docker

# Variables
BINARY_NAME=looking-glass
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Build application
build:
	@echo "?? Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) main.go

# Install dependencies
deps:
	@echo "?? Installing dependencies..."
	go mod tidy
	go mod download

# Run tests
test:
	@echo "?? Running tests..."
	go test -v ./...

# Build for multiple platforms
build-all:
	@echo "?? Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go

# Docker build
docker:
	@echo "?? Building Docker image..."
	docker build -t looking-glass:$(VERSION) .
	docker tag looking-glass:$(VERSION) looking-glass:latest

# Install system-wide
install: build
	@echo "?? Installing system-wide..."
	sudo ./scripts/install.sh

# Clean build artifacts
clean:
	@echo "?? Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean

# Setup development environment
setup:
	@echo "?? Setting up development environment..."
	mkdir -p config public/images logs
	cp config/config.example.json config/config.json || true
	cp config/company.example.json config/company.json || true
	cp config/routers.example.json config/routers.json || true
	@echo "? Setup complete! Edit config files as needed."

# Run locally
run: build
	@echo "?? Starting Looking Glass..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Test coverage
test-coverage:
	@echo "?? Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Lint code
lint:
	@echo "?? Linting code..."
	golangci-lint run || echo "Install golangci-lint for better linting"
