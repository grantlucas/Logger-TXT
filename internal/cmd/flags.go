package cmd

import (
	"fmt"
	"strings"

	"github.com/grantlucas/Logger-TXT/internal/entry"
	"github.com/spf13/cobra"
)

func addDateRangeFlags(cmd *cobra.Command) {
	cmd.Flags().String("start", "", "start date (DD/MM/YY or DD/MM/YY HH:MM)")
	cmd.Flags().String("end", "", "end date (DD/MM/YY or DD/MM/YY HH:MM)")
}

func parseDateRangeFlags(cmd *cobra.Command) (start, end string, hasRange bool, err error) {
	start, _ = cmd.Flags().GetString("start")
	end, _ = cmd.Flags().GetString("end")

	if start == "" && end == "" {
		return "", "", false, nil
	}

	if start == "" || end == "" {
		return "", "", false, fmt.Errorf("--start and --end must both be provided")
	}

	return start, end, true, nil
}

func addTypeProjectFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("type", "t", "", "filter by entry type")
	cmd.Flags().StringP("project", "p", "", "filter by project")
}

func parseTypeProjectFlags(cmd *cobra.Command) (typeName, project string) {
	typeName, _ = cmd.Flags().GetString("type")
	project, _ = cmd.Flags().GetString("project")
	return strings.ToUpper(typeName), strings.ToUpper(project)
}

func buildEntryFilter(typeName, project string) func(entry.Entry) bool {
	if typeName == "" && project == "" {
		return nil
	}
	return func(e entry.Entry) bool {
		if typeName != "" && !strings.EqualFold(e.Type, typeName) {
			return false
		}
		if project != "" && !strings.EqualFold(e.Project, project) {
			return false
		}
		return true
	}
}
