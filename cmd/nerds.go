package cmd

import (
	"github.com/dyammarcano/clonr/internal/nerds"
	"github.com/spf13/cobra"
)

var nerdsCmd = &cobra.Command{
	Use:   "nerds",
	Short: "Interactively select and run nerd tools",
	Long: `Shows a list of available nerd tools from clonr and allows you 
to select one or more to run interactively.`,
	RunE: nerds.Tools,
}

func init() {
	rootCmd.AddCommand(nerdsCmd)
}
