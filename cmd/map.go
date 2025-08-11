package cmd

import (
	"github.com/dyammarcano/clonr/internal/git"
	"github.com/spf13/cobra"
)

// mapCmd represents the map command
var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: git.MapRepos,
}

func init() {
	rootCmd.AddCommand(mapCmd)
}
