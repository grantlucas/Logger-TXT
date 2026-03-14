package cmd

import (
	"github.com/grantlucas/Logger-TXT/internal/config"
	"github.com/spf13/cobra"
)

var filePath string

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "logger-txt",
		Short:   "A simple command-line logging tool",
		Long:    rootLong,
		Example: rootExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default action: run show with default count
			path := config.ResolveFilePath(filePath)
			return runShow(cmd.OutOrStdout(), path, 10, nil)
		},
	}

	rootCmd.SetHelpTemplate(rootHelpTemplate)
	rootCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "path to log file")

	subCmds := []*cobra.Command{
		newAddCmd(),
		newShowCmd(),
		newSearchCmd(),
		newDeleteCmd(),
		newVersionCmd(),
	}
	for _, sub := range subCmds {
		sub.SetHelpTemplate(subcommandHelpTemplate)
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}

// Execute runs the root command.
func Execute() error {
	return NewRootCmd().Execute()
}
