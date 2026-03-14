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

	start, startDateOnly, err := entry.ParseInputTime(startStr, loc)
	if err != nil {
		return fmt.Errorf("invalid --start value: %w", err)
	}

	end, endDateOnly, err := entry.ParseInputTime(endStr, loc)
	if err != nil {
		return fmt.Errorf("invalid --end value: %w", err)
	}

	if endDateOnly {
		end = entry.EndOfDay(end)
	}
	_ = startDateOnly

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
			startStr, _ := cmd.Flags().GetString("start")
			endStr, _ := cmd.Flags().GetString("end")
			path := config.ResolveFilePath(filePath)

			if startStr != "" || endStr != "" {
				if startStr == "" || endStr == "" {
					return fmt.Errorf("--start and --end must both be provided")
				}
				return runShowRange(cmd.OutOrStdout(), path, count, startStr, endStr)
			}

			return runShow(cmd.OutOrStdout(), path, count)
		},
	}

	cmd.Flags().IntP("count", "c", 10, "number of entries to display")
	cmd.Flags().String("start", "", "start date (DD/MM/YY or DD/MM/YY HH:MM)")
	cmd.Flags().String("end", "", "end date (DD/MM/YY or DD/MM/YY HH:MM)")

	return cmd
}
