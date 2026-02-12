# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CLI client for the S.EE content sharing platform (short URLs, text snippets, file uploads). Built in Go 1.26+ using Cobra for CLI framework and `github.com/sdotee/sdk.go` (v1.1.1) for API interactions. Module path: `s.ee/cli`.

## Build & Development Commands

```bash
make build              # Build binary (CGO_ENABLED=0, injects version/build time via ldflags)
make install            # Build and install to GOPATH
make test               # Run all tests: go test ./...
make clean              # Remove build artifacts
make all                # clean + test + build
make darwin_universal   # Universal macOS binary (arm64 + amd64)
make release            # GoReleaser snapshot build
make build_docker_image # Docker multi-stage build
```

Run a single test:
```bash
go test ./cmd/ -run TestFunctionName
```

## Architecture

**Entry point**: `main.go` delegates to `cmd.Execute()`. Version info (`BuildVersion`, `BuildTime`) injected via ldflags at build time.

**Command structure** (`cmd/` package, all in one package):
- `root.go` — Root `see` command, global flags (`--api-key`, `--base-url`, `--timeout`, `--json`), API client initialization in `PersistentPreRunE`
- `shorturl.go` — `see shorturl {create,update,delete}` subcommands
- `text.go` — `see text {create,update,delete}` subcommands
- `file.go` — `see file {upload,delete,domains,history}` subcommands
- `domains.go` — `see domains`
- `tags.go` — `see tags`
- `utils.go` — Shared helpers: `printJSON`, `readContent` (file/stdin with MIME validation), `ensureTextContent`

**Pattern**: Each feature file defines an options struct, command variables, and registers flags in `init()`. All commands use `RunE` for error propagation.

**Configuration priority**: CLI flags → environment variables (`SEE_API_KEY`, `SEE_BASE_URL`, `SEE_TIMEOUT`) → SDK defaults.

## Testing

Tests live in `cmd/` alongside source files (`file_test.go`, `utils_test.go`). Tests cover command structure/flag verification and utility function unit tests (content validation, MIME detection, JSON output). No integration tests requiring API access.

## Release

Tags trigger GoReleaser via GitHub Actions (`.github/workflows/release.yml`), producing cross-platform binaries (linux/darwin/windows, amd64/arm64) and Linux packages (deb/rpm via nfpm, package name `see-cli`). CI builds run on push/PR to main (`.github/workflows/go.yml`). Docker image uses Go 1.26 with UTC timezone.
