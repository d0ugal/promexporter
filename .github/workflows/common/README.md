# Common Reusable Workflows

This directory contains reusable GitHub Actions workflows that are shared across all exporter repositories.

## Workflows

### `dev-build.yml`
Development build workflow that builds and pushes Docker images to GHCR with dev tags.

**Inputs:**
- `image_name` (required): Full image name (e.g., `ghcr.io/d0ugal/slzb-exporter`)

**Usage:**
```yaml
jobs:
  call-dev-build:
    uses: d0ugal/promexporter/.github/workflows/common/dev-build.yml@main
    with:
      image_name: ${{ github.repository }}
    secrets: inherit
```

### `ci.yml`
CI workflow that runs tests, linting, builds, and security scans.

**Inputs:**
- `binary_name` (required): Name of the binary to build (e.g., `slzb-exporter`)
- `go_version` (optional): Go version to use (default: `1.25.4`)

**Usage:**
```yaml
jobs:
  call-ci:
    uses: d0ugal/promexporter/.github/workflows/common/ci.yml@main
    with:
      binary_name: ${{ github.event.repository.name }}
      go_version: '1.25.4'
    secrets: inherit
```

### `release-please.yml`
Release workflow that uses release-please to create releases and build/push Docker images.

**Inputs:**
- `image_name` (required): Full image name (e.g., `ghcr.io/d0ugal/slzb-exporter`)

**Usage:**
```yaml
jobs:
  call-release-please:
    uses: d0ugal/promexporter/.github/workflows/common/release-please.yml@main
    with:
      image_name: ${{ github.repository }}
    secrets: inherit
```

**Required Secrets:**
- `RELEASE_TOKEN`: GitHub token with permissions to create releases and PRs

### `auto-format.yml`
Auto-format workflow that runs `make fmt` and commits formatting changes.

**Usage:**
```yaml
jobs:
  call-auto-format:
    uses: d0ugal/promexporter/.github/workflows/common/auto-format.yml@main
    secrets: inherit
```

## Versioning

These workflows are versioned alongside promexporter releases. To pin to a specific version, use the promexporter version tag:

```yaml
uses: d0ugal/promexporter/.github/workflows/common/dev-build.yml@v1.13.7
```

For the latest version, use `@main`:

```yaml
uses: d0ugal/promexporter/.github/workflows/common/dev-build.yml@main
```

## Maintenance

When updating these workflows:
1. Make changes in this directory
2. Test with one exporter repository first
3. Once verified, update all exporter repositories to use the new version
4. Consider creating a new promexporter release to tag the workflow changes

