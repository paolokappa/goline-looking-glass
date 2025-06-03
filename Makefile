.PHONY: build install clean test docker

# Variables
BINARY_NAME=looking-glass
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Build application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) main.go

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Setup development environment
setup:
	@echo "Setting up development environment..."
	mkdir -p config public/images logs
	cp config/config.example.json config/config.json || true
	cp config/company.example.json config/company.json || true
	cp config/routers.example.json config/routers.json || true
	@echo "Setup complete! Edit config files as needed."

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean
