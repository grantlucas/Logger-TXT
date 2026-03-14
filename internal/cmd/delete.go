package cmd

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/grantlucas/Logger-TXT/internal/config"
	"github.com/grantlucas/Logger-TXT/internal/logger"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete the last log entry",
		Long:    deleteLong,
		Example: deleteExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			yes, _ := cmd.Flags().GetBool("yes")
			path := config.ResolveFilePath(filePath)

			// Preview the last line
			lines, err := logger.Tail(path, 1, nil)
			if err != nil {
				return err
			}
			if len(lines) == 0 {
				return fmt.Errorf("file is empty")
			}

			lastLine := lines[0]
			out := cmd.OutOrStdout()

			if !yes {
				fmt.Fprintln(out)
				fmt.Fprintln(out, "Warning: You are removing the line below which appears at the end of the log file.")
				fmt.Fprintln(out)
				fmt.Fprintln(out, "-------------------")
				fmt.Fprintln(out, lastLine)
				fmt.Fprintln(out, "-------------------")
				fmt.Fprintln(out)
				fmt.Fprint(out, "Do you wish to continue? (Y/n) ")

				reader := bufio.NewReader(cmd.InOrStdin())
				answer, _ := reader.ReadString('\n')
				answer = strings.TrimSpace(answer)

				if answer != "Y" {
					fmt.Fprintln(out)
					fmt.Fprintln(out, "No line deleted")
					return nil
				}
			}

			_, err = logger.DeleteLast(path)
			if err != nil {
				return err
			}

			fmt.Fprintln(out)
			fmt.Fprintln(out, "Deleted last line from file")

			return nil
		},
	}

	cmd.Flags().BoolP("yes", "y", false, "skip confirmation prompt")

	return cmd
}
