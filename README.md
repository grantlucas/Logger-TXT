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

For faster usage, add a shell alias:

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

```bash
logger-txt show
logger-txt show -c 20
```

Running `logger-txt` with no subcommand is equivalent to `logger-txt show`.

### Searching entries

```bash
logger-txt search "coffee"
logger-txt search --case-sensitive "API"
logger-txt search -c 5 "deploy"
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
