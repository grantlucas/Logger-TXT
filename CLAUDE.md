# CLAUDE.md

## Project Overview

Logger-TXT is a CLI tool for logging daily activities to a timestamped text
file with optional type and project categorization. Written in Go.

## Key Paths

- `cmd/logger-txt/main.go` - Entry point
- `internal/cmd/` - Cobra CLI commands
- `internal/entry/` - Entry struct, formatting, parsing
- `internal/logger/` - File operations (append, tail, search, delete)
- `internal/config/` - Log file path resolution
- `Makefile` - Build, test, lint, coverage targets

## Architecture

- All logic in `internal/` - no public Go API
- Core packages (`entry`, `logger`, `config`) have zero CLI dependencies
- `internal/cmd/` contains thin Cobra wrappers only
- Table-driven tests throughout
- 100% unit test coverage target on `internal/` packages

## Development Workflow

**All code changes MUST use the `/tdd` skill.** Write failing tests first,
implement minimal code to pass, then refactor. Commit after tests pass to
checkpoint before moving on.

### Essential Commands

- `make test` - Run all tests
- `make coverage` - Tests + fail if below 100% coverage

## Development Notes

- Version injected at build time via ldflags (see Makefile)
- Cross-platform: use `filepath.Join`, `os.UserHomeDir()`, handle `\r\n`
- Tests use `t.TempDir()` for isolated file operations
- CI runs on ubuntu-latest and macos-latest
- Prefer `ReverseLineScanner` over forward scanning (`bufio.Scanner`) when
  reading from the log file — log files can span 10+ years, and reverse
  scanning stops as soon as enough results are collected
