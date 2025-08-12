package cmd

import (
	"os"

	"github.com/dyammarcano/clonr/internal/git"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clonr",
	Short: "clonr - A Git wrapper to clone, monitor, and manage repositories.",
	Long: `clonr is a command-line tool to efficiently clone, monitor, and manage multiple Git repositories.

Features:
- Clone repositories
- List registered repositories
- Remove repositories from the registry
- Monitor repositories via a built-in server
- Map existing local repositories to the registry

For more information on each command, use 'clonr [command] --help'.`,
	Args: cobra.MaximumNArgs(2),
	RunE: git.CloneRepo,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
