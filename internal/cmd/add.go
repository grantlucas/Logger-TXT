package cmd

import (
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [message]",
		Short: "Add a new log entry",
		Long:  "Add a new timestamped entry to the log file.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement in Phase 5
			return nil
		},
	}

	cmd.Flags().StringP("type", "t", "", "entry type (uppercased)")
	cmd.Flags().StringP("project", "p", "", "project name (uppercased)")

	return cmd
}
