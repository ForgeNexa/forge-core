# Forge Core — Shared Go Systems Foundation

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)](LICENSE)
[![Go Module](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Race Detector](https://img.shields.io/badge/Race%20Detector-100%25%20PASS-brightgreen?style=flat-square)](https://golang.org)
[![Security](https://img.shields.io/badge/AES--256--GCM-Authenticated-blue?style=flat-square)](internal/security)
[![CI](https://github.com/ForgeNexa/forge-core/actions/workflows/ci.yml/badge.svg)](https://github.com/ForgeNexa/forge-core/actions/workflows/ci.yml)

> **"Rock-solid systems foundation for modern developer tools."**  
> An open-source Go module providing bounded concurrency worker pools, cryptographic security routines, atomic file operations, and Git porcelain helpers.

---

## 🏷️ Tags & Hashtags

`#ForgeNexa` `#ForgeCore` `#GoLang` `#Concurrency` `#Cryptography` `#Security` `#AtomicIO` `#DeveloperTools` `#DevTools` `#OpenSource`

---

## 📦 Packages

- **`internal/concurrency`**: Bounded worker pool (`Pool`) with context cancellation and graceful job draining.
- **`internal/security`**:
  - `SafePath`: Directory traversal boundary enforcement.
  - `Encrypt` / `Decrypt`: Authenticated AES-256-GCM credential encryption.
  - `MaskToken`: Redacts sensitive secrets and API tokens.
  - `AtomicWriteFile`: Atomic filesystem operations preventing partial file corruption.
- **`internal/git`**: High-performance porcelain wrappers for status parsing, commit histories, and branch queries.

---

## 🛠️ Testing & Verification

```bash
# Run all tests with Go race detector
go test -v -race ./...
```

---

## 🌐 Part of ForgeNexa

Forge Core powers all backend services across the [ForgeNexa Developer Tools Ecosystem](https://forgenexa.github.io/):
- [Forge IDE](https://github.com/ForgeNexa/forge-ide)
- [Forge API](https://github.com/ForgeNexa/forge-api)
- [Forge Terminal](https://github.com/ForgeNexa/forge-terminal)
- [Forge Git](https://github.com/ForgeNexa/forge-git)
- [Forge Hub](https://github.com/ForgeNexa/forge-hub)

---

## 📄 License

MIT © [ForgeNexa](https://forgenexa.github.io)
