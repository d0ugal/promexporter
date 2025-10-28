# Random Exporter Example

This is a comprehensive example exporter built using the `promexporter` library. It demonstrates all the features and capabilities of promexporter, including:

- **Configuration Management**: Environment-based and YAML-based configuration
- **Metrics Collection**: All Prometheus metric types (Counter, Gauge, Histogram, Summary, Info)
- **OpenTelemetry Tracing**: Complete tracing integration with spans, events, and error tracking
- **Logging**: Structured logging with configurable levels and formats
- **HTTP Server**: Built-in metrics endpoint and health checks
- **Graceful Shutdown**: Signal handling and clean shutdown
- **Error Handling**: Random error generation for testing error scenarios

## Features Demonstrated

### 1. Configuration Management
- Environment variable configuration loading
- YAML configuration file support
- Configuration validation and defaults
- Configuration display (`--show-config` flag)

### 2. Metrics Types
- **Counters**: `random_counter_total`, `random_counter_rate_total`
- **Gauges**: `random_gauge`, `random_temperature_celsius`, `random_memory_usage_bytes`
- **Histograms**: `random_latency_seconds`, `random_response_time_seconds`
- **Summaries**: `random_duration_seconds`, `random_processing_time_seconds`
- **Info**: `random_info`

### 3. OpenTelemetry Tracing
- Service name configuration
- OTLP endpoint configuration
- Span creation for collection cycles
- Sub-spans for different metric types
- Event logging within spans
- Error recording with attributes
- Context propagation

### 4. Logging
- Structured JSON logging
- Configurable log levels (debug, info, warn, error)
- Configurable log formats (json, text)
- Contextual logging with fields

### 5. HTTP Server
- Metrics endpoint at `/metrics`
- Health check endpoint at `/health`
- Version information endpoint
- Configurable server timeouts
- Request tracing middleware

## Usage

### Command Line Options

```bash
# Show version information
./random-exporter --version

# Show loaded configuration
./random-exporter --show-config

# Load configuration from environment variables
./random-exporter --config-from-env

# Specify configuration file
./random-exporter --config /path/to/config.yaml
```

### Environment Variables

#### Server Configuration
- `SERVER_HOST`: Server host (default: "0.0.0.0")
- `SERVER_PORT`: Server port (default: 8080)

#### Logging Configuration
- `LOG_LEVEL`: Log level - debug, info, warn, error (default: "info")
- `LOG_FORMAT`: Log format - json, text (default: "json")

#### Tracing Configuration
- `TRACING_ENABLED`: Enable tracing - true, false (default: false)
- `TRACING_SERVICE_NAME`: Service name for traces (default: "promexporter")
- `TRACING_ENDPOINT`: OTLP endpoint URL (default: "")

#### Random Exporter Configuration
- `RANDOM_COLLECTION_INTERVAL`: Collection interval (default: "10s")
- `RANDOM_METRIC_COUNT`: Number of metrics to generate per collection (default: 20)
- `RANDOM_ENABLE_ERRORS`: Enable random errors (default: false)
- `RANDOM_ERROR_PROBABILITY`: Error probability 0.0-1.0 (default: 0.1)

### Example Configuration

```bash
# Enable tracing
export TRACING_ENABLED=true
export TRACING_SERVICE_NAME=random-exporter-demo
export TRACING_ENDPOINT=http://localhost:4318/v1/traces

# Configure collection
export RANDOM_COLLECTION_INTERVAL=5s
export RANDOM_METRIC_COUNT=50
export RANDOM_ENABLE_ERRORS=true
export RANDOM_ERROR_PROBABILITY=0.05

# Run the exporter
./random-exporter --config-from-env
```

## Building and Running

### Prerequisites
- Go 1.21 or later
- Access to the promexporter library

### Build
```bash
cd examples/random-exporter
go mod init random-exporter
go mod tidy
go build -o random-exporter main.go
```

### Run
```bash
# Basic run with defaults
./random-exporter

# Run with tracing enabled
TRACING_ENABLED=true TRACING_SERVICE_NAME=random-exporter-demo ./random-exporter

# Run with custom configuration
RANDOM_COLLECTION_INTERVAL=5s RANDOM_METRIC_COUNT=100 ./random-exporter
```

## Integration with Observability Stack

This example integrates well with:
- **Prometheus**: For metrics collection
- **Grafana**: For metrics visualization  
- **Tempo/Jaeger**: For trace storage and visualization
- **Loki**: For log aggregation
- **AlertManager**: For alerting on metrics

## Tracing Integration

The random exporter demonstrates comprehensive tracing integration:

### Span Hierarchy
```
random-collector (collect-metrics)
├── random-collector (generate-counters)
├── random-collector (generate-gauges)
├── random-collector (generate-histograms)
├── random-collector (generate-summaries)
└── random-collector (generate-info)
```

### Events and Attributes
- Collection start/end events
- Metric generation counts
- Error events with error types
- Service and operation attributes

### Error Tracking
- Random error generation for testing
- Error recording with context
- Error probability configuration

## Metrics Generated

### Counter Metrics
- `random_counter_total`: Incrementing counters by service and region
- `random_counter_rate_total`: Counters with varying increment rates

### Gauge Metrics
- `random_gauge`: Fluctuating values by instance and type
- `random_temperature_celsius`: Temperature readings by sensor and location
- `random_memory_usage_bytes`: Memory usage by process and type

### Histogram Metrics
- `random_latency_seconds`: Request latency by service and endpoint
- `random_response_time_seconds`: Response times by method and status

### Summary Metrics
- `random_duration_seconds`: Operation duration by operation and priority
- `random_processing_time_seconds`: Processing times by task and worker

### Info Metrics
- `random_info`: Static information about the exporter

## Testing

### Manual Testing
```bash
# Test metrics endpoint
curl http://localhost:8080/metrics

# Test health endpoint
curl http://localhost:8080/health

# Test version endpoint
curl http://localhost:8080/version
```

### Tracing Testing
1. Start Tempo or another OTLP-compatible backend
2. Configure tracing environment variables
3. Run the exporter
4. Check traces in your tracing backend

### Error Testing
```bash
# Enable random errors
export RANDOM_ENABLE_ERRORS=true
export RANDOM_ERROR_PROBABILITY=0.5
./random-exporter
```

## Integration with Observability Stack

This example integrates well with:
- **Prometheus**: For metrics collection
- **Grafana**: For metrics visualization
- **Tempo/Jaeger**: For trace storage and visualization
- **Loki**: For log aggregation
- **AlertManager**: For alerting on metrics

## Code Structure

The example demonstrates proper Go project structure:
- Clear separation of concerns
- Interface-based design
- Error handling patterns
- Context propagation
- Resource management
- Graceful shutdown

This serves as a comprehensive reference for building exporters with the promexporter library.
