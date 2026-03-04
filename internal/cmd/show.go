package cmd

import (
	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show recent log entries",
		Long:  "Display the most recent entries from the log file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement in Phase 5
			return nil
		},
	}

	cmd.Flags().IntP("count", "c", 10, "number of entries to display")

	return cmd
}
