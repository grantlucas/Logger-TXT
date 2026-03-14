package cmd

import (
	"fmt"

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
