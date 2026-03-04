# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with
code in this repository.

## Project Overview

Logger-TXT is a simple command-line logging tool that allows users to log
activities throughout the day to a portable text file with timestamps. The tool
supports categorization by type and project.

The project is being rewritten from Bash to Go (see `PLAN-v2.md`). Both the
original Bash script and the new Go implementation live in this repo during the
transition.

## Key Files

### Go (v2 - in progress)

- `cmd/logger-txt/main.go` - Entry point
- `internal/cmd/` - Cobra CLI commands (root, add, show, search, delete,
  version)
- `internal/entry/` - Entry struct, formatting, parsing
- `internal/logger/` - Core log file operations (append, tail, search, delete)
- `internal/config/` - Config resolution (flags > env > defaults)
- `Makefile` - Build, test, lint, coverage targets
- `PLAN-v2.md` - Detailed implementation plan

### Bash (v1 - reference implementation)

- `logger-txt` - Original Bash script
- `test-logger-txt.sh` - Bash test suite

## Commands

```bash
# Build the Go binary
make build

# Run all tests
make test

# Run go vet
make vet

# Run linter (requires golangci-lint)
make lint

# Run tests with coverage, fail if below 100%
make coverage

# Clean build artifacts
make clean

# Install binary to GOPATH/bin
make install

# Run the original Bash script
./logger-txt "message"

# Run the original Bash tests
./test-logger-txt.sh
```

## Architecture (Go v2)

### Package Structure

- **`internal/entry`** - `Entry` struct with `Format()` and `ParseEntry()`.
  Owns the log line format: `DD/MM/YY HH:MM +ZZZZ - TYPE (PROJECT) - Message`
- **`internal/logger`** - Pure file operations: `Append`, `Tail`, `Search`,
  `DeleteLast`, `EnsureFile`. No cobra dependency. Uses `io.Writer`/file paths.
- **`internal/config`** - Resolves log file path. Precedence: `--file` flag >
  `LOGGERTXT_PATH` env > `./log.txt`
- **`internal/cmd`** - Thin cobra wrappers that parse flags and call into
  `logger`/`entry`

### Design Principles

- All logic in `internal/` - no public Go API
- Core packages (`entry`, `logger`, `config`) have zero CLI dependencies
- Table-driven tests throughout
- 100% unit test coverage target on `internal/` packages

### Log Format (unchanged from v1)

```text
DD/MM/YY HH:MM +ZZZZ - TYPE (PROJECT) - Message text
DD/MM/YY HH:MM +ZZZZ - TYPE - Message text
DD/MM/YY HH:MM +ZZZZ - (PROJECT) - Message text
DD/MM/YY HH:MM +ZZZZ - Message text
```

### Environment Configuration

- `LOGGERTXT_PATH` - Log file location (env var)
- `--file` / `-f` - Log file location (flag, takes precedence over env)
- Default: `./log.txt` (current working directory)

## Development Workflow

**All code changes MUST use the `/tdd` skill.** When implementing new features,
fixing bugs, or refactoring code, always invoke the TDD skill to follow the
red-green-refactor cycle. Write failing tests first, then implement the minimal
code to pass, then refactor. After tests pass, commit the changes to create a
checkpoint before moving on to the next step.

## Development Notes

- Version injected at build time via ldflags (see Makefile)
- Cross-platform: use `filepath.Join`, `os.UserHomeDir()`, handle `\r\n`
- Tests use `t.TempDir()` for isolated file operations
- CI runs on ubuntu-latest and macos-latest
