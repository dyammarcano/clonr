package cmd

import (
	"github.com/dyammarcano/clonr/internal/svc"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered Git repositories.",
	Long:  `Display all Git repositories currently registered in the clonr database, showing their remote URL and local path.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := svc.PrettiListRepos(true)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
