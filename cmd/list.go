package cmd

import (
	"fmt"

	"github.com/dyammarcano/clonr/internal/git"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered Git repositories.",
	Long:  `Display all Git repositories currently registered in the clonr database, showing their remote URL and local path.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		list, err := git.ListRepos()
		if err != nil {
			return err
		}

		for _, repo := range list {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s -> %s\n", repo.URL, repo.Path)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
