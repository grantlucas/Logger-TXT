package cmd

import (
	"fmt"
	"io"
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
			path := config.ResolveFilePath(filePath)

			typeName, project := parseTypeProjectFlags(cmd)
			entryFilter := buildEntryFilter(typeName, project)

			startStr, endStr, hasRange, err := parseDateRangeFlags(cmd)
			if err != nil {
				return err
			}

			if hasRange {
				countChanged := cmd.Flags().Changed("count")
				return runSearchRange(cmd.OutOrStdout(), path, args[0], caseSensitive, count, countChanged, startStr, endStr, entryFilter)
			}

			lines, err := logger.Search(path, args[0], caseSensitive, count, entryFilter)
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
	addDateRangeFlags(cmd)
	addTypeProjectFlags(cmd)

	return cmd
}

func runSearchRange(out io.Writer, path, term string, caseSensitive bool, count int, countChanged bool, startStr, endStr string, entryFilter func(entry.Entry) bool) error {
	loc := time.Now().Location()

	start, end, err := entry.ParseDateRange(startStr, endStr, loc)
	if err != nil {
		return err
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
		if !strings.Contains(line, searchTerm) {
			return false
		}
		if entryFilter != nil && !entryFilter(e) {
			return false
		}
		return true
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
