package cmd

import (
	"fmt"

	"github.com/grantlucas/Logger-TXT/internal/config"
	"github.com/grantlucas/Logger-TXT/internal/logger"
	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show recent log entries",
		Long:  "Display the most recent entries from the log file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			count, _ := cmd.Flags().GetInt("count")
			path := config.ResolveFilePath(filePath)

			lines, err := logger.Tail(path, count)
			if err != nil {
				return err
			}

			for _, line := range lines {
				fmt.Fprintln(cmd.OutOrStdout(), line)
			}

			return nil
		},
	}

	cmd.Flags().IntP("count", "c", 10, "number of entries to display")

	return cmd
}
