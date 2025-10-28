# Makefile for promexporter library

.PHONY: help test lint fmt build clean lint-only dev-tag

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
	@echo "  dev-tag  - Generate dev tag for Docker image"

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

# Generate dev tag for Docker image
dev-tag:
	@SHORT_SHA=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	LAST_TAG=$$(git describe --tags --abbrev=0 --match="v[0-9]*.[0-9]*.[0-9]*" 2>/dev/null || echo ""); \
	if [ -z "$$LAST_TAG" ]; then \
		VERSION="0.0.0"; \
		COMMIT_COUNT=$$(git rev-list --count HEAD); \
	else \
		VERSION=$${LAST_TAG#v}; \
		COMMIT_COUNT=$$(git rev-list --count $${LAST_TAG}..HEAD); \
	fi; \
	DEV_TAG="v$${VERSION}-dev.$${COMMIT_COUNT}.$${SHORT_SHA}"; \
	echo "$$DEV_TAG"
