# Repository Guidelines

## Project Overview

`see-cli` is the Go command-line client for the S.EE content sharing platform. It manages short URLs, text snippets, files, bio pages, dynamic QR codes, account usage, and API tokens. The module path is `s.ee/cli`; the CLI uses Cobra and `github.com/sdotee/sdk.go` for API access.

Read [README.md](README.md) for user-facing command and configuration documentation.

## Architecture

- `main.go` passes build metadata to `cmd.Execute`.
- `cmd/root.go` owns persistent flags, environment defaults, SDK client initialization, version output, and top-level command registration.
- Feature commands live in the single `cmd` package: `shorturl.go`, `text.go`, `file.go`, `bio.go`, `qrcode.go`, `account.go`, `domains.go`, and `tags.go`.
- `cmd/utils.go` contains shared output, pagination, stdin/file reading, and text MIME validation helpers.
- Feature files define option structs and Cobra commands, register flags and children in `init`, and return errors from `RunE`.
- Configuration precedence is explicit CLI flags, then `SEE_API_KEY`, `SEE_BASE_URL`, and `SEE_TIMEOUT`, then SDK defaults.

## Implementation Conventions

- Follow existing Cobra patterns and keep command-specific behavior in the owning feature file.
- Use the shared `apiClient`; do not create per-command SDK clients.
- Set `Args` validation on commands and validate semantic requirements before making API calls.
- Wrap SDK failures with operation context using `%w`.
- Write command output through `cmd.OutOrStdout()` so it remains testable.
- Honor the global `--json` flag for new API-backed output. Use `printJSON` for data and `printMessage` for mutation responses where appropriate.
- Reuse `addPageFlag`, `readContent`, and other shared helpers instead of duplicating their behavior.
- Preserve stdin support for content and file workflows. Do not close streams owned by Cobra or the caller.
- Format changed Go files with `gofmt`. Avoid unrelated refactors and public CLI changes.

## Build and Test

```bash
make test                 # go test ./...
make build                # static binary with version/build-time ldflags
make all                  # clean, test, and build
go test ./cmd/ -run TestName
```

Run the narrowest relevant test first, then `make test` for changes that affect shared command behavior. Tests live beside source files in `cmd`; prefer unit tests that do not require API access. Add coverage for argument validation, parsing, flags, output modes, and reusable helpers when those behaviors change.

## Documentation and Release

- Keep [README.md](README.md) synchronized when commands, flags, configuration, or user-visible output change.
- Keep this file and [CLAUDE.md](CLAUDE.md) aligned; `CLAUDE.md` imports this file as the canonical agent guidance.
- CI builds with Go 1.27 on pushes and pull requests to `main`.
- Tags trigger GoReleaser. `make release` creates a local snapshot, and `make darwin_universal` builds a universal macOS binary.
