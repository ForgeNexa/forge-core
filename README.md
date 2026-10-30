# Forge Core (`forge-core`)

[![CI](https://github.com/ForgeNexa/forge-core/actions/workflows/ci.yml/badge.svg)](https://github.com/ForgeNexa/forge-core/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Forge Core** is the shared, high-performance Go foundation powering backend operations, process supervision, bounded concurrency, filesystem safety, and Git orchestration across the entire **ForgeNexa** developer tools ecosystem:

- **Forge IDE** (Code editor & workspace daemon)
- **Forge API** (Local-first HTTP/gRPC engine)
- **Forge Terminal** (PTY process manager)
- **Forge Git** (Topological commit graph & diff engine)
- **Forge Hub** (Ecosystem orchestration & launcher)

---

## Architecture & Subpackages

- `internal/concurrency/`: Bounded worker pools, context cancellation, and graceful task draining.
- `internal/security/`: Path traversal prevention (`SafePath`), AES-256-GCM secret encryption, secret redaction, and atomic file writes.
- `internal/git/`: Topological commit graph traversal, branch and staging orchestration, and merge conflict resolution.
- `internal/process/`: Cross-platform process supervision and clean teardown.
- `cmd/forge-core/`: Standalone core engine binary.

---

## Development & Testing

```bash
# Run unit and race-condition tests
go test -v -race ./...

# Build core daemon binary
go build -o bin/forge-core ./cmd/forge-core
```
