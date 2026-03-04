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

  # Search entries
  logger-txt search roadmap

  # Delete the last entry (with confirmation)
  logger-txt delete`

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
