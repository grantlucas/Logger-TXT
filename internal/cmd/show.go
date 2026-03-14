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

func runShow(out io.Writer, path string, count int, filter func(entry.Entry) bool) error {
	lines, err := logger.Tail(path, count, filter)
	if err != nil {
		return err
	}

	for _, line := range lines {
		fmt.Fprintln(out, line)
	}

	return nil
}

func runShowRange(out io.Writer, path string, count int, countChanged bool, startStr, endStr string, filter func(entry.Entry) bool) error {
	loc := time.Now().Location()

	start, end, err := entry.ParseDateRange(startStr, endStr, loc)
	if err != nil {
		return err
	}

	lines, err := logger.Range(path, start, end, filter)
	if err != nil {
		return err
	}

	if countChanged && len(lines) > count {
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

			typeName, project := parseTypeProjectFlags(cmd)
			filter := buildEntryFilter(typeName, project)

			startStr, endStr, hasRange, err := parseDateRangeFlags(cmd)
			if err != nil {
				return err
			}

			if hasRange {
				countChanged := cmd.Flags().Changed("count")
				return runShowRange(cmd.OutOrStdout(), path, count, countChanged, startStr, endStr, filter)
			}

			return runShow(cmd.OutOrStdout(), path, count, filter)
		},
	}

	cmd.Flags().IntP("count", "c", 10, "number of entries to display")
	addDateRangeFlags(cmd)
	addTypeProjectFlags(cmd)

	return cmd
}
