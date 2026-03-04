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

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [message]",
		Short: "Add a new log entry",
		Long:  "Add a new timestamped entry to the log file.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			typeName, _ := cmd.Flags().GetString("type")
			project, _ := cmd.Flags().GetString("project")
			path := config.ResolveFilePath(filePath)

			typeName = strings.ToUpper(typeName)
			project = strings.ToUpper(project)
			message := strings.Join(args, " ")

			e := entry.Entry{
				Time:    time.Now(),
				Type:    typeName,
				Project: project,
				Message: message,
			}

			if err := logger.Append(path, e); err != nil {
				return err
			}

			// Print confirmation
			confirmation := fmt.Sprintf("%q logged", message)
			if typeName != "" && project != "" {
				confirmation = fmt.Sprintf("%q logged under the type %s in the project %s", message, typeName, project)
			} else if typeName != "" {
				confirmation = fmt.Sprintf("%q logged under the type %s", message, typeName)
			} else if project != "" {
				confirmation = fmt.Sprintf("%q logged in the project %s", message, project)
			}
			fmt.Fprintln(cmd.OutOrStdout(), confirmation)

			return nil
		},
	}

	cmd.Flags().StringP("type", "t", "", "entry type (uppercased)")
	cmd.Flags().StringP("project", "p", "", "project name (uppercased)")

	return cmd
}
