# Promexporter

A Go library for building Prometheus exporters with common functionality.

## Features

- **Application Bootstrap**: Simple builder pattern for setting up exporters
- **HTTP Server**: Gin-based server with standard routes (`/`, `/metrics`, `/health`)
- **Configuration Management**: YAML file and environment variable support with auto-detection
- **Structured Logging**: slog-based logging with configurable levels and formats
- **Metrics Registry**: Prometheus metrics with UI metadata tracking
- **Web Dashboard**: Modern, responsive HTML dashboard for all exporters
- **Graceful Shutdown**: Proper signal handling and cleanup
- **OpenTelemetry Tracing**: Optional distributed tracing support with OTLP export

## Quick Start

```go
package main

import (
    "github.com/d0ugal/promexporter/app"
    "github.com/d0ugal/promexporter/config"
    "github.com/d0ugal/promexporter/metrics"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig("config.yaml", false)
    if err != nil {
        log.Fatal(err)
    }

    // Create metrics registry
    metricsRegistry := metrics.NewRegistry()

    // Create and run application
    app := app.New("my-exporter").
        WithConfig(cfg).
        WithMetrics(metricsRegistry).
        Build()

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## Configuration

The library supports both YAML configuration files and environment variables:

### YAML Configuration

```yaml
server:
  host: "0.0.0.0"
  port: 8080

logging:
  level: "info"
  format: "json"

metrics:
  collection:
    default_interval: "30s"

tracing:
  enabled: false  # Set to true to enable tracing
  service_name: "my-exporter"
  endpoint: "http://localhost:4318/v1/traces"
  headers:
    authorization: "Bearer your-token"
```

### Environment Variables

```bash
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
LOG_LEVEL=info
LOG_FORMAT=json
METRICS_DEFAULT_INTERVAL=30s
TRACING_ENABLED=false
TRACING_SERVICE_NAME=my-exporter
TRACING_ENDPOINT=http://localhost:4318/v1/traces
```

## Tracing Support

The library includes optional OpenTelemetry tracing support for distributed tracing. When enabled, it automatically traces HTTP requests and provides utilities for tracing collector operations.

### Using Tracing in Collectors

```go
func (c *MyCollector) Start(ctx context.Context) {
    go c.run(ctx)
}

func (c *MyCollector) run(ctx context.Context) {
    // Get tracer from the app
    tracer := c.app.GetTracer()
    
    for {
        // Create a span for each collection cycle
        collectorSpan := tracer.NewCollectorSpan(ctx, "my-collector", "collect")
        
        // Use the context with the span
        ctx := collectorSpan.Context()
        
        // Your collection logic here
        if err := c.collectData(ctx); err != nil {
            collectorSpan.RecordError(err)
        }
        
        // Add events and attributes
        collectorSpan.AddEvent("collection_completed")
        collectorSpan.SetAttributes(
            attribute.String("collection.type", "scheduled"),
            attribute.Int("items.collected", 42),
        )
        
        collectorSpan.End()
        
        time.Sleep(30 * time.Second)
    }
}
```

### Tracing Configuration

- **enabled**: Set to `true` to enable tracing (default: `false`)
- **service_name**: Name of your service for trace identification
- **endpoint**: OTLP endpoint URL (required when tracing is enabled)
- **headers**: Additional headers for OTLP requests (optional)

## License

MIT
