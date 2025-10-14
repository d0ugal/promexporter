# Makefile for promexporter library

.PHONY: help test lint fmt build clean

# Default target
help:
	@echo "Available targets:"
	@echo "  test    - Run tests"
	@echo "  lint    - Run linter"
	@echo "  fmt     - Format code"
	@echo "  build   - Build the library"
	@echo "  clean   - Clean build artifacts"

# Test
test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Lint
lint:
	golangci-lint run

# Format
fmt:
	go fmt ./...
	goimports -w .

# Build
build:
	go build ./...

# Clean
clean:
	go clean
	rm -f coverage.txt
