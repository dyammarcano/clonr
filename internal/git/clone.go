package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"clonr/internal/db"

	"github.com/spf13/cobra"
)

func CloneRepo(cmd *cobra.Command, url, path string) error {
	initDB, err := db.InitDB()
	if err != nil {
		return fmt.Errorf("starting server: %w", err)
	}

	savePath := filepath.Join(path, extractRepoName(url))

	runCmd := exec.Command("git", "clone", url, savePath)
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone error: %v - %s", err, string(output))
	}

	if err := initDB.SaveRepo(url, savePath); err != nil {
		return fmt.Errorf("error saving repo to database: %w", err)
	}

	cmd.Printf("Cloned repo at %s\n", savePath)

	return nil
}

func PullRepo(path string) error {
	cmd := exec.Command("git", "-C", path, "pull")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error git pull: %v, output: %s", err, string(output))
	}
	return nil
}

func extractRepoName(url string) string {
	parts := strings.Split(url, "/")
	last := parts[len(parts)-1]
	return strings.TrimSuffix(last, ".git")
}
