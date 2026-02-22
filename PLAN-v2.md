# Logger-TXT v2: Go CLI Rewrite Plan

## 1. Why Go

The current bash script is ~255 lines with 20 bash-based tests. It works, but:
- Adding complex features (structured output, multiple log files, config files) is painful in bash
- Testing is fragile — shell test scripts can't mock, isolate, or test edge cases well
- Cross-platform support is impossible (Windows users are out)
- No type safety, no dependency management, no compiled binary distribution

Go solves all of these: single static binary, cross-compiles to every OS, excellent CLI libraries, built-in testing, and Homebrew has first-class support for Go formulae.

---

## 2. Feature Parity Checklist

Every behavior from v1 must work identically in v2 before any new features are added.

| v1 Feature | v1 Flag | v2 Command / Flag | Notes |
|---|---|---|---|
| Add log entry | positional args | `logger-txt add "message"` | Subcommand-based |
| Type categorization | `-t TYPE` | `--type TYPE` / `-t TYPE` | Uppercased |
| Project categorization | `-p PROJECT` | `--project PROJECT` / `-p PROJECT` | Uppercased, wrapped in `()` |
| Show recent entries | (default, no args) | `logger-txt show` | Default subcommand |
| Display count | `-c COUNT` | `--count COUNT` / `-c COUNT` on `show` | Default: 10 |
| Case-insensitive search | `-s text` | `logger-txt search "text"` | Default behavior |
| Case-sensitive search | `-S text` | `logger-txt search --case-sensitive "text"` | Opt-in flag |
| Delete last entry | `-x` | `logger-txt delete` | Interactive confirm |
| Custom file path | `-f path` | `--file PATH` / `-f PATH` | Global flag |
| Env var: LOGGERTXT_PATH | `$LOGGERTXT_PATH` | `$LOGGERTXT_PATH` | Same behavior |
| Help | `-h` | `--help` / `-h` | Built-in with cobra |
| Version | `-V` / `-v` | `--version` / `-v` | Built-in with cobra |

### Log Format (Unchanged)

The on-disk format must remain identical for backward compatibility:

```
DD/MM/YY HH:MM +ZZZZ - TYPE (PROJECT) - Message text
DD/MM/YY HH:MM +ZZZZ - TYPE - Message text
DD/MM/YY HH:MM +ZZZZ - (PROJECT) - Message text
DD/MM/YY HH:MM +ZZZZ - Message text
```

Existing log files must work with the new binary with zero migration.

---

## 3. Project Structure

```
Logger-TXT/
├── cmd/
│   └── logger-txt/
│       └── main.go              # Entry point, wires up root command
├── internal/
│   ├── cmd/
│   │   ├── root.go              # Root command, global flags (--file)
│   │   ├── add.go               # `add` subcommand
│   │   ├── add_test.go
│   │   ├── show.go              # `show` subcommand (default)
│   │   ├── show_test.go
│   │   ├── search.go            # `search` subcommand
│   │   ├── search_test.go
│   │   ├── delete.go            # `delete` subcommand
│   │   ├── delete_test.go
│   │   └── version.go           # `version` subcommand
│   ├── logger/
│   │   ├── logger.go            # Core logging logic (write, read, search, delete)
│   │   └── logger_test.go       # Unit tests for core logic
│   ├── entry/
│   │   ├── entry.go             # Entry struct, formatting, parsing
│   │   └── entry_test.go        # Unit tests for entry formatting
│   └── config/
│       ├── config.go            # Env var and flag resolution
│       └── config_test.go
├── logger-txt                   # Original bash script (kept in repo)
├── test-logger-txt.sh           # Original test script (kept in repo)
├── go.mod
├── go.sum
├── Makefile                     # Build, test, lint, install targets
├── goreleaser.yml               # Cross-platform release builds
├── README.md
├── LICENSE
└── .github/
    └── workflows/
        ├── ci.yml               # Test + lint on PRs
        ├── release.yml          # GoReleaser on tag push
        └── pr-verify.yml        # (existing, kept for bash script)
```

### Key Design Decisions

- **`internal/` package**: All logic is internal — no public Go API to maintain. The binary is the interface.
- **`logger/` package**: Pure logic layer. Takes an `io.Writer`/`io.Reader` or file path. No cobra dependency. Fully unit-testable.
- **`entry/` package**: Handles formatting a log entry string and parsing one from a line. This is where the `DD/MM/YY HH:MM +ZZZZ - TYPE (PROJECT) - Message` format lives.
- **`cmd/` package**: Thin cobra command wrappers that call into `logger/`. Only responsible for flag parsing and calling the right function.
- **`config/` package**: Resolves the effective config from env vars, flags, and defaults. Precedence: flag > env var > default.

---

## 4. Dependencies

| Dependency | Purpose | Rationale |
|---|---|---|
| [cobra](https://github.com/spf13/cobra) | CLI framework | Industry standard, subcommand support, auto-generated help, shell completions |
| Standard library only for core logic | File I/O, string formatting, time | No unnecessary deps in the hot path |

That's it. No logging library, no config library. Keep it minimal.

---

## 5. Implementation Phases

### Phase 1: Project Scaffolding
- Initialize Go module (`go mod init github.com/grantlucas/Logger-TXT`)
- Set up directory structure
- Add `Makefile` with targets: `build`, `test`, `lint`, `install`
- Add basic CI workflow (go test, go vet, staticcheck/golangci-lint)

### Phase 2: Core Entry Logic (`internal/entry/`)
- `Entry` struct: `Time`, `Type`, `Project`, `Message`
- `Entry.Format() string` — produces the exact log line format
- `ParseEntry(line string) (Entry, error)` — parses a log line back into an Entry
- Unit tests covering every format combination:
  - Message only
  - Type + message
  - Project + message
  - Type + project + message
  - Timestamp formatting with timezone

### Phase 3: Logger Operations (`internal/logger/`)
- `Append(path string, entry Entry) error` — append entry to file, create file if needed
- `Tail(path string, n int) ([]string, error)` — last N lines
- `Search(path string, term string, caseSensitive bool, limit int) ([]string, error)`
- `DeleteLast(path string) (string, error)` — returns the deleted line
- `EnsureFile(path string) error` — create + permission check
- Unit tests using `t.TempDir()` for isolated file operations

### Phase 4: Config Resolution (`internal/config/`)
- Resolve log file path: `--file` flag > `LOGGERTXT_PATH` env > `./log.txt` (current working directory)
- Unit tests for precedence behavior

### Phase 5: CLI Commands (`internal/cmd/`)

**Root command** (`root.go`):
- Global persistent flag: `--file` / `-f`
- Default action (no subcommand): runs `show`

```bash
# No args — shows last 10 entries (same as `logger-txt show`)
$ logger-txt
22/02/26 09:15 -0500 - WORK (API) - Fixed auth token refresh bug
22/02/26 09:30 -0500 - Grabbed a coffee
22/02/26 10:00 -0500 - WORK (API) - Deployed v1.3.2 to staging

# Global --file flag works with any subcommand
$ logger-txt --file ~/Dropbox/log.txt
$ logger-txt -f /tmp/test.txt add "Testing with custom file"
```

---

**`add` subcommand** (`add.go`):
- Flags: `--type` / `-t`, `--project` / `-p`
- Takes remaining args as the message
- Uppercases type and project
- Prints confirmation: `"message" logged under the type TYPE in the project PROJECT`

```bash
# Simple entry — no type, no project
$ logger-txt add "Grabbed a coffee"
"Grabbed a coffee" logged
# Writes: 22/02/26 10:30 -0500 - Grabbed a coffee

# With type only
$ logger-txt add -t work "Fixed login bug"
"Fixed login bug" logged under the type WORK
# Writes: 22/02/26 10:31 -0500 - WORK - Fixed login bug

# With project only
$ logger-txt add -p api "Deployed v1.3.2"
"Deployed v1.3.2" logged in the project API
# Writes: 22/02/26 10:32 -0500 - (API) - Deployed v1.3.2

# With both type and project (note: input is case-insensitive, stored uppercase)
$ logger-txt add -t work -p api "Reviewed pull request"
"Reviewed pull request" logged under the type WORK in the project API
# Writes: 22/02/26 10:33 -0500 - WORK (API) - Reviewed pull request

# Multi-word messages don't need quotes if unambiguous
$ logger-txt add -t personal Picked up groceries
"Picked up groceries" logged under the type PERSONAL
```

---

**`show` subcommand** (`show.go`):
- Flag: `--count` / `-c` (default 10)
- Prints last N lines of log file

```bash
# Show last 10 entries (default)
$ logger-txt show
22/02/26 08:00 -0500 - WORK (API) - Started morning standup
22/02/26 08:30 -0500 - WORK (API) - Fixed auth token refresh bug
...

# Show last 3 entries
$ logger-txt show -c 3
22/02/26 10:31 -0500 - WORK - Fixed login bug
22/02/26 10:32 -0500 - (API) - Deployed v1.3.2
22/02/26 10:33 -0500 - WORK (API) - Reviewed pull request

# Show last 20 entries
$ logger-txt show --count 20

# Remember: no args at all is equivalent to `show`
$ logger-txt
# (same as logger-txt show)
```

---

**`search` subcommand** (`search.go`):
- Flag: `--case-sensitive` (default false)
- Flag: `--count` / `-c` (default 10)
- Takes search term as argument
- Prints matching lines (last N matches)

```bash
# Case-insensitive search (default) — finds "work", "Work", "WORK"
$ logger-txt search "work"
22/02/26 08:00 -0500 - WORK (API) - Started morning standup
22/02/26 08:30 -0500 - WORK (API) - Fixed auth token refresh bug
22/02/26 10:31 -0500 - WORK - Fixed login bug
22/02/26 10:33 -0500 - WORK (API) - Reviewed pull request

# Case-sensitive search — only exact case
$ logger-txt search --case-sensitive "WORK"
22/02/26 08:00 -0500 - WORK (API) - Started morning standup
22/02/26 10:31 -0500 - WORK - Fixed login bug

# Limit search results to last 2 matches
$ logger-txt search -c 2 "api"
22/02/26 10:32 -0500 - (API) - Deployed v1.3.2
22/02/26 10:33 -0500 - WORK (API) - Reviewed pull request

# Search for a phrase
$ logger-txt search "pull request"
22/02/26 10:33 -0500 - WORK (API) - Reviewed pull request
```

---

**`delete` subcommand** (`delete.go`):
- Shows last line, prompts for confirmation (`Y` to confirm)
- Deletes last line if confirmed
- Accepts `--yes` / `-y` flag to skip confirmation (useful for scripting)

```bash
# Interactive delete — shows the line and asks for confirmation
$ logger-txt delete

Warning: You are removing the line below which appears at the end of the log file.

-------------------
22/02/26 10:33 -0500 - WORK (API) - Reviewed pull request
-------------------

Do you wish to continue? (Y/n)
Y

Deleted last line from file

# Cancelled delete
$ logger-txt delete
...
Do you wish to continue? (Y/n)
n

No line deleted

# Skip confirmation (for scripting)
$ logger-txt delete --yes
Deleted last line from file

# Works with custom file path
$ logger-txt -f /tmp/test.txt delete
```

---

**`version` subcommand** (or `--version` flag):
- Prints version, author, dates
- Version injected at build time via `-ldflags`

```bash
$ logger-txt version
Logger-TXT
Version 2.0.0
Author: Grant Lucas (contact@grantlucas.com)
Last updated: 22/02/2026
Release date: 26/07/2010
License: GPL, http://www.gnu.org/copyleft/gpl.html
Release: https://github.com/grantlucas/Logger-TXT/releases/tag/v2.0.0

# Also works as a flag
$ logger-txt --version
Logger-TXT version 2.0.0
```

---

#### v1 to v2 Command Mapping (Quick Reference)

```bash
# v1 (bash)                          →  v2 (Go)
logger-txt "message"                  →  logger-txt add "message"
logger-txt -t work "message"          →  logger-txt add -t work "message"
logger-txt -t work -p proj "message"  →  logger-txt add -t work -p proj "message"
logger-txt                            →  logger-txt show       (or just: logger-txt)
logger-txt -c 20                      →  logger-txt show -c 20
logger-txt -s "term"                  →  logger-txt search "term"
logger-txt -S "term"                  →  logger-txt search --case-sensitive "term"
logger-txt -x                         →  logger-txt delete
logger-txt -f path "message"          →  logger-txt -f path add "message"
logger-txt -V                         →  logger-txt version    (or: logger-txt --version)
logger-txt -h                         →  logger-txt --help     (or: logger-txt add --help, etc.)
```

### Phase 6: Integration Tests
- End-to-end tests that invoke the built binary
- Mirror every test case from `test-logger-txt.sh`:
  - Basic logging
  - Append (not overwrite)
  - Type categorization (uppercase)
  - Project categorization (uppercase, parentheses)
  - Type + project combined
  - Display functionality
  - Display count
  - Default 10-line limit
  - Search across multiple entries
  - Case-insensitive search
  - Case-sensitive search (match and no-match)
  - Delete with confirmation (Y, n, random)
  - File creation
  - Help and version output
  - Timestamp format
  - Paths with spaces

### Phase 7: Build & Release Pipeline
- **`Makefile`**: `build`, `test`, `lint`, `clean`, `install` targets
- **GoReleaser config** (`goreleaser.yml`):
  - Builds for: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, `windows/amd64`, `windows/arm64`
  - Injects version via ldflags
  - Creates GitHub release with archives
  - Generates checksums
  - Homebrew tap configuration (see below)
- **GitHub Actions** (`release.yml`):
  - Triggers on version tag push (`v*`)
  - Runs GoReleaser
  - Publishes to GitHub Releases

### Phase 8: Homebrew Distribution

**Option A — Your Own Tap (start here)**:
- Create repo `grantlucas/homebrew-tap`
- GoReleaser auto-publishes formula to the tap on release
- Users install with: `brew install grantlucas/tap/logger-txt`
- Low barrier, full control, immediate

**Option B — Homebrew Core (goal)**:
- Requirements from Homebrew:
  - Notable project (enough GitHub stars/users, or evidence of real-world usage)
  - Stable release history
  - Must build from source (Go formulae do this naturally)
  - No vendored dependencies that duplicate Homebrew packages
  - Must follow Homebrew formula conventions
- Process: Submit a PR to `homebrew-core` with the formula
- Realistically, start with your own tap and submit to core once the project has traction

**Recommendation**: Ship with your own tap first. It's zero friction and GoReleaser handles formula generation automatically. Submit to homebrew-core later as a separate effort when usage justifies it.

---

## 6. Testing Strategy

### Unit Tests (in every `_test.go` file)
- **Entry formatting**: Every combination of type/project/message
- **Entry parsing**: Valid lines, malformed lines, edge cases
- **Logger operations**: Append, tail, search, delete — all using `t.TempDir()`
- **Config resolution**: `--file` flag > `LOGGERTXT_PATH` env > `./log.txt`
- **Table-driven tests** throughout

### Integration Tests
- Build the binary, run it as a subprocess, assert on stdout/stderr/exit code/file contents
- Mirror the 20 existing bash test cases 1:1
- Add cross-platform path handling tests

### CI
- `go test ./...` on every PR
- `go vet ./...`
- `golangci-lint run`
- Tests run on ubuntu-latest (matches current CI); consider adding macos-latest and windows-latest matrix

### Coverage Target
- **100% unit test coverage** across all packages (`internal/logger/`, `internal/entry/`, `internal/config/`, `internal/cmd/`)
- Use `go test -coverprofile` and enforce with `-covermode=atomic`
- Add `make coverage` target that fails if coverage drops below 100%
- Integration tests provide additional confidence but are not counted toward the coverage gate

---

## 7. Cross-Platform Considerations

| Concern | Approach |
|---|---|
| File paths | Use `filepath.Join`, never hardcode `/` |
| Home directory | Use `os.UserHomeDir()` |
| Line endings | Write `\n` always; when reading, handle both `\n` and `\r\n` |
| Timestamps | Use Go's `time.Now().Format("02/01/06 15:04 -0700")` |
| Terminal colors | Only for interactive use, detect with `os.Stdout.Stat()` |
| File permissions | Use `os.OpenFile` with appropriate mode; skip `chmod` on Windows |

---

## 8. Migration Path for Existing Users

1. The Go binary is a drop-in replacement — same binary name (`logger-txt`), same `LOGGERTXT_PATH` env var, same log file format
2. Existing log files work without any changes
3. **Default file location changed**: v1 defaults to `~/log.txt`, v2 defaults to `./log.txt` (current working directory). Users who relied on the v1 default should set `LOGGERTXT_PATH` or pass `--file` explicitly.
4. **Old flags removed**: v1 flags like `-s`, `-S`, `-c`, `-x`, `-t`, `-p` are not carried forward. v2 uses subcommands (`add`, `show`, `search`, `delete`). No hidden aliases or deprecation shims — this is a clean break.
5. Document all breaking changes in the release notes and README

---

## 9. Implementation Order Summary

```
Phase 1: Scaffolding          → go.mod, directory structure, Makefile, CI
Phase 2: Entry formatting     → entry.go + entry_test.go
Phase 3: Logger operations    → logger.go + logger_test.go
Phase 4: Config resolution    → config.go + config_test.go
Phase 5: CLI commands         → cobra commands + command tests
Phase 6: Integration tests    → end-to-end binary tests
Phase 7: Build & release      → GoReleaser + GitHub Actions
Phase 8: Homebrew             → Own tap first, homebrew-core later
```

Each phase is independently testable. No phase depends on a later phase. The bash script and its tests remain in the repo throughout — they're the reference implementation until v2 is validated.
