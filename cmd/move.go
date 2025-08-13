package cmd

import (
	"fmt"

	"github.com/dyammarcano/clonr/internal/git"
	"github.com/spf13/cobra"
)

var destFolder string

// moveCmd represents the move command
var moveCmd = &cobra.Command{
	Use:   "move <uid>",
	Short: "Move a local repository to another folder",
	Long: `Move a local cloned repository to a new folder using its UID.
Example:
    clonr move abc123 --dest /path/to/new/location`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if destFolder == "" {
			return fmt.Errorf("--dest is required")
		}

		repos, err := git.PrettiListRepos(false)
		if err != nil {
			return fmt.Errorf("failed to list repos: %w", err)
		}

		if err := git.MoveRepo(args[0], repos, destFolder); err != nil {
			return fmt.Errorf("failed to update database: %w", err)
		}

		return nil
	},
}

func init() {
	moveCmd.Flags().StringVar(&destFolder, "dest", "", "Destination folder to move the repository to")
	rootCmd.AddCommand(moveCmd)
}
