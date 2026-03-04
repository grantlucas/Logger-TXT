package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Set via ldflags at build time.
var (
	version = "dev"
	date    = "unknown"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			fmt.Fprintln(out, "Logger-TXT")
			fmt.Fprintf(out, "Version %s\n", version)
			fmt.Fprintln(out, "Author: Grant Lucas (contact@grantlucas.com)")
			fmt.Fprintf(out, "Last updated: %s\n", date)
			fmt.Fprintln(out, "Release date: 26/07/2010")
			fmt.Fprintln(out, "License: GPL, http://www.gnu.org/copyleft/gpl.html")
			if version != "dev" {
				fmt.Fprintf(out, "Release: https://github.com/grantlucas/Logger-TXT/releases/tag/v%s\n", version)
			}
		},
	}
}
