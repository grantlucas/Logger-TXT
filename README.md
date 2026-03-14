# Logger TXT

Logger TXT is a simple command-line tool to log activities throughout the day
to a portable text file with timestamps. Options are available to log a
specific entry under a type and project. Whether you track purchases, what you
ate that day, progress on projects at work or all of the above and more, you
will always have a simple, solid way of storing that information and a tool
that gets out of your way to get it there.

## Installation

### Homebrew

```bash
brew install grantlucas/tap/logger-txt
```

### Download Binary

Download a pre-built binary from the
[GitHub Releases](https://github.com/grantlucas/Logger-TXT/releases) page.

### Build from Source

```bash
git clone https://github.com/grantlucas/Logger-TXT.git
cd Logger-TXT
make build
make install
```

## Quick Access

For faster usage, add an alias:

```bash
alias l="logger-txt"
```

## Usage

### Adding entries

```bash
logger-txt add "This is a general log entry"
logger-txt add -t personal "Entry with a type"
logger-txt add -p project "Entry with a project"
logger-txt add -t personal -p project "Entry with both"
```

### Showing recent entries

Use `show` to browse entries by type or project. For free-text keyword
matching, use `search` instead.

```bash
logger-txt show
logger-txt show -c 20
logger-txt show -t personal
logger-txt show -p myproject
logger-txt show --start "01/03/26" --end "14/03/26"
```

Running `logger-txt` with no subcommand is equivalent to `logger-txt show`.

### Searching entries

Use `search` for free-text keyword searches across log entries. To browse
entries by type or project without a search term, use `show -t` / `show -p`
instead.

The search term matches anywhere in the full log line, including the
timestamp, type, project, and message. Note that `search WORK` is not the
same as `show -t WORK` — the search term will also match entries that mention
"WORK" in their message body. Use `show -t` to filter strictly by the type
tag.

```bash
logger-txt search "coffee"
logger-txt search --case-sensitive "API"
logger-txt search -c 5 "deploy"
logger-txt search -t work "meeting"
logger-txt search --start "01/03/26" --end "14/03/26" "deploy"
```

### Filtering

Both `show` and `search` support the following filters:

- `--type` / `-t` - filter by entry type
- `--project` / `-p` - filter by project
- `--start` - start date (`DD/MM/YY` or `DD/MM/YY HH:MM`)
- `--end` - end date (`DD/MM/YY` or `DD/MM/YY HH:MM`)
- `--count` / `-c` - number of entries to display (default: 10)

Date ranges require both `--start` and `--end`. When a date range is active,
all matching entries are returned — the default count of 10 does not apply.
Pass `-c` explicitly to limit results within the range.

Filters can be combined:

```bash
logger-txt show -t work -p api --start "01/03/26" -c 50
logger-txt search -p backend --start "01/03/26" --end "07/03/26" "deploy"
```

### Deleting the last entry

```bash
logger-txt delete
logger-txt delete -y
```

### Version

```bash
logger-txt version
```

## Example Output in log.txt

```text
31/01/26 13:30 -0600 - PERSONAL (PROJECT) - This is a log note with a type and project
31/01/26 13:35 -0600 - PERSONAL - This is a log note with just a type
31/01/26 13:40 -0600 - (PROJECT) - This is a log note with just a project
31/01/26 13:45 -0600 - This is just a general event
```

## Configuration

The log file path is resolved in the following order:

1. `--file` / `-f` flag
2. `./log.txt` in the current directory (if the file already exists)
3. `LOGGERTXT_PATH` environment variable
4. `./log.txt` (default)

## Related Projects

- [Logger TXT for macOS](https://github.com/grantlucas/Logger-TXT-macos) -
  A native macOS GUI for Logger TXT

## Main Goals

The main goal of this project is to provide a simple logging tool which can be
accessed quickly from the command line. By storing all data in a plain text
file, you're not locked into always using this tool or limited to only viewing
log entries with it. The data portability that a text file offers between
tools, operating systems and environments is crucial to having a smooth
workflow that is extremely dependable.

## What do you use it for anyways?!?

Over time the act of logging will become habitual. Over the course of a day
you may log any of the following and anything else you deem important.

- Progress of tasks related to work and/or specific projects
  - Extremely handy when it comes to filling in hours with an employer as you
    can easily look up what projects were worked on, on that Tuesday two weeks
    ago.
- Progress of personal tasks or projects
  - Progress logging is the main use of this tool
- Purchases made
  - Extremely useful when the credit card bill comes with cryptic names of
    companies.
- Log important events or anything where the time that it happened is
  important.
  - Had an important conversation with someone? Log that you had it so you can
    also know when it exactly happened.
- Log anything!
