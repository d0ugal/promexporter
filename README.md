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
```

### Environment Variables

```bash
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
LOG_LEVEL=info
LOG_FORMAT=json
METRICS_DEFAULT_INTERVAL=30s
```

## License

MIT
