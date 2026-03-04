package cmd

import (
	"github.com/grantlucas/Logger-TXT/internal/config"
	"github.com/spf13/cobra"
)

var filePath string

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "logger-txt",
		Short: "A simple command-line logging tool",
		Long:  "Logger-TXT is a simple command-line logging tool that allows you to log activities throughout the day to a portable text file with timestamps.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default action: run show with default count
			path := config.ResolveFilePath(filePath)
			return runShow(cmd.OutOrStdout(), path, 10)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "path to log file")

	rootCmd.AddCommand(newAddCmd())
	rootCmd.AddCommand(newShowCmd())
	rootCmd.AddCommand(newSearchCmd())
	rootCmd.AddCommand(newDeleteCmd())
	rootCmd.AddCommand(newVersionCmd())

	return rootCmd
}

// Execute runs the root command.
func Execute() error {
	return NewRootCmd().Execute()
}
