# Random Exporter Makefile

.PHONY: build run test clean help

# Default target
all: build

# Build the random exporter
build:
	@echo "Building random-exporter..."
	go mod tidy
	go build -o random-exporter main.go
	@echo "Build complete!"

# Run the random exporter locally
run: build
	@echo "Running random-exporter..."
	./random-exporter

# Run with tracing enabled
run-tracing: build
	@echo "Running random-exporter with tracing..."
	LOG_LEVEL=debug TRACING_ENABLED=true TRACING_SERVICE_NAME=random-exporter-demo ./random-exporter

# Run with debug logging
run-debug: build
	@echo "Running random-exporter with debug logging..."
	LOG_LEVEL=debug ./random-exporter

# Run with custom configuration
run-custom: build
	@echo "Running random-exporter with custom config..."
	SERVER_PORT=9999 RANDOM_COLLECTION_INTERVAL=3s RANDOM_METRIC_COUNT=100 RANDOM_ENABLE_ERRORS=true ./random-exporter

# Test the application
test:
	@echo "Running tests..."
	go test -v ./...

# Show configuration
show-config: build
	@echo "Showing configuration..."
	./random-exporter --show-config

# Show version
version: build
	@echo "Showing version..."
	./random-exporter --version

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f random-exporter
	go clean

# Test endpoints
test-endpoints:
	@echo "Testing endpoints..."
	@echo "Metrics endpoint:"
	@curl -s http://localhost:8080/metrics | head -20
	@echo "\nHealth endpoint:"
	@curl -s http://localhost:8080/health
	@echo "\nVersion endpoint:"
	@curl -s http://localhost:8080/version

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the random exporter"
	@echo "  run            - Run the random exporter locally"
	@echo "  run-tracing    - Run with tracing enabled"
	@echo "  run-debug      - Run with debug logging"
	@echo "  run-custom     - Run with custom configuration"
	@echo "  test           - Run tests"
	@echo "  show-config    - Show loaded configuration"
	@echo "  version        - Show version information"
	@echo "  clean          - Clean build artifacts"
	@echo "  test-endpoints - Test HTTP endpoints"
	@echo "  help           - Show this help"
