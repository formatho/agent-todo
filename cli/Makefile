.PHONY: build clean install test

BINARY_NAME=agent-todo
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Build for current platform
build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe .

# Install to GOPATH/bin
install:
	go install $(LDFLAGS) .

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Run linter and formatter
check: fmt

# Generate completion scripts
completions:
	@mkdir -p bin/completions
	./bin/$(BINARY_NAME) completion bash > bin/completions/$(BINARY_NAME).bash
	./bin/$(BINARY_NAME) completion zsh > bin/completions/$(BINARY_NAME).zsh
	./bin/$(BINARY_NAME) completion fish > bin/completions/$(BINARY_NAME).fish
	./bin/$(BINARY_NAME) completion powershell > bin/completions/$(BINARY_NAME).ps1

# Development build with race detection
dev:
	go build -race -o bin/$(BINARY_NAME) .
