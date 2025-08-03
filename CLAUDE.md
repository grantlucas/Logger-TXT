# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Logger-TXT is a simple command-line logging tool written in Bash that allows users to log activities throughout the day to a portable text file with timestamps. The tool supports categorization by type and project, making it easy to track work progress, personal tasks, purchases, or any time-sensitive events.

## Key Files

- `logger-txt` - Main executable Bash script containing all functionality
- `README.md` - Comprehensive documentation with installation and usage instructions  
- `LICENSE` - GNU General Public License v3
- `.github/release.yml` - GitHub release automation configuration

## Commands

Since this is a simple Bash script project, there are no build or test commands. The main script is executed directly:

```bash
# Make executable (if needed)
chmod +x logger-txt

# Basic usage examples
./logger-txt "Simple log entry"
./logger-txt -t work -p project "Categorized log entry"
./logger-txt -c 20  # Show last 20 log entries
./logger-txt -s "search term"  # Search logs
./logger-txt -x  # Delete last entry (with confirmation)
```

## Architecture

This is a single-file Bash application with a straightforward structure:

### Core Functions
- `usage()` - Shows brief usage syntax
- `help()` - Displays comprehensive help text
- `version()` - Shows version and author information
- `check_log_file()` - Ensures log file exists and is writable
- `search_log()` / `search_log_sensitive()` - Case-insensitive/sensitive log searching
- `deleteLast()` / `confirmDeleteLast()` - Remove last log entry with user confirmation

### Script Flow
1. Set default variables and parse environment (LOGGERTXT_PATH)
2. Process command-line options using getopts
3. Execute appropriate action based on options:
   - Add new log entry with optional type/project categorization
   - Display recent log entries (default: last 10)
   - Search existing logs
   - Delete last entry

### Log Format
Entries are stored in plain text with the format:
```
DD/MM/YY HH:MM TZ - [TYPE] [(PROJECT)] - Log message text
```

### Environment Configuration
- `LOGGERTXT_PATH` - Optional environment variable to specify log file location
- Default log location: `~/log.txt` or same directory as script

## Development Notes

- The script uses standard Bash features and should be compatible with most Unix-like systems
- No external dependencies beyond basic Unix tools (grep, tail, sed, mv, date)
- File operations include proper error checking and permission validation
- User input validation is minimal - designed for trusted local use
- The script maintains backward compatibility with existing log files