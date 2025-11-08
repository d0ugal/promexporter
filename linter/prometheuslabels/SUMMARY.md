# Shared Prometheus Labels Linter - Summary

## What Was Done

The custom `prometheuslabels` linter has been moved to `promexporter` as a shared module that can be used across all 8 exporters and promexporter itself.

## Structure

```
promexporter/
└── linter/
    └── prometheuslabels/
        ├── analyzer.go                    # Reusable analyzer
        ├── cmd/
        │   └── prometheuslabels/
        │       └── main.go                # Standalone tool entry point
        ├── README.md                      # Usage documentation
        ├── INTEGRATION.md                 # Integration guide
        ├── SUMMARY.md                     # This file
        └── add-to-exporters.sh            # Script to add to all exporters
```

## Current Status

- ✅ Linter created in `promexporter/linter/prometheuslabels/`
- ✅ `filesystem-exporter` updated to use shared linter
- ✅ Integration script created
- ✅ Documentation created

## Adding to Other Exporters

### Option 1: Use the Script (Recommended)

```bash
cd /home/hoose/Code/promexporter/linter/prometheuslabels
./add-to-exporters.sh
```

This will automatically add the linter configuration to all exporter `.golangci.yml` files.

### Option 2: Manual Integration

For each exporter, add this to `.golangci.yml` before the `formatters:` section:

```yaml
# Custom linter: prometheuslabels
# This linter forbids calls to WithLabelValues from Prometheus metrics
# and suggests using With(prometheus.Labels{...}) instead.
custom:
  prometheuslabels:
    path: ../../promexporter/linter/prometheuslabels/cmd/prometheuslabels
    description: Forbids calls to WithLabelValues from Prometheus metrics
    original-url: github.com/d0ugal/promexporter/linter/prometheuslabels
```

## Testing

Build and test the linter:

```bash
cd /home/hoose/Code/promexporter
go build -o /tmp/prometheuslabels ./linter/prometheuslabels/cmd/prometheuslabels
/tmp/prometheuslabels /path/to/exporter/...
```

## Exporters to Update

1. ✅ filesystem-exporter - **DONE**
2. ⬜ mqtt-exporter
3. ⬜ zigbee2mqtt-exporter
4. ⬜ brother-exporter
5. ⬜ ghcr-exporter
6. ⬜ github-exporter
7. ⬜ internet-perf-exporter
8. ⬜ slzb-exporter
9. ⬜ promexporter (if it has code to lint)

## Benefits

- **Single source of truth**: Linter code lives in one place
- **Easy updates**: Update once, affects all exporters
- **No duplication**: No need to copy-paste code
- **Consistent**: All exporters use the same linter version

## Next Steps

1. Run `./add-to-exporters.sh` to add to all exporters
2. Test each exporter's linting setup
3. Commit changes to each exporter repository
4. Update CI/CD pipelines if needed

