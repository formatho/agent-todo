# Makefile for agent-todo CLI

# Variables
BINARY_NAME=agent-todo
CLI_PATH=./cli
BUILD_DIR=./bin
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILT_BY=$(shell whoami)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE) -X main.BuiltBy=$(BUILTBy)"

# Platforms for cross-compilation
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

.PHONY: all build clean test install release

all: build

# Build for current platform
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@cd $(CLI_PATH) && go build $(LDFLAGS) -o ../$(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Built: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for all platforms
release:
	@echo "Building release binaries..."
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		OS=$${platform%/*}; \
		ARCH=$${platform#*/}; \
		echo "Building for $$OS/$$ARCH..."; \
		cd $(CLI_PATH) && GOOS=$$OS GOARCH=$$ARCH go build $(LDFLAGS) -o ../$(BUILD_DIR)/$(BINARY_NAME)-$$OS-$$ARCH .; \
		cd ..; \
	done
	@echo "Release binaries built in $(BUILD_DIR)/"

# Install locally
install:
	@echo "Installing $(BINARY_NAME)..."
	@cd $(CLI_PATH) && go install $(LDFLAGS) .
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@cd $(CLI_PATH) && go clean
	@echo "Cleaned"

# Run tests
test:
	@echo "Running tests..."
	@cd $(CLI_PATH) && go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	@cd $(CLI_PATH) && go fmt ./...
	@echo "Formatted"

# Run linter
lint:
	@echo "Running linter..."
	@cd $(CLI_PATH) && golangci-lint run ./...
	@echo "Linting complete"

# Build with version info
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"
	@echo "Built by: $(BUILT_BY)"

# Show help
help:
	@echo "Available targets:"
	@echo "  build    - Build for current platform"
	@echo "  release  - Build for all platforms"
	@echo "  install  - Install locally"
	@echo "  clean    - Clean build artifacts"
	@echo "  test     - Run tests"
	@echo "  fmt      - Format code"
	@echo "  lint     - Run linter"
	@echo "  version  - Show version info"
	@echo "  help     - Show this help message"
