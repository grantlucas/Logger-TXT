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
			fmt.Println("Logger-TXT")
			fmt.Printf("Version %s\n", version)
			fmt.Println("Author: Grant Lucas (contact@grantlucas.com)")
			fmt.Printf("Last updated: %s\n", date)
			fmt.Println("Release date: 26/07/2010")
			fmt.Println("License: GPL, http://www.gnu.org/copyleft/gpl.html")
			if version != "dev" {
				fmt.Printf("Release: https://github.com/grantlucas/Logger-TXT/releases/tag/v%s\n", version)
			}
		},
	}
}
