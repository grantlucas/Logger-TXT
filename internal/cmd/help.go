package cmd

const rootLong = `Logger-TXT logs daily activities to a plain text file with timestamps.

Each entry is a single line with a timestamp, optional type and project tags,
and a message. The log file is portable — it is just a text file you can read,
grep, back up, or version-control with standard tools.

When run without a subcommand, logger-txt behaves like "show" and prints the
10 most recent entries.`

const rootExample = `  # Add a simple entry
  logger-txt add Had coffee with the team

  # Add with type and project
  logger-txt add -t meeting -p acme Discussed Q3 roadmap

  # Show last 5 entries
  logger-txt show -c 5

  # Show all entries from a specific date
  logger-txt show --start 14/03/26 --end 14/03/26

  # Search entries
  logger-txt search roadmap

  # Search within a date range
  logger-txt search deploy --start 01/03/26 --end 14/03/26

  # Delete the last entry (with confirmation)
  logger-txt delete`

const addLong = `Add a new timestamped entry to the log file.

The -t (type) and -p (project) flags are automatically uppercased, so
"-t meeting" becomes "MEETING" in the log. The message is logged as-is.`

const addExample = `  # Simple message
  logger-txt add Had coffee with the team
  # => 04/03/26 14:30 -0500 - Had coffee with the team

  # With a type
  logger-txt add -t meeting Standup with the team
  # => 04/03/26 14:30 -0500 - MEETING - Standup with the team

  # With a project
  logger-txt add -p acme Deployed v2.1
  # => 04/03/26 14:30 -0500 - (ACME) - Deployed v2.1

  # With both type and project
  logger-txt add -t dev -p acme Fixed login bug
  # => 04/03/26 14:30 -0500 - DEV (ACME) - Fixed login bug`

const showLong = `Display the most recent entries from the log file, newest last.

Output goes to stdout with one entry per line, so it works well in a pipe:

  logger-txt show | grep MEETING
  logger-txt show -c 50 | wc -l

Use --start and --end to display entries within a date range. Dates use
DD/MM/YY format; add HH:MM for exact times (quote the value on the CLI).
Date-only --end values default to 23:59. Both flags must be provided together.

Running "logger-txt" with no subcommand is equivalent to "logger-txt show".`

const showExample = `  # Show the default last 10 entries
  logger-txt show

  # Show the last 25 entries
  logger-txt show -c 25

  # Show all entries from a single day
  logger-txt show --start 14/03/26 --end 14/03/26

  # Show entries for a date range (e.g. first two weeks of March)
  logger-txt show --start 01/03/26 --end 14/03/26

  # Show entries within a specific time window (quote to include time)
  logger-txt show --start "14/03/26 09:00" --end "14/03/26 17:00"

  # Show the last 5 entries from a date range
  logger-txt show --start 01/03/26 --end 14/03/26 -c 5

  # Equivalent — bare command defaults to show
  logger-txt`

const searchLong = `Search log entries for a term, case-insensitive by default.

The search term is matched anywhere in the full log line, including the
timestamp, type, project, and message. Use --case-sensitive for exact
case matching. Results are limited to the most recent matches (default 10).

Use --start and --end to restrict the search to a date range. Dates use
DD/MM/YY format; add HH:MM for exact times (quote the value on the CLI).
Date-only --end values default to 23:59. Both flags must be provided together.`

const searchExample = `  # Find all entries mentioning "deploy"
  logger-txt search deploy

  # Case-sensitive search
  logger-txt search --case-sensitive MEETING

  # Return up to 20 matches
  logger-txt search -c 20 bug

  # Search within a date range (e.g. all March meetings)
  logger-txt search meeting --start 01/03/26 --end 31/03/26

  # Search a specific day's afternoon entries for bugs
  logger-txt search bug --start "14/03/26 12:00" --end "14/03/26 17:00"

  # Combine date range with result limit
  logger-txt search deploy --start 01/01/26 --end 14/03/26 -c 5`

const deleteLong = `Remove the last entry from the log file.

The entry is shown and you are asked to confirm before it is deleted. You must
type an uppercase "Y" to confirm — any other input cancels the operation.
Use --yes (-y) to skip the prompt for scripting.

This command only removes the last line of the log file. To delete an
older entry, edit the file directly.`

const deleteExample = `  # Delete with interactive confirmation
  logger-txt delete

  # Skip confirmation (for scripts)
  logger-txt delete --yes`

const subcommandHelpTemplate = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

const rootHelpTemplate = `{{.Long}}

Log entry format:
  DD/MM/YY HH:MM -0700 - Message
  DD/MM/YY HH:MM -0700 - TYPE - Message
  DD/MM/YY HH:MM -0700 - (PROJECT) - Message
  DD/MM/YY HH:MM -0700 - TYPE (PROJECT) - Message

File resolution order:
  1. --file flag value
  2. ./log.txt if it exists in the current directory
  3. LOGGERTXT_PATH environment variable
  4. ./log.txt (default)

{{if .Example}}Examples:
{{.Example}}

{{end}}Usage:
  {{.UseLine}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
