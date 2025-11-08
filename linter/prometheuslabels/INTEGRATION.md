# Integrating Prometheus Labels Linter into Exporters

This guide shows how to add the shared `prometheuslabels` linter to any exporter project.

## Quick Setup

### 1. Update `.golangci.yml`

Add the custom linter configuration to your `.golangci.yml`:

```yaml
# Custom linter: prometheuslabels
# This linter forbids calls to WithLabelValues from Prometheus metrics
# and suggests using With(prometheus.Labels{...}) instead.
# 
# The linter is built from the promexporter module dependency (see Makefile).
custom:
  prometheuslabels:
    path: ./bin/prometheuslabels
    description: Forbids calls to WithLabelValues from Prometheus metrics
    original-url: github.com/d0ugal/promexporter/linter/prometheuslabels
```

### 2. Update Makefile

Add a target to build the linter (see `filesystem-exporter/Makefile` for reference):

```makefile
build-prometheuslabels-linter:
	@echo "Building prometheuslabels linter from promexporter module..."
	@mkdir -p bin
	@go build -o bin/prometheuslabels github.com/d0ugal/promexporter/linter/prometheuslabels/cmd/prometheuslabels
```

Then make your `lint`, `fmt`, and `lint-only` targets depend on it:

```makefile
lint: build-prometheuslabels-linter
	# ... your lint command
```

### 2. Test the Linter

Build and run the linter standalone:

```bash
cd /home/hoose/Code/promexporter
go build -o prometheuslabels ./linter/prometheuslabels/cmd/prometheuslabels
./prometheuslabels /path/to/your/exporter/...
```

### 3. Add to Makefile (Optional)

Add a lint target to your Makefile:

```makefile
.PHONY: lint-prometheus
lint-prometheus:
	@cd ../../promexporter && go build -o /tmp/prometheuslabels ./linter/prometheuslabels/cmd/prometheuslabels
	@/tmp/prometheuslabels ./...
```

## Exporter Projects

The following exporters should integrate this linter:

1. ✅ filesystem-exporter
2. ⬜ mqtt-exporter
3. ⬜ zigbee2mqtt-exporter
4. ⬜ brother-exporter
5. ⬜ ghcr-exporter
6. ⬜ github-exporter
7. ⬜ internet-perf-exporter
8. ⬜ slzb-exporter
9. ⬜ promexporter (if it has its own code to lint)

## Standalone Usage

If golangci-lint custom linter support doesn't work, you can:

1. **Build once and reuse:**
   ```bash
   cd /home/hoose/Code/promexporter
   go build -o ~/bin/prometheuslabels ./linter/prometheuslabels/cmd/prometheuslabels
   ```

2. **Add to PATH and use:**
   ```bash
   prometheuslabels ./...
   ```

3. **Use in CI/CD:**
   ```bash
   go install github.com/d0ugal/promexporter/linter/prometheuslabels/cmd/prometheuslabels@latest
   prometheuslabels ./...
   ```

