# Makefile for promexporter library

.PHONY: help test lint fmt build clean lint-only

# Docker image versions
GOLANGCI_LINT_VERSION := v2.5.0

# Default target
help:
	@echo "Available targets:"
	@echo "  test      - Run tests"
	@echo "  lint      - Format code and run golangci-lint"
	@echo "  fmt       - Format code using golangci-lint"
	@echo "  lint-only - Run golangci-lint without formatting"
	@echo "  build     - Build the library"
	@echo "  clean     - Clean build artifacts"

# Run tests
test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./... || true

# Format code using golangci-lint formatters (faster than separate tools)
fmt:
	docker run --rm \
		-u "$(shell id -u):$(shell id -g)" \
		-e GOCACHE=/tmp/go-cache \
		-e GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache \
		-v "$(PWD):/app" \
		-v "$(HOME)/.cache:/home/cache" \
		-w /app \
		golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) \
		golangci-lint run --fix

# Run golangci-lint (formats first, then lints)
lint:
	docker run --rm \
		-u "$(shell id -u):$(shell id -g)" \
		-e GOCACHE=/tmp/go-cache \
		-e GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache \
		-v "$(PWD):/app" \
		-v "$(HOME)/.cache:/home/cache" \
		-w /app \
		golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) \
		golangci-lint run --fix

# Run only linting without formatting
lint-only:
	docker run --rm \
		-u "$(shell id -u):$(shell id -g)" \
		-e GOCACHE=/tmp/go-cache \
		-e GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache \
		-v "$(PWD):/app" \
		-v "$(HOME)/.cache:/home/cache" \
		-w /app \
		golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) \
		golangci-lint run

# Build
build:
	go build ./...

# Clean build artifacts
clean:
	go clean
	rm -f coverage.txt
