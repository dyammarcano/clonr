package cmd

import (
	"github.com/dyammarcano/clonr/internal/core"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: core.Init,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// init open menu for initialization
// * set default directory
// * set default editor
// * set default terminal
