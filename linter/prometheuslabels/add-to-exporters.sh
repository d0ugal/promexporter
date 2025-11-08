#!/bin/bash
# Script to add prometheuslabels linter to all exporter projects

set -euo pipefail

# Base directory for all exporters
CODE_DIR="${HOME}/Code"
PROMEXPORTER_DIR="${CODE_DIR}/promexporter"

# List of exporters (excluding promexporter itself)
EXPORTERS=(
    "filesystem-exporter"
    "mqtt-exporter"
    "zigbee2mqtt-exporter"
    "brother-exporter"
    "ghcr-exporter"
    "github-exporter"
    "internet-perf-exporter"
    "slzb-exporter"
)

# Template for the custom linter configuration
LINTER_CONFIG='# Custom linter: prometheuslabels
# This linter forbids calls to WithLabelValues from Prometheus metrics
# and suggests using With(prometheus.Labels{...}) instead.
# 
# The linter is built from the promexporter module dependency (see Makefile).
# Path is relative to the project root and will be built before linting.
# In CI, install with: go install github.com/d0ugal/promexporter/linter/prometheuslabels/cmd/prometheuslabels@latest
custom:
  prometheuslabels:
    path: ./bin/prometheuslabels
    description: Forbids calls to WithLabelValues from Prometheus metrics
    original-url: github.com/d0ugal/promexporter/linter/prometheuslabels
'

echo "Adding prometheuslabels linter to exporter projects..."
echo ""

for exporter in "${EXPORTERS[@]}"; do
    exporter_dir="${CODE_DIR}/${exporter}"
    golangci_file="${exporter_dir}/.golangci.yml"
    
    if [ ! -d "${exporter_dir}" ]; then
        echo "⚠️  Skipping ${exporter}: directory not found"
        continue
    fi
    
    if [ ! -f "${golangci_file}" ]; then
        echo "⚠️  Skipping ${exporter}: .golangci.yml not found"
        continue
    fi
    
    # Check if already configured
    if grep -q "prometheuslabels:" "${golangci_file}" 2>/dev/null; then
        echo "✓ ${exporter}: already configured"
        continue
    fi
    
    # Add the configuration before the formatters section
    if grep -q "^formatters:" "${golangci_file}"; then
        # Insert before formatters section
        awk -v config="${LINTER_CONFIG}" '
            /^formatters:/ {
                print config
            }
            { print }
        ' "${golangci_file}" > "${golangci_file}.tmp" && mv "${golangci_file}.tmp" "${golangci_file}"
        echo "✓ ${exporter}: added linter configuration"
    else
        # Append to end of file
        echo "" >> "${golangci_file}"
        echo "${LINTER_CONFIG}" >> "${golangci_file}"
        echo "✓ ${exporter}: added linter configuration (appended)"
    fi
done

echo ""
echo "Done! All exporters have been updated."
echo ""
echo "To test the linter, run:"
echo "  cd ${PROMEXPORTER_DIR}"
echo "  go build -o /tmp/prometheuslabels ./linter/prometheuslabels/cmd/prometheuslabels"
echo "  /tmp/prometheuslabels <exporter-path>/..."

