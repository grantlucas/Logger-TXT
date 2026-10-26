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
		Use:     "add [message]",
		Short:   "Add a new log entry",
		Long:    addLong,
		Example: addExample,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			typeName, _ := cmd.Flags().GetString("type")
			project, _ := cmd.Flags().GetString("project")
			at, _ := cmd.Flags().GetString("at")
			backfill := cmd.Flags().Changed("at")
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
			if err := e.Validate(); err != nil {
				return err
			}

			// Backfilled entries are inserted in chronological position
			write := func() error { return logger.Append(path, e, cmd.ErrOrStderr()) }
			if backfill {
				t, err := entry.ParseAtTime(at, e.Time.Location(), e.Time)
				if err != nil {
					return err
				}
				e.Time = t
				write = func() error { return logger.Insert(path, e) }
			}

			if err := write(); err != nil {
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
			if backfill {
				confirmation += " at " + e.Time.Format(entry.TimeLayout)
			}
			fmt.Fprintln(cmd.OutOrStdout(), confirmation)

			return nil
		},
	}

	cmd.Flags().StringP("type", "t", "", "entry type (uppercased)")
	cmd.Flags().StringP("project", "p", "", "project name (uppercased)")
	cmd.Flags().String("at", "", "backfill time (DD/MM/YY HH:MM); inserted in chronological order")

	return cmd
}
