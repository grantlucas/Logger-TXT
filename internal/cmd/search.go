package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/grantlucas/Logger-TXT/internal/config"
	"github.com/grantlucas/Logger-TXT/internal/entry"
	"github.com/grantlucas/Logger-TXT/internal/logger"
	"github.com/spf13/cobra"
)

func newSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "search [term]",
		Short:   "Search log entries",
		Long:    searchLong,
		Example: searchExample,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			caseSensitive, _ := cmd.Flags().GetBool("case-sensitive")
			count, _ := cmd.Flags().GetInt("count")
			startStr, _ := cmd.Flags().GetString("start")
			endStr, _ := cmd.Flags().GetString("end")
			path := config.ResolveFilePath(filePath)

			if startStr != "" || endStr != "" {
				if startStr == "" || endStr == "" {
					return fmt.Errorf("--start and --end must both be provided")
				}
				return runSearchRange(cmd, path, args[0], caseSensitive, count, startStr, endStr)
			}

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
	cmd.Flags().String("start", "", "start date (DD/MM/YY or DD/MM/YY HH:MM)")
	cmd.Flags().String("end", "", "end date (DD/MM/YY or DD/MM/YY HH:MM)")

	return cmd
}

func runSearchRange(cmd *cobra.Command, path, term string, caseSensitive bool, count int, startStr, endStr string) error {
	loc := time.Now().Location()

	start, _, err := entry.ParseInputTime(startStr, loc)
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

	searchTerm := term
	if !caseSensitive {
		searchTerm = strings.ToLower(searchTerm)
	}

	filter := func(e entry.Entry) bool {
		line := e.Format()
		if !caseSensitive {
			line = strings.ToLower(line)
		}
		return strings.Contains(line, searchTerm)
	}

	lines, err := logger.Range(path, start, end, filter)
	if err != nil {
		return err
	}

	if len(lines) > count {
		lines = lines[len(lines)-count:]
	}

	for _, line := range lines {
		fmt.Fprintln(cmd.OutOrStdout(), line)
	}

	return nil
}
