package cmd

import (
	"github.com/dyammarcano/clonr/internal/svc"
	"github.com/spf13/cobra"
)

// mapCmd represents the map command
var mapCmd = &cobra.Command{
	Use:   "map [directory]",
	Short: "Scan for existing Git repositories and add them to the clonr database.",
	Long: `Recursively search the specified directory (or current directory if not specified) for Git repositories.
For each found repository, add it to the clonr database if not already present. This allows you to manage and update
repositories that were not originally cloned with clonr.`,
	RunE: svc.MapRepos,
}

func init() {
	rootCmd.AddCommand(mapCmd)
}
