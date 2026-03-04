package cmd

import (
	"github.com/spf13/cobra"
)

func newSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [term]",
		Short: "Search log entries",
		Long:  "Search for entries matching the given term. Case-insensitive by default.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement in Phase 5
			return nil
		},
	}

	cmd.Flags().Bool("case-sensitive", false, "perform case-sensitive search")
	cmd.Flags().IntP("count", "c", 10, "maximum number of results")

	return cmd
}
