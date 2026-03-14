package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/grantlucas/Logger-TXT/internal/config"
	"github.com/grantlucas/Logger-TXT/internal/entry"
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

func runShowRange(out io.Writer, path string, count int, startStr, endStr string) error {
	loc := time.Now().Location()

	start, end, err := entry.ParseDateRange(startStr, endStr, loc)
	if err != nil {
		return err
	}

	lines, err := logger.Range(path, start, end, nil)
	if err != nil {
		return err
	}

	if len(lines) > count {
		lines = lines[len(lines)-count:]
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

			startStr, endStr, hasRange, err := parseDateRangeFlags(cmd)
			if err != nil {
				return err
			}

			if hasRange {
				return runShowRange(cmd.OutOrStdout(), path, count, startStr, endStr)
			}

			return runShow(cmd.OutOrStdout(), path, count)
		},
	}

	cmd.Flags().IntP("count", "c", 10, "number of entries to display")
	addDateRangeFlags(cmd)

	return cmd
}
