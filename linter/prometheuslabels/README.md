# Prometheus Labels Linter

A custom Go linter that forbids calls to `WithLabelValues` from Prometheus metrics and suggests using `With(prometheus.Labels{...})` instead.

## Usage

### Standalone

Build and run the linter directly:

```bash
go build -o prometheuslabels ./linter/prometheuslabels/cmd/prometheuslabels
./prometheuslabels ./path/to/file.go
```

### With golangci-lint

Add to `.golangci.yml` in each exporter project:

```yaml
custom:
  prometheuslabels:
    path: ../../promexporter/linter/prometheuslabels/cmd/prometheuslabels
    description: Forbids calls to WithLabelValues from Prometheus metrics
```

Or use an absolute path:

```yaml
custom:
  prometheuslabels:
    path: /home/hoose/Code/promexporter/linter/prometheuslabels/cmd/prometheuslabels
    description: Forbids calls to WithLabelValues from Prometheus metrics
```

## Example

**Forbidden:**
```go
metric.WithLabelValues("label1", "label2", "label3").Inc()
```

**Suggested:**
```go
metric.With(prometheus.Labels{
    "label1": "value1",
    "label2": "value2", 
    "label3": "value3",
}).Inc()
```

## Implementation

The linter uses the `go/analysis` framework and is structured as:
- `analyzer.go`: Contains the reusable `Analyzer` variable
- `cmd/prometheuslabels/main.go`: Uses `singlechecker` to create a standalone tool

