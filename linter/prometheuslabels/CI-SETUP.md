# CI Setup for Prometheus Labels Linter

## Overview

The `prometheuslabels` linter is part of the `promexporter` module. Once `promexporter` is released with the linter included, exporters can use it via their existing dependency.

## Setup Options

### Option 1: Using `go install` (Recommended for CI)

Once `promexporter` is released with the linter, add this to your CI workflow:

```yaml
# GitHub Actions example
- name: Install prometheuslabels linter
  run: go install github.com/d0ugal/promexporter/linter/prometheuslabels/cmd/prometheuslabels@latest

- name: Run golangci-lint
  run: make lint
```

Then update `.golangci.yml` to use the installed binary:

```yaml
custom:
  prometheuslabels:
    path: $(go env GOPATH)/bin/prometheuslabels
    description: Forbids calls to WithLabelValues from Prometheus metrics
```

### Option 2: Build from Module Cache (Current Implementation)

The Makefile builds the linter from the module cache. This works once `promexporter` includes the linter:

```makefile
build-prometheuslabels-linter:
	go build -o bin/prometheuslabels github.com/d0ugal/promexporter/linter/prometheuslabels/cmd/prometheuslabels
```

### Option 3: Use Replace Directive (For Development)

For local development before the release, add to `go.mod`:

```go
replace github.com/d0ugal/promexporter => ../promexporter
```

Then the Makefile will build from the local source.

## CI Workflow Example

```yaml
name: Lint

on: [push, pull_request]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      
      - name: Install prometheuslabels linter
        run: go install github.com/d0ugal/promexporter/linter/prometheuslabels/cmd/prometheuslabels@latest
      
      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: v2.6.1
```

## Current Status

⚠️ **Note**: The linter is not yet in the published `promexporter` module. Once `promexporter` v1.12.0+ (or next release) includes the linter, the CI setup above will work.

For now, you can:
1. Use a replace directive for local development
2. Wait for the promexporter release
3. Use `go install` with a commit hash: `go install github.com/d0ugal/promexporter/linter/prometheuslabels/cmd/prometheuslabels@main`

