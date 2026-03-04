package cmd

import (
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete the last log entry",
		Long:  "Remove the last entry from the log file. Shows the entry and asks for confirmation before deleting.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement in Phase 5
			return nil
		},
	}

	cmd.Flags().BoolP("yes", "y", false, "skip confirmation prompt")

	return cmd
}
