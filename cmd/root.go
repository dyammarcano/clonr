package cmd

import (
	"os"

	"clonr/internal/git"

	"github.com/spf13/cobra"
)

var (
	path string
)

var rootCmd = &cobra.Command{
	Use:   "clonr",
	Short: "clonr - git wrapper para clonar y monitorear repos",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repoURL := args[0]
		return git.CloneRepo(cmd, repoURL, path)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&path, "path", "p", "", "Ruta donde se clonará el repositorio")
	_ = rootCmd.MarkFlagRequired("path")
}
