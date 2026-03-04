package cmd

import (
	"fmt"
	"io"

	"github.com/grantlucas/Logger-TXT/internal/config"
	"github.com/grantlucas/Logger-TXT/internal/logger"
	"github.com/spf13/cobra"
)

func runShow(out io.Writer, path string, count int) error {
	lines, err := logger.Tail(path, count)
	if err != nil {
		return err
	}

	for _, line := range lines {
		fmt.Fprintln(out, line)
	}

	return nil
}

func newShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Short:   "Show recent log entries",
		Long:    showLong,
		Example: showExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			count, _ := cmd.Flags().GetInt("count")
			path := config.ResolveFilePath(filePath)
			return runShow(cmd.OutOrStdout(), path, count)
		},
	}

	cmd.Flags().IntP("count", "c", 10, "number of entries to display")

	return cmd
}
