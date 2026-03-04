package cmd

import (
	"fmt"

	"github.com/grantlucas/Logger-TXT/internal/config"
	"github.com/grantlucas/Logger-TXT/internal/logger"
	"github.com/spf13/cobra"
)

func newSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "search [term]",
		Short:   "Search log entries",
		Long:    searchLong,
		Example: searchExample,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			caseSensitive, _ := cmd.Flags().GetBool("case-sensitive")
			count, _ := cmd.Flags().GetInt("count")
			path := config.ResolveFilePath(filePath)

			lines, err := logger.Search(path, args[0], caseSensitive, count)
			if err != nil {
				return err
			}

			for _, line := range lines {
				fmt.Fprintln(cmd.OutOrStdout(), line)
			}

			return nil
		},
	}

	cmd.Flags().Bool("case-sensitive", false, "perform case-sensitive search")
	cmd.Flags().IntP("count", "c", 10, "maximum number of results")

	return cmd
}
